package database

import (
	"gocbt/internal/models"
	"time"
)

// TestSessionRepository implements the models.TestSessionRepository interface
type TestSessionRepository struct {
	db *DB
}

// NewTestSessionRepository creates a new test session repository
func NewTestSessionRepository(db *DB) models.TestSessionRepository {
	return &TestSessionRepository{db: db}
}

// Create creates a new test session
func (r *TestSessionRepository) Create(session *models.TestSession) error {
	query := `
		INSERT INTO test_sessions (test_id, user_id, session_token, status, expires_at, time_remaining, current_question_index)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	if r.db.Driver == "postgres" {
		query = `
			INSERT INTO test_sessions (test_id, user_id, session_token, status, expires_at, time_remaining, current_question_index)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, created_at, updated_at
		`
	}

	if r.db.Driver == "postgres" {
		err := r.db.QueryRow(query, session.TestID, session.UserID, session.SessionToken,
			session.Status, session.ExpiresAt, session.TimeRemaining, session.CurrentQuestionIndex).Scan(
			&session.ID, &session.CreatedAt, &session.UpdatedAt)
		return err
	}

	result, err := r.db.Exec(query, session.TestID, session.UserID, session.SessionToken,
		session.Status, session.ExpiresAt, session.TimeRemaining, session.CurrentQuestionIndex)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	session.ID = int(id)
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()
	return nil
}

// GetByID retrieves a test session by ID
func (r *TestSessionRepository) GetByID(id int) (*models.TestSession, error) {
	query := `
		SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
		FROM test_sessions WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
			FROM test_sessions WHERE id = $1
		`
	}

	row := r.db.QueryRow(query, id)
	return models.ScanTestSession(row)
}

// GetByToken retrieves a test session by token
func (r *TestSessionRepository) GetByToken(token string) (*models.TestSession, error) {
	query := `
		SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
		FROM test_sessions WHERE session_token = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
			FROM test_sessions WHERE session_token = $1
		`
	}

	row := r.db.QueryRow(query, token)
	return models.ScanTestSession(row)
}

// GetByUserAndTest retrieves a test session by user and test
func (r *TestSessionRepository) GetByUserAndTest(userID, testID int) (*models.TestSession, error) {
	query := `
		SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
		FROM test_sessions WHERE user_id = ? AND test_id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
			FROM test_sessions WHERE user_id = $1 AND test_id = $2
		`
	}

	row := r.db.QueryRow(query, userID, testID)
	return models.ScanTestSession(row)
}

// Update updates a test session
func (r *TestSessionRepository) Update(session *models.TestSession) error {
	query := `
		UPDATE test_sessions 
		SET status = ?, started_at = ?, submitted_at = ?, time_remaining = ?, current_question_index = ?, updated_at = ?
		WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE test_sessions 
			SET status = $1, started_at = $2, submitted_at = $3, time_remaining = $4, current_question_index = $5, updated_at = $6
			WHERE id = $7
		`
	}

	session.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, session.Status, session.StartedAt, session.SubmittedAt,
		session.TimeRemaining, session.CurrentQuestionIndex, session.UpdatedAt, session.ID)
	return err
}

// Delete deletes a test session
func (r *TestSessionRepository) Delete(id int) error {
	query := "DELETE FROM test_sessions WHERE id = ?"
	if r.db.Driver == "postgres" {
		query = "DELETE FROM test_sessions WHERE id = $1"
	}

	_, err := r.db.Exec(query, id)
	return err
}

// GetActiveSessionsByTest retrieves active sessions for a test
func (r *TestSessionRepository) GetActiveSessionsByTest(testID int) ([]*models.TestSession, error) {
	query := `
		SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
		FROM test_sessions 
		WHERE test_id = ? AND status IN ('not_started', 'in_progress')
		ORDER BY created_at DESC
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
			FROM test_sessions 
			WHERE test_id = $1 AND status IN ('not_started', 'in_progress')
			ORDER BY created_at DESC
		`
	}

	rows, err := r.db.Query(query, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.TestSession
	for rows.Next() {
		session, err := models.ScanTestSession(rows)
		if err != nil {
			return nil, err
		}
		if session != nil {
			sessions = append(sessions, session)
		}
	}

	return sessions, rows.Err()
}

// GetUserSessions retrieves sessions for a user with pagination
func (r *TestSessionRepository) GetUserSessions(userID int, limit, offset int) ([]*models.TestSession, error) {
	query := `
		SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
		FROM test_sessions 
		WHERE user_id = ? 
		ORDER BY created_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, test_id, user_id, session_token, status, started_at, submitted_at, expires_at, time_remaining, current_question_index, created_at, updated_at
			FROM test_sessions 
			WHERE user_id = $1 
			ORDER BY created_at DESC LIMIT $2 OFFSET $3
		`
	}

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*models.TestSession
	for rows.Next() {
		session, err := models.ScanTestSession(rows)
		if err != nil {
			return nil, err
		}
		if session != nil {
			sessions = append(sessions, session)
		}
	}

	return sessions, rows.Err()
}

// ExpireOldSessions marks expired sessions as expired
func (r *TestSessionRepository) ExpireOldSessions() error {
	now := time.Now()
	query := `
		UPDATE test_sessions 
		SET status = 'expired', updated_at = ?
		WHERE expires_at < ? AND status IN ('not_started', 'in_progress')
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE test_sessions 
			SET status = 'expired', updated_at = $1
			WHERE expires_at < $2 AND status IN ('not_started', 'in_progress')
		`
	}

	_, err := r.db.Exec(query, now, now)
	return err
}
