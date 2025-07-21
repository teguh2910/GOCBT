package models

import (
	"database/sql"
	"time"
)

// SessionStatus represents the status of a test session
type SessionStatus string

const (
	SessionStatusNotStarted SessionStatus = "not_started"
	SessionStatusInProgress SessionStatus = "in_progress"
	SessionStatusCompleted  SessionStatus = "completed"
	SessionStatusSubmitted  SessionStatus = "submitted"
	SessionStatusExpired    SessionStatus = "expired"
)

// TestSession represents a user's test session
type TestSession struct {
	ID                   int           `json:"id" db:"id"`
	TestID               int           `json:"test_id" db:"test_id"`
	UserID               int           `json:"user_id" db:"user_id"`
	SessionToken         string        `json:"session_token" db:"session_token"`
	Status               SessionStatus `json:"status" db:"status"`
	StartedAt            *time.Time    `json:"started_at" db:"started_at"`
	SubmittedAt          *time.Time    `json:"submitted_at" db:"submitted_at"`
	ExpiresAt            time.Time     `json:"expires_at" db:"expires_at"`
	TimeRemaining        *int          `json:"time_remaining" db:"time_remaining"` // in seconds
	CurrentQuestionIndex int           `json:"current_question_index" db:"current_question_index"`
	CreatedAt            time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at" db:"updated_at"`
	
	// Related data (not stored in database)
	Test    *Test        `json:"test,omitempty"`
	User    *User        `json:"user,omitempty"`
	Answers []*UserAnswer `json:"answers,omitempty"`
}

// UserAnswer represents a user's answer to a question
type UserAnswer struct {
	ID               int       `json:"id" db:"id"`
	SessionID        int       `json:"session_id" db:"session_id"`
	QuestionID       int       `json:"question_id" db:"question_id"`
	AnswerText       *string   `json:"answer_text" db:"answer_text"`
	SelectedOptionID *int      `json:"selected_option_id" db:"selected_option_id"`
	IsCorrect        *bool     `json:"is_correct" db:"is_correct"`
	MarksAwarded     int       `json:"marks_awarded" db:"marks_awarded"`
	AnsweredAt       time.Time `json:"answered_at" db:"answered_at"`
	
	// Related data (not stored in database)
	Question       *Question       `json:"question,omitempty"`
	SelectedOption *QuestionOption `json:"selected_option,omitempty"`
}

// TestSessionRepository defines the interface for test session data operations
type TestSessionRepository interface {
	Create(session *TestSession) error
	GetByID(id int) (*TestSession, error)
	GetByToken(token string) (*TestSession, error)
	GetByUserAndTest(userID, testID int) (*TestSession, error)
	Update(session *TestSession) error
	Delete(id int) error
	GetActiveSessionsByTest(testID int) ([]*TestSession, error)
	GetUserSessions(userID int, limit, offset int) ([]*TestSession, error)
	ExpireOldSessions() error
}

// UserAnswerRepository defines the interface for user answer data operations
type UserAnswerRepository interface {
	Create(answer *UserAnswer) error
	GetByID(id int) (*UserAnswer, error)
	GetBySessionAndQuestion(sessionID, questionID int) (*UserAnswer, error)
	GetBySession(sessionID int) ([]*UserAnswer, error)
	Update(answer *UserAnswer) error
	Delete(id int) error
}

// TestSessionService defines the interface for test session business logic
type TestSessionService interface {
	StartSession(userID, testID int) (*TestSession, error)
	GetSession(sessionToken string) (*TestSession, error)
	SubmitAnswer(sessionToken string, questionID int, answerText *string, selectedOptionID *int) (*UserAnswer, error)
	GetSessionAnswers(sessionToken string) ([]*UserAnswer, error)
	SubmitSession(sessionToken string) (*TestSession, error)
	GetUserSessions(userID int, limit, offset int) ([]*TestSession, error)
	UpdateSessionProgress(sessionToken string, currentQuestionIndex int) error
}

// IsExpired checks if the session has expired
func (s *TestSession) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsActive checks if the session is active (in progress and not expired)
func (s *TestSession) IsActive() bool {
	return s.Status == SessionStatusInProgress && !s.IsExpired()
}

// GetRemainingTime returns the remaining time in seconds
func (s *TestSession) GetRemainingTime() int {
	if s.IsExpired() {
		return 0
	}
	
	if s.TimeRemaining != nil {
		return *s.TimeRemaining
	}
	
	remaining := int(time.Until(s.ExpiresAt).Seconds())
	if remaining < 0 {
		return 0
	}
	return remaining
}

// ScanTestSession scans database row into TestSession struct
func ScanTestSession(row interface {
	Scan(dest ...interface{}) error
}) (*TestSession, error) {
	session := &TestSession{}
	err := row.Scan(
		&session.ID,
		&session.TestID,
		&session.UserID,
		&session.SessionToken,
		&session.Status,
		&session.StartedAt,
		&session.SubmittedAt,
		&session.ExpiresAt,
		&session.TimeRemaining,
		&session.CurrentQuestionIndex,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return session, nil
}

// ScanUserAnswer scans database row into UserAnswer struct
func ScanUserAnswer(row interface {
	Scan(dest ...interface{}) error
}) (*UserAnswer, error) {
	answer := &UserAnswer{}
	err := row.Scan(
		&answer.ID,
		&answer.SessionID,
		&answer.QuestionID,
		&answer.AnswerText,
		&answer.SelectedOptionID,
		&answer.IsCorrect,
		&answer.MarksAwarded,
		&answer.AnsweredAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return answer, nil
}
