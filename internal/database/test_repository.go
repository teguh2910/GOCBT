package database

import (
	"gocbt/internal/models"
	"time"
)

// TestRepository implements the models.TestRepository interface
type TestRepository struct {
	db *DB
}

// NewTestRepository creates a new test repository
func NewTestRepository(db *DB) models.TestRepository {
	return &TestRepository{db: db}
}

// Create creates a new test
func (r *TestRepository) Create(test *models.Test) error {
	query := `
		INSERT INTO tests (title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if r.db.Driver == "postgres" {
		query = `
			INSERT INTO tests (title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, created_at, updated_at
		`
	}

	if r.db.Driver == "postgres" {
		err := r.db.QueryRow(query, test.Title, test.Description, test.CreatedBy,
			test.DurationMinutes, test.TotalMarks, test.PassingMarks, test.Instructions,
			test.IsActive, test.StartTime, test.EndTime).Scan(
			&test.ID, &test.CreatedAt, &test.UpdatedAt)
		return err
	}

	result, err := r.db.Exec(query, test.Title, test.Description, test.CreatedBy,
		test.DurationMinutes, test.TotalMarks, test.PassingMarks, test.Instructions,
		test.IsActive, test.StartTime, test.EndTime)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	test.ID = int(id)
	test.CreatedAt = time.Now()
	test.UpdatedAt = time.Now()
	return nil
}

// GetByID retrieves a test by ID
func (r *TestRepository) GetByID(id int) (*models.Test, error) {
	query := `
		SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
		FROM tests WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
			FROM tests WHERE id = $1
		`
	}

	row := r.db.QueryRow(query, id)
	return models.ScanTest(row)
}

// Update updates a test
func (r *TestRepository) Update(test *models.Test) error {
	query := `
		UPDATE tests 
		SET title = ?, description = ?, duration_minutes = ?, total_marks = ?, passing_marks = ?, instructions = ?, is_active = ?, start_time = ?, end_time = ?, updated_at = ?
		WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE tests 
			SET title = $1, description = $2, duration_minutes = $3, total_marks = $4, passing_marks = $5, instructions = $6, is_active = $7, start_time = $8, end_time = $9, updated_at = $10
			WHERE id = $11
		`
	}

	test.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, test.Title, test.Description, test.DurationMinutes,
		test.TotalMarks, test.PassingMarks, test.Instructions, test.IsActive,
		test.StartTime, test.EndTime, test.UpdatedAt, test.ID)
	return err
}

// Delete deletes a test
func (r *TestRepository) Delete(id int) error {
	query := "DELETE FROM tests WHERE id = ?"
	if r.db.Driver == "postgres" {
		query = "DELETE FROM tests WHERE id = $1"
	}

	_, err := r.db.Exec(query, id)
	return err
}

// List retrieves a list of tests with pagination
func (r *TestRepository) List(limit, offset int) ([]*models.Test, error) {
	query := `
		SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
		FROM tests ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
			FROM tests ORDER BY created_at DESC LIMIT $1 OFFSET $2
		`
	}

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*models.Test
	for rows.Next() {
		test, err := models.ScanTest(rows)
		if err != nil {
			return nil, err
		}
		if test != nil {
			tests = append(tests, test)
		}
	}

	return tests, rows.Err()
}

// GetByCreator retrieves tests by creator with pagination
func (r *TestRepository) GetByCreator(creatorID int, limit, offset int) ([]*models.Test, error) {
	query := `
		SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
		FROM tests WHERE created_by = ? ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
			FROM tests WHERE created_by = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3
		`
	}

	rows, err := r.db.Query(query, creatorID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*models.Test
	for rows.Next() {
		test, err := models.ScanTest(rows)
		if err != nil {
			return nil, err
		}
		if test != nil {
			tests = append(tests, test)
		}
	}

	return tests, rows.Err()
}

// GetActiveTests retrieves active tests with pagination
func (r *TestRepository) GetActiveTests(limit, offset int) ([]*models.Test, error) {
	query := `
		SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
		FROM tests WHERE is_active = true ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
			FROM tests WHERE is_active = true ORDER BY created_at DESC LIMIT $1 OFFSET $2
		`
	}

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*models.Test
	for rows.Next() {
		test, err := models.ScanTest(rows)
		if err != nil {
			return nil, err
		}
		if test != nil {
			tests = append(tests, test)
		}
	}

	return tests, rows.Err()
}

// GetAvailableTests retrieves tests available for a user (active and within time window)
func (r *TestRepository) GetAvailableTests(userID int, limit, offset int) ([]*models.Test, error) {
	now := time.Now()
	query := `
		SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
		FROM tests 
		WHERE is_active = true 
		AND (start_time IS NULL OR start_time <= ?)
		AND (end_time IS NULL OR end_time >= ?)
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, title, description, created_by, duration_minutes, total_marks, passing_marks, instructions, is_active, start_time, end_time, created_at, updated_at
			FROM tests 
			WHERE is_active = true 
			AND (start_time IS NULL OR start_time <= $1)
			AND (end_time IS NULL OR end_time >= $2)
			ORDER BY created_at DESC LIMIT $3 OFFSET $4
		`
	}

	rows, err := r.db.Query(query, now, now, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []*models.Test
	for rows.Next() {
		test, err := models.ScanTest(rows)
		if err != nil {
			return nil, err
		}
		if test != nil {
			tests = append(tests, test)
		}
	}

	return tests, rows.Err()
}
