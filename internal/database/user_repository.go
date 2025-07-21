package database

import (
	"gocbt/internal/models"
	"time"
)

// UserRepository implements the models.UserRepository interface
type UserRepository struct {
	db *DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) models.UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	if r.db.Driver == "postgres" {
		query = `
			INSERT INTO users (username, email, password_hash, first_name, last_name, role, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, created_at, updated_at
		`
	}

	if r.db.Driver == "postgres" {
		err := r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash,
			user.FirstName, user.LastName, user.Role, user.IsActive).Scan(
			&user.ID, &user.CreatedAt, &user.UpdatedAt)
		return err
	}

	result, err := r.db.Exec(query, user.Username, user.Email, user.PasswordHash,
		user.FirstName, user.LastName, user.Role, user.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
			FROM users WHERE id = $1
		`
	}

	row := r.db.QueryRow(query, id)
	return models.ScanUser(row)
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users WHERE username = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
			FROM users WHERE username = $1
		`
	}

	row := r.db.QueryRow(query, username)
	return models.ScanUser(row)
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users WHERE email = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
			FROM users WHERE email = $1
		`
	}

	row := r.db.QueryRow(query, email)
	return models.ScanUser(row)
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users 
		SET username = ?, email = ?, first_name = ?, last_name = ?, role = ?, is_active = ?, updated_at = ?
		WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE users 
			SET username = $1, email = $2, first_name = $3, last_name = $4, role = $5, is_active = $6, updated_at = $7
			WHERE id = $8
		`
	}

	user.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, user.Username, user.Email, user.FirstName,
		user.LastName, user.Role, user.IsActive, user.UpdatedAt, user.ID)
	return err
}

// Delete deletes a user
func (r *UserRepository) Delete(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	if r.db.Driver == "postgres" {
		query = "DELETE FROM users WHERE id = $1"
	}

	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves a list of users with pagination
func (r *UserRepository) List(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
			FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2
		`
	}

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := models.ScanUser(rows)
		if err != nil {
			return nil, err
		}
		if user != nil {
			users = append(users, user)
		}
	}

	return users, rows.Err()
}

// GetByRole retrieves users by role with pagination
func (r *UserRepository) GetByRole(role models.UserRole, limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
		FROM users WHERE role = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, username, email, password_hash, first_name, last_name, role, is_active, created_at, updated_at
			FROM users WHERE role = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
		`
	}

	rows, err := r.db.Query(query, role, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := models.ScanUser(rows)
		if err != nil {
			return nil, err
		}
		if user != nil {
			users = append(users, user)
		}
	}

	return users, rows.Err()
}
