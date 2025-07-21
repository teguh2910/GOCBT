package database

import (
	"database/sql"
	"fmt"
	"gocbt/internal/config"
	"time"

	_ "github.com/lib/pq"           // PostgreSQL driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB holds the database connection
type DB struct {
	*sql.DB
	Driver string
}

// Connect establishes a database connection based on configuration
func Connect(cfg *config.DatabaseConfig) (*DB, error) {
	var dsn string
	var err error

	switch cfg.Driver {
	case "sqlite":
		dsn = cfg.FilePath
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool with security considerations
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute) // Rotate connections regularly
	db.SetConnMaxIdleTime(1 * time.Minute) // Close idle connections quickly

	return &DB{
		DB:     db,
		Driver: cfg.Driver,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// RunMigrations runs database migrations
func (db *DB) RunMigrations(migrationsDir string) error {
	migrator := NewMigrator(db.DB, db.Driver)
	return migrator.Up(migrationsDir)
}
