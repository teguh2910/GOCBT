package database

import (
	"gocbt/internal/models"
	"time"
)

// TestResultRepository implements the models.TestResultRepository interface
type TestResultRepository struct {
	db *DB
}

// NewTestResultRepository creates a new test result repository
func NewTestResultRepository(db *DB) models.TestResultRepository {
	return &TestResultRepository{db: db}
}

// Create creates a new test result
func (r *TestResultRepository) Create(result *models.TestResult) error {
	query := `
		INSERT INTO test_results (session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	if r.db.Driver == "postgres" {
		query = `
			INSERT INTO test_results (session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			RETURNING id, completed_at
		`
	}

	if r.db.Driver == "postgres" {
		err := r.db.QueryRow(query, result.SessionID, result.TestID, result.UserID,
			result.TotalQuestions, result.AnsweredQuestions, result.CorrectAnswers,
			result.TotalMarks, result.MarksObtained, result.Percentage, result.Grade,
			result.IsPassed, result.TimeTaken).Scan(&result.ID, &result.CompletedAt)
		return err
	}

	res, err := r.db.Exec(query, result.SessionID, result.TestID, result.UserID,
		result.TotalQuestions, result.AnsweredQuestions, result.CorrectAnswers,
		result.TotalMarks, result.MarksObtained, result.Percentage, result.Grade,
		result.IsPassed, result.TimeTaken)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	result.ID = int(id)
	result.CompletedAt = time.Now()
	return nil
}

// GetByID retrieves a test result by ID
func (r *TestResultRepository) GetByID(id int) (*models.TestResult, error) {
	query := `
		SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
		FROM test_results WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
			FROM test_results WHERE id = $1
		`
	}

	row := r.db.QueryRow(query, id)
	return models.ScanTestResult(row)
}

// GetBySessionID retrieves a test result by session ID
func (r *TestResultRepository) GetBySessionID(sessionID int) (*models.TestResult, error) {
	query := `
		SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
		FROM test_results WHERE session_id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
			FROM test_results WHERE session_id = $1
		`
	}

	row := r.db.QueryRow(query, sessionID)
	return models.ScanTestResult(row)
}

// GetByUserAndTest retrieves a test result by user and test
func (r *TestResultRepository) GetByUserAndTest(userID, testID int) (*models.TestResult, error) {
	query := `
		SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
		FROM test_results WHERE user_id = ? AND test_id = ? ORDER BY completed_at DESC LIMIT 1
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
			FROM test_results WHERE user_id = $1 AND test_id = $2 ORDER BY completed_at DESC LIMIT 1
		`
	}

	row := r.db.QueryRow(query, userID, testID)
	return models.ScanTestResult(row)
}

// GetByUser retrieves test results by user with pagination
func (r *TestResultRepository) GetByUser(userID int, limit, offset int) ([]*models.TestResult, error) {
	query := `
		SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
		FROM test_results WHERE user_id = ? ORDER BY completed_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
			FROM test_results WHERE user_id = $1 ORDER BY completed_at DESC LIMIT $2 OFFSET $3
		`
	}

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.TestResult
	for rows.Next() {
		result, err := models.ScanTestResult(rows)
		if err != nil {
			return nil, err
		}
		if result != nil {
			results = append(results, result)
		}
	}

	return results, rows.Err()
}

// GetByTest retrieves test results by test with pagination
func (r *TestResultRepository) GetByTest(testID int, limit, offset int) ([]*models.TestResult, error) {
	query := `
		SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
		FROM test_results WHERE test_id = ? ORDER BY completed_at DESC LIMIT ? OFFSET ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, session_id, test_id, user_id, total_questions, answered_questions, correct_answers, total_marks, marks_obtained, percentage, grade, is_passed, time_taken, completed_at
			FROM test_results WHERE test_id = $1 ORDER BY completed_at DESC LIMIT $2 OFFSET $3
		`
	}

	rows, err := r.db.Query(query, testID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*models.TestResult
	for rows.Next() {
		result, err := models.ScanTestResult(rows)
		if err != nil {
			return nil, err
		}
		if result != nil {
			results = append(results, result)
		}
	}

	return results, rows.Err()
}

// Update updates a test result
func (r *TestResultRepository) Update(result *models.TestResult) error {
	query := `
		UPDATE test_results 
		SET total_questions = ?, answered_questions = ?, correct_answers = ?, total_marks = ?, marks_obtained = ?, percentage = ?, grade = ?, is_passed = ?, time_taken = ?
		WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE test_results 
			SET total_questions = $1, answered_questions = $2, correct_answers = $3, total_marks = $4, marks_obtained = $5, percentage = $6, grade = $7, is_passed = $8, time_taken = $9
			WHERE id = $10
		`
	}

	_, err := r.db.Exec(query, result.TotalQuestions, result.AnsweredQuestions,
		result.CorrectAnswers, result.TotalMarks, result.MarksObtained, result.Percentage,
		result.Grade, result.IsPassed, result.TimeTaken, result.ID)
	return err
}

// Delete deletes a test result
func (r *TestResultRepository) Delete(id int) error {
	query := "DELETE FROM test_results WHERE id = ?"
	if r.db.Driver == "postgres" {
		query = "DELETE FROM test_results WHERE id = $1"
	}

	_, err := r.db.Exec(query, id)
	return err
}

// GetTestStatistics retrieves statistics for a test
func (r *TestResultRepository) GetTestStatistics(testID int) (*models.TestStatistics, error) {
	query := `
		SELECT
			COUNT(*) as total_attempts,
			COUNT(CASE WHEN percentage >= 0 THEN 1 END) as completed_attempts,
			COUNT(CASE WHEN is_passed = true THEN 1 END) as passed_attempts,
			COALESCE(AVG(percentage), 0) as average_score,
			COALESCE(MAX(percentage), 0) as highest_score,
			COALESCE(MIN(percentage), 0) as lowest_score,
			COALESCE(ROUND(AVG(time_taken)), 0) as average_time_taken
		FROM test_results
		WHERE test_id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT
				COUNT(*) as total_attempts,
				COUNT(CASE WHEN percentage >= 0 THEN 1 END) as completed_attempts,
				COUNT(CASE WHEN is_passed = true THEN 1 END) as passed_attempts,
				COALESCE(AVG(percentage), 0) as average_score,
				COALESCE(MAX(percentage), 0) as highest_score,
				COALESCE(MIN(percentage), 0) as lowest_score,
				COALESCE(ROUND(AVG(time_taken)), 0) as average_time_taken
			FROM test_results
			WHERE test_id = $1
		`
	}

	stats := &models.TestStatistics{TestID: testID}
	err := r.db.QueryRow(query, testID).Scan(
		&stats.TotalAttempts,
		&stats.CompletedAttempts,
		&stats.PassedAttempts,
		&stats.AverageScore,
		&stats.HighestScore,
		&stats.LowestScore,
		&stats.AverageTimeTaken,
	)

	if err != nil {
		return nil, err
	}

	return stats, nil
}
