package database

import (
	"gocbt/internal/models"
	"time"
)

// QuestionRepository implements the models.QuestionRepository interface
type QuestionRepository struct {
	db *DB
}

// NewQuestionRepository creates a new question repository
func NewQuestionRepository(db *DB) models.QuestionRepository {
	return &QuestionRepository{db: db}
}

// Create creates a new question
func (r *QuestionRepository) Create(question *models.Question) error {
	query := `
		INSERT INTO questions (test_id, question_text, question_type, marks, order_index)
		VALUES (?, ?, ?, ?, ?)
	`

	if r.db.Driver == "postgres" {
		query = `
			INSERT INTO questions (test_id, question_text, question_type, marks, order_index)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at, updated_at
		`
	}

	if r.db.Driver == "postgres" {
		err := r.db.QueryRow(query, question.TestID, question.QuestionText,
			question.QuestionType, question.Marks, question.OrderIndex).Scan(
			&question.ID, &question.CreatedAt, &question.UpdatedAt)
		return err
	}

	result, err := r.db.Exec(query, question.TestID, question.QuestionText,
		question.QuestionType, question.Marks, question.OrderIndex)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	question.ID = int(id)
	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()
	return nil
}

// GetByID retrieves a question by ID
func (r *QuestionRepository) GetByID(id int) (*models.Question, error) {
	query := `
		SELECT id, test_id, question_text, question_type, marks, order_index, created_at, updated_at
		FROM questions WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, test_id, question_text, question_type, marks, order_index, created_at, updated_at
			FROM questions WHERE id = $1
		`
	}

	row := r.db.QueryRow(query, id)
	return models.ScanQuestion(row)
}

// GetByTestID retrieves questions by test ID
func (r *QuestionRepository) GetByTestID(testID int) ([]*models.Question, error) {
	query := `
		SELECT id, test_id, question_text, question_type, marks, order_index, created_at, updated_at
		FROM questions WHERE test_id = ? ORDER BY order_index ASC
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, test_id, question_text, question_type, marks, order_index, created_at, updated_at
			FROM questions WHERE test_id = $1 ORDER BY order_index ASC
		`
	}

	rows, err := r.db.Query(query, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []*models.Question
	for rows.Next() {
		question, err := models.ScanQuestion(rows)
		if err != nil {
			return nil, err
		}
		if question != nil {
			questions = append(questions, question)
		}
	}

	return questions, rows.Err()
}

// Update updates a question
func (r *QuestionRepository) Update(question *models.Question) error {
	query := `
		UPDATE questions 
		SET question_text = ?, question_type = ?, marks = ?, order_index = ?, updated_at = ?
		WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE questions 
			SET question_text = $1, question_type = $2, marks = $3, order_index = $4, updated_at = $5
			WHERE id = $6
		`
	}

	question.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, question.QuestionText, question.QuestionType,
		question.Marks, question.OrderIndex, question.UpdatedAt, question.ID)
	return err
}

// Delete deletes a question
func (r *QuestionRepository) Delete(id int) error {
	query := "DELETE FROM questions WHERE id = ?"
	if r.db.Driver == "postgres" {
		query = "DELETE FROM questions WHERE id = $1"
	}

	_, err := r.db.Exec(query, id)
	return err
}

// CreateOption creates a new question option
func (r *QuestionRepository) CreateOption(option *models.QuestionOption) error {
	query := `
		INSERT INTO question_options (question_id, option_text, is_correct, order_index)
		VALUES (?, ?, ?, ?)
	`

	if r.db.Driver == "postgres" {
		query = `
			INSERT INTO question_options (question_id, option_text, is_correct, order_index)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at
		`
	}

	if r.db.Driver == "postgres" {
		err := r.db.QueryRow(query, option.QuestionID, option.OptionText,
			option.IsCorrect, option.OrderIndex).Scan(&option.ID, &option.CreatedAt)
		return err
	}

	result, err := r.db.Exec(query, option.QuestionID, option.OptionText,
		option.IsCorrect, option.OrderIndex)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	option.ID = int(id)
	option.CreatedAt = time.Now()
	return nil
}

