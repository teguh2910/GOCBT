package models

import (
	"database/sql"
	"time"
)

// QuestionType represents the type of question
type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeTrueFalse      QuestionType = "true_false"
	QuestionTypeShortAnswer    QuestionType = "short_answer"
)

// Question represents a question in a test
type Question struct {
	ID           int          `json:"id" db:"id"`
	TestID       int          `json:"test_id" db:"test_id"`
	QuestionText string       `json:"question_text" db:"question_text"`
	QuestionType QuestionType `json:"question_type" db:"question_type"`
	Marks        int          `json:"marks" db:"marks"`
	OrderIndex   int          `json:"order_index" db:"order_index"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	
	// Related data (not stored in database)
	Options        []*QuestionOption `json:"options,omitempty"`
	CorrectAnswers []*CorrectAnswer  `json:"correct_answers,omitempty"`
}

// QuestionOption represents an option for multiple choice questions
type QuestionOption struct {
	ID         int       `json:"id" db:"id"`
	QuestionID int       `json:"question_id" db:"question_id"`
	OptionText string    `json:"option_text" db:"option_text"`
	IsCorrect  bool      `json:"is_correct" db:"is_correct"`
	OrderIndex int       `json:"order_index" db:"order_index"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

// CorrectAnswer represents correct answers for non-multiple choice questions
type CorrectAnswer struct {
	ID              int       `json:"id" db:"id"`
	QuestionID      int       `json:"question_id" db:"question_id"`
	AnswerText      string    `json:"answer_text" db:"answer_text"`
	IsCaseSensitive bool      `json:"is_case_sensitive" db:"is_case_sensitive"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// QuestionRepository defines the interface for question data operations
type QuestionRepository interface {
	Create(question *Question) error
	GetByID(id int) (*Question, error)
	GetByTestID(testID int) ([]*Question, error)
	Update(question *Question) error
	Delete(id int) error
	CreateOption(option *QuestionOption) error
	GetOptionsByQuestionID(questionID int) ([]*QuestionOption, error)
	UpdateOption(option *QuestionOption) error
	DeleteOption(id int) error
	CreateCorrectAnswer(answer *CorrectAnswer) error
	GetCorrectAnswersByQuestionID(questionID int) ([]*CorrectAnswer, error)
	UpdateCorrectAnswer(answer *CorrectAnswer) error
	DeleteCorrectAnswer(id int) error
}

// QuestionService defines the interface for question business logic
type QuestionService interface {
	CreateQuestion(testID int, questionText string, questionType QuestionType, marks, orderIndex int) (*Question, error)
	GetQuestion(questionID int) (*Question, error)
	GetTestQuestions(testID int) ([]*Question, error)
	UpdateQuestion(questionID int, questionText string, marks, orderIndex int) (*Question, error)
	DeleteQuestion(questionID int) error
	AddOption(questionID int, optionText string, isCorrect bool, orderIndex int) (*QuestionOption, error)
	UpdateOption(optionID int, optionText string, isCorrect bool, orderIndex int) (*QuestionOption, error)
	DeleteOption(optionID int) error
	AddCorrectAnswer(questionID int, answerText string, isCaseSensitive bool) (*CorrectAnswer, error)
	UpdateCorrectAnswer(answerID int, answerText string, isCaseSensitive bool) (*CorrectAnswer, error)
	DeleteCorrectAnswer(answerID int) error
}

// IsValidType checks if the question type is valid
func (qt QuestionType) IsValid() bool {
	switch qt {
	case QuestionTypeMultipleChoice, QuestionTypeTrueFalse, QuestionTypeShortAnswer:
		return true
	default:
		return false
	}
}

// ScanQuestion scans database row into Question struct
func ScanQuestion(row interface {
	Scan(dest ...interface{}) error
}) (*Question, error) {
	question := &Question{}
	err := row.Scan(
		&question.ID,
		&question.TestID,
		&question.QuestionText,
		&question.QuestionType,
		&question.Marks,
		&question.OrderIndex,
		&question.CreatedAt,
		&question.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return question, nil
}

// ScanQuestionOption scans database row into QuestionOption struct
func ScanQuestionOption(row interface {
	Scan(dest ...interface{}) error
}) (*QuestionOption, error) {
	option := &QuestionOption{}
	err := row.Scan(
		&option.ID,
		&option.QuestionID,
		&option.OptionText,
		&option.IsCorrect,
		&option.OrderIndex,
		&option.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return option, nil
}

// ScanCorrectAnswer scans database row into CorrectAnswer struct
func ScanCorrectAnswer(row interface {
	Scan(dest ...interface{}) error
}) (*CorrectAnswer, error) {
	answer := &CorrectAnswer{}
	err := row.Scan(
		&answer.ID,
		&answer.QuestionID,
		&answer.AnswerText,
		&answer.IsCaseSensitive,
		&answer.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return answer, nil
}
