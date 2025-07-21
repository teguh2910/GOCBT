package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// Migrator handles database migrations
type Migrator struct {
	db     *sql.DB
	driver string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, driver string) *Migrator {
	return &Migrator{
		db:     db,
		driver: driver,
	}
}

// CreateMigrationsTable creates the migrations tracking table
func (m *Migrator) CreateMigrationsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`

	if m.driver == "postgres" {
		query = `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`
	}

	_, err := m.db.Exec(query)
	return err
}

// GetAppliedMigrations returns a list of applied migration versions
func (m *Migrator) GetAppliedMigrations() ([]int, error) {
	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var versions []int
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, rows.Err()
}

// LoadMigrations loads migration files from the migrations directory
func (m *Migrator) LoadMigrations(migrationsDir string) ([]Migration, error) {
	var migrations []Migration

	// Choose the appropriate directory based on driver
	dir := migrationsDir
	if m.driver == "postgres" {
		postgresDir := filepath.Join(migrationsDir, "postgres")
		if _, err := ioutil.ReadDir(postgresDir); err == nil {
			dir = postgresDir
		}
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Extract version from filename (e.g., "001_create_users.sql" -> 1)
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			continue
		}

		var version int
		if _, err := fmt.Sscanf(parts[0], "%d", &version); err != nil {
			continue
		}

		// Read migration content
		content, err := ioutil.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, Migration{
			Version: version,
			Name:    strings.TrimSuffix(file.Name(), ".sql"),
			SQL:     string(content),
		})
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// Up runs pending migrations
func (m *Migrator) Up(migrationsDir string) error {
	if err := m.CreateMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	appliedVersions, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	migrations, err := m.LoadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	appliedMap := make(map[int]bool)
	for _, version := range appliedVersions {
		appliedMap[version] = true
	}

	for _, migration := range migrations {
		if appliedMap[migration.Version] {
			continue // Skip already applied migrations
		}

		fmt.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)

		// Execute migration in a transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}

		if _, err := tx.Exec(migration.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d: %w", migration.Version, err)
		}

		// Record migration as applied
		var recordQuery string
		if m.driver == "postgres" {
			recordQuery = "INSERT INTO schema_migrations (version) VALUES ($1)"
		} else {
			recordQuery = "INSERT INTO schema_migrations (version) VALUES (?)"
		}

		if _, err := tx.Exec(recordQuery, migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		fmt.Printf("Migration %d applied successfully\n", migration.Version)
	}

	return nil
}
