package models

import (
	"database/sql"
	"fmt"
	"time"
)

// TestResult represents the result of a completed test
type TestResult struct {
	ID                int       `json:"id" db:"id"`
	SessionID         int       `json:"session_id" db:"session_id"`
	TestID            int       `json:"test_id" db:"test_id"`
	UserID            int       `json:"user_id" db:"user_id"`
	TotalQuestions    int       `json:"total_questions" db:"total_questions"`
	AnsweredQuestions int       `json:"answered_questions" db:"answered_questions"`
	CorrectAnswers    int       `json:"correct_answers" db:"correct_answers"`
	TotalMarks        int       `json:"total_marks" db:"total_marks"`
	MarksObtained     int       `json:"marks_obtained" db:"marks_obtained"`
	Percentage        float64   `json:"percentage" db:"percentage"`
	Grade             *string   `json:"grade" db:"grade"`
	IsPassed          bool      `json:"is_passed" db:"is_passed"`
	TimeTaken         *int      `json:"time_taken" db:"time_taken"` // in seconds
	CompletedAt       time.Time `json:"completed_at" db:"completed_at"`

	// Related data (not stored in database)
	Test    *Test        `json:"test,omitempty"`
	User    *User        `json:"user,omitempty"`
	Session *TestSession `json:"session,omitempty"`
}

// TestResultRepository defines the interface for test result data operations
type TestResultRepository interface {
	Create(result *TestResult) error
	GetByID(id int) (*TestResult, error)
	GetBySessionID(sessionID int) (*TestResult, error)
	GetByUserAndTest(userID, testID int) (*TestResult, error)
	GetByUser(userID int, limit, offset int) ([]*TestResult, error)
	GetByTest(testID int, limit, offset int) ([]*TestResult, error)
	Update(result *TestResult) error
	Delete(id int) error
	GetTestStatistics(testID int) (*TestStatistics, error)
}

// TestResultService defines the interface for test result business logic
type TestResultService interface {
	CalculateResult(sessionID int) (*TestResult, error)
	GetResult(resultID int) (*TestResult, error)
	GetUserResults(userID int, limit, offset int) ([]*TestResult, error)
	GetTestResults(testID int, limit, offset int) ([]*TestResult, error)
	GetTestStatistics(testID int) (*TestStatistics, error)
	GetResultBySession(sessionID int) (*TestResult, error)
	GetResultByUserAndTest(userID, testID int) (*TestResult, error)
}

// TestStatistics represents statistics for a test
type TestStatistics struct {
	TestID            int     `json:"test_id"`
	TotalAttempts     int     `json:"total_attempts"`
	CompletedAttempts int     `json:"completed_attempts"`
	PassedAttempts    int     `json:"passed_attempts"`
	AverageScore      float64 `json:"average_score"`
	HighestScore      float64 `json:"highest_score"`
	LowestScore       float64 `json:"lowest_score"`
	AverageTimeTaken  int     `json:"average_time_taken"` // in seconds
}

// CalculateGrade calculates the grade based on percentage
func (r *TestResult) CalculateGrade() string {
	switch {
	case r.Percentage >= 90:
		return "A+"
	case r.Percentage >= 80:
		return "A"
	case r.Percentage >= 70:
		return "B+"
	case r.Percentage >= 60:
		return "B"
	case r.Percentage >= 50:
		return "C"
	case r.Percentage >= 40:
		return "D"
	default:
		return "F"
	}
}

// GetTimeTakenFormatted returns formatted time taken (e.g., "1h 30m 45s")
func (r *TestResult) GetTimeTakenFormatted() string {
	if r.TimeTaken == nil {
		return "N/A"
	}

	duration := time.Duration(*r.TimeTaken) * time.Second
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	} else {
		return fmt.Sprintf("%ds", seconds)
	}
}

// ScanTestResult scans database row into TestResult struct
func ScanTestResult(row interface {
	Scan(dest ...interface{}) error
}) (*TestResult, error) {
	result := &TestResult{}
	err := row.Scan(
		&result.ID,
		&result.SessionID,
		&result.TestID,
		&result.UserID,
		&result.TotalQuestions,
		&result.AnsweredQuestions,
		&result.CorrectAnswers,
		&result.TotalMarks,
		&result.MarksObtained,
		&result.Percentage,
		&result.Grade,
		&result.IsPassed,
		&result.TimeTaken,
		&result.CompletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return result, nil
}
