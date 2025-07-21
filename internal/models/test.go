package models

import (
	"database/sql"
	"time"
)

// Test represents a test in the system
type Test struct {
	ID             int       `json:"id" db:"id"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	CreatedBy      int       `json:"created_by" db:"created_by"`
	DurationMinutes int      `json:"duration_minutes" db:"duration_minutes"`
	TotalMarks     int       `json:"total_marks" db:"total_marks"`
	PassingMarks   int       `json:"passing_marks" db:"passing_marks"`
	Instructions   string    `json:"instructions" db:"instructions"`
	IsActive       bool      `json:"is_active" db:"is_active"`
	StartTime      *time.Time `json:"start_time" db:"start_time"`
	EndTime        *time.Time `json:"end_time" db:"end_time"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
	
	// Related data (not stored in database)
	Creator   *User       `json:"creator,omitempty"`
	Questions []*Question `json:"questions,omitempty"`
}

// TestRepository defines the interface for test data operations
type TestRepository interface {
	Create(test *Test) error
	GetByID(id int) (*Test, error)
	Update(test *Test) error
	Delete(id int) error
	List(limit, offset int) ([]*Test, error)
	GetByCreator(creatorID int, limit, offset int) ([]*Test, error)
	GetActiveTests(limit, offset int) ([]*Test, error)
	GetAvailableTests(userID int, limit, offset int) ([]*Test, error)
}

// TestService defines the interface for test business logic
type TestService interface {
	CreateTest(creatorID int, title, description, instructions string, durationMinutes, totalMarks, passingMarks int, startTime, endTime *time.Time) (*Test, error)
	GetTest(testID int) (*Test, error)
	UpdateTest(testID int, title, description, instructions string, durationMinutes, totalMarks, passingMarks int, startTime, endTime *time.Time) (*Test, error)
	DeleteTest(testID int) error
	ListTests(creatorID int, limit, offset int) ([]*Test, error)
	GetAvailableTests(userID int, limit, offset int) ([]*Test, error)
	ActivateTest(testID int) error
	DeactivateTest(testID int) error
}

// IsAvailable checks if the test is currently available for taking
func (t *Test) IsAvailable() bool {
	if !t.IsActive {
		return false
	}
	
	now := time.Now()
	
	// Check if test has started
	if t.StartTime != nil && now.Before(*t.StartTime) {
		return false
	}
	
	// Check if test has ended
	if t.EndTime != nil && now.After(*t.EndTime) {
		return false
	}
	
	return true
}

// GetDuration returns the test duration as a time.Duration
func (t *Test) GetDuration() time.Duration {
	return time.Duration(t.DurationMinutes) * time.Minute
}

// ScanTest scans database row into Test struct
func ScanTest(row interface {
	Scan(dest ...interface{}) error
}) (*Test, error) {
	test := &Test{}
	err := row.Scan(
		&test.ID,
		&test.Title,
		&test.Description,
		&test.CreatedBy,
		&test.DurationMinutes,
		&test.TotalMarks,
		&test.PassingMarks,
		&test.Instructions,
		&test.IsActive,
		&test.StartTime,
		&test.EndTime,
		&test.CreatedAt,
		&test.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return test, nil
}