// GetOptionsByQuestionID retrieves options by question ID
func (r *QuestionRepository) GetOptionsByQuestionID(questionID int) ([]*models.QuestionOption, error) {
	query := `
		SELECT id, question_id, option_text, is_correct, order_index, created_at
		FROM question_options WHERE question_id = ? ORDER BY order_index ASC
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, question_id, option_text, is_correct, order_index, created_at
			FROM question_options WHERE question_id = $1 ORDER BY order_index ASC
		`
	}

	rows, err := r.db.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var options []*models.QuestionOption
	for rows.Next() {
		option, err := models.ScanQuestionOption(rows)
		if err != nil {
			return nil, err
		}
		if option != nil {
			options = append(options, option)
		}
	}

	return options, rows.Err()
}

// UpdateOption updates a question option
func (r *QuestionRepository) UpdateOption(option *models.QuestionOption) error {
	query := `
		UPDATE question_options 
		SET option_text = ?, is_correct = ?, order_index = ?
		WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE question_options 
			SET option_text = $1, is_correct = $2, order_index = $3
			WHERE id = $4
		`
	}

	_, err := r.db.Exec(query, option.OptionText, option.IsCorrect,
		option.OrderIndex, option.ID)
	return err
}

// DeleteOption deletes a question option
func (r *QuestionRepository) DeleteOption(id int) error {
	query := "DELETE FROM question_options WHERE id = ?"
	if r.db.Driver == "postgres" {
		query = "DELETE FROM question_options WHERE id = $1"
	}

	_, err := r.db.Exec(query, id)
	return err
}

// CreateCorrectAnswer creates a new correct answer
func (r *QuestionRepository) CreateCorrectAnswer(answer *models.CorrectAnswer) error {
	query := `
		INSERT INTO correct_answers (question_id, answer_text, is_case_sensitive)
		VALUES (?, ?, ?)
	`

	if r.db.Driver == "postgres" {
		query = `
			INSERT INTO correct_answers (question_id, answer_text, is_case_sensitive)
			VALUES ($1, $2, $3)
			RETURNING id, created_at
		`
	}

	if r.db.Driver == "postgres" {
		err := r.db.QueryRow(query, answer.QuestionID, answer.AnswerText,
			answer.IsCaseSensitive).Scan(&answer.ID, &answer.CreatedAt)
		return err
	}

	result, err := r.db.Exec(query, answer.QuestionID, answer.AnswerText,
		answer.IsCaseSensitive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	answer.ID = int(id)
	answer.CreatedAt = time.Now()
	return nil
}

// GetCorrectAnswersByQuestionID retrieves correct answers by question ID
func (r *QuestionRepository) GetCorrectAnswersByQuestionID(questionID int) ([]*models.CorrectAnswer, error) {
	query := `
		SELECT id, question_id, answer_text, is_case_sensitive, created_at
		FROM correct_answers WHERE question_id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			SELECT id, question_id, answer_text, is_case_sensitive, created_at
			FROM correct_answers WHERE question_id = $1
		`
	}

	rows, err := r.db.Query(query, questionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []*models.CorrectAnswer
	for rows.Next() {
		answer, err := models.ScanCorrectAnswer(rows)
		if err != nil {
			return nil, err
		}
		if answer != nil {
			answers = append(answers, answer)
		}
	}

	return answers, rows.Err()
}

// UpdateCorrectAnswer updates a correct answer
func (r *QuestionRepository) UpdateCorrectAnswer(answer *models.CorrectAnswer) error {
	query := `
		UPDATE correct_answers 
		SET answer_text = ?, is_case_sensitive = ?
		WHERE id = ?
	`

	if r.db.Driver == "postgres" {
		query = `
			UPDATE correct_answers 
			SET answer_text = $1, is_case_sensitive = $2
			WHERE id = $3
		`
	}

	_, err := r.db.Exec(query, answer.AnswerText, answer.IsCaseSensitive, answer.ID)
	return err
}

// DeleteCorrectAnswer deletes a correct answer
func (r *QuestionRepository) DeleteCorrectAnswer(id int) error {
	query := "DELETE FROM correct_answers WHERE id = ?"
	if r.db.Driver == "postgres" {
		query = "DELETE FROM correct_answers WHERE id = $1"
	}

	_, err := r.db.Exec(query, id)
	return err
}
