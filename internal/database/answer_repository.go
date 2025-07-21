package database

import (
	"gocbt/internal/models"
	"time"
)

// UserAnswerRepository implements the models.UserAnswerRepository interface
type UserAnswerRepository struct {
	db *DB
}

// NewUserAnswerRepository creates a new user answer repository
func NewUserAnswerRepository(db *DB) models.UserAnswerRepository {
	return &UserAnswerRepository{db: db}
}

// Create creates a new user answer
func (r *UserAnswerRepository) Create(answer *models.UserAnswer) error {
	query := `
		INSERT INTO user_answers (session_id, question_id, answer_text, selected_option_id, is_correct, marks_awarded)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	if r.db.Driver == "postgres" {
		query = `
			INSERT INTO user_answers (session_id, question_id, answer_text, selected_option_id, is_correct, marks_awarded)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, answered_at
		`
	}

	if r.db.Driver == "postgres" {
		err := r.db.QueryRow(query, answer.SessionID, answer.QuestionID, answer.AnswerText,
			answer.SelectedOptionID, answer.IsCorrect, answer.MarksAwarded).Scan(
			&answer.ID, &answer.AnsweredAt)
		return err
	}

	result, err := r.db.Exec(query, answer.SessionID, answer.QuestionID, answer.AnswerText,
		answer.SelectedOptionID, answer.IsCorrect, answer.MarksAwarded)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	answer.ID = int(id)
	answer.AnsweredAt = time.Now()
	return nil
}

// GetByID retrieves a user answer by ID
func (r *UserAnswerRepository) GetByID(id int) (*models.UserAnswer, error) {
	query := `
		SELECT id, session_id, question_id, answer_text, selected_option_id, is_correct, marks_awarded, answered_at
		FROM user_answers WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, session_id, question_id, answer_text, selected_option_id, is_correct, marks_awarded, answered_at
			FROM user_answers WHERE id = $1
		`
	}

	row := r.db.QueryRow(query, id)
	return models.ScanUserAnswer(row)
}

// GetBySessionAndQuestion retrieves a user answer by session and question
func (r *UserAnswerRepository) GetBySessionAndQuestion(sessionID, questionID int) (*models.UserAnswer, error) {
	query := `
		SELECT id, session_id, question_id, answer_text, selected_option_id, is_correct, marks_awarded, answered_at
		FROM user_answers WHERE session_id = ? AND question_id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, session_id, question_id, answer_text, selected_option_id, is_correct, marks_awarded, answered_at
			FROM user_answers WHERE session_id = $1 AND question_id = $2
		`
	}

	row := r.db.QueryRow(query, sessionID, questionID)
	return models.ScanUserAnswer(row)
}

// GetBySession retrieves all user answers for a session
func (r *UserAnswerRepository) GetBySession(sessionID int) ([]*models.UserAnswer, error) {
	query := `
		SELECT id, session_id, question_id, answer_text, selected_option_id, is_correct, marks_awarded, answered_at
		FROM user_answers WHERE session_id = ? ORDER BY answered_at ASC
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, session_id, question_id, answer_text, selected_option_id, is_correct, marks_awarded, answered_at
			FROM user_answers WHERE session_id = $1 ORDER BY answered_at ASC
		`
	}

	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*models.UserAnswer
	for rows.Next() {
		answer, err := models.ScanUserAnswer(rows)
		if err != nil {
			return nil, err
		}
		if answer != nil {
			answers = append(answers, answer)
		}
	}

	return answers, rows.Err()
}

// Update updates a user answer
func (r *UserAnswerRepository) Update(answer *models.UserAnswer) error {
	query := `
		UPDATE user_answers 
		SET answer_text = ?, selected_option_id = ?, is_correct = ?, marks_awarded = ?
		WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE user_answers 
			SET answer_text = $1, selected_option_id = $2, is_correct = $3, marks_awarded = $4
			WHERE id = $5
		`
	}

	_, err := r.db.Exec(query, answer.AnswerText, answer.SelectedOptionID,
		answer.IsCorrect, answer.MarksAwarded, answer.ID)
	return err
}

// Delete deletes a user answer
func (r *UserAnswerRepository) Delete(id int) error {
	query := "DELETE FROM user_answers WHERE id = ?"
	if r.db.Driver == "postgres" {
		query = "DELETE FROM user_answers WHERE id = $1"
	}

	_, err := r.db.Exec(query, id)
	return err
}
