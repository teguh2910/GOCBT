package services

import (
	"database/sql"
	"gocbt/internal/auth"
	"gocbt/internal/models"
	"strings"
	"time"
)

// TestService implements the models.TestService interface
type TestService struct {
	testRepo models.TestRepository
}

// NewTestService creates a new test service
func NewTestService(testRepo models.TestRepository) models.TestService {
	return &TestService{
		testRepo: testRepo,
	}
}

// CreateTest creates a new test
func (s *TestService) CreateTest(creatorID int, title, description, instructions string, durationMinutes, totalMarks, passingMarks int, startTime, endTime *time.Time) (*models.Test, error) {
	// Validate input
	if strings.TrimSpace(title) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if durationMinutes <= 0 {
		return nil, auth.ErrInvalidCredentials
	}
	if totalMarks <= 0 {
		return nil, auth.ErrInvalidCredentials
	}
	if passingMarks < 0 || passingMarks > totalMarks {
		return nil, auth.ErrInvalidCredentials
	}

	// Validate time window
	if startTime != nil && endTime != nil && startTime.After(*endTime) {
		return nil, auth.ErrInvalidCredentials
	}

	test := &models.Test{
		Title:           strings.TrimSpace(title),
		Description:     strings.TrimSpace(description),
		CreatedBy:       creatorID,
		DurationMinutes: durationMinutes,
		TotalMarks:      totalMarks,
		PassingMarks:    passingMarks,
		Instructions:    strings.TrimSpace(instructions),
		IsActive:        true,
		StartTime:       startTime,
		EndTime:         endTime,
	}

	if err := s.testRepo.Create(test); err != nil {
		return nil, err
	}

	return test, nil
}

// GetTest retrieves a test by ID
func (s *TestService) GetTest(testID int) (*models.Test, error) {
	test, err := s.testRepo.GetByID(testID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound // Reusing error for "not found"
		}
		return nil, err
	}

	if test == nil {
		return nil, auth.ErrUserNotFound
	}

	return test, nil
}

// UpdateTest updates a test
func (s *TestService) UpdateTest(testID int, title, description, instructions string, durationMinutes, totalMarks, passingMarks int, startTime, endTime *time.Time) (*models.Test, error) {
	// Get existing test
	test, err := s.GetTest(testID)
	if err != nil {
		return nil, err
	}

	// Validate input
	if strings.TrimSpace(title) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if durationMinutes <= 0 {
		return nil, auth.ErrInvalidCredentials
	}
	if totalMarks <= 0 {
		return nil, auth.ErrInvalidCredentials
	}
	if passingMarks < 0 || passingMarks > totalMarks {
		return nil, auth.ErrInvalidCredentials
	}

	// Validate time window
	if startTime != nil && endTime != nil && startTime.After(*endTime) {
		return nil, auth.ErrInvalidCredentials
	}

	// Update test fields
	test.Title = strings.TrimSpace(title)
	test.Description = strings.TrimSpace(description)
	test.Instructions = strings.TrimSpace(instructions)
	test.DurationMinutes = durationMinutes
	test.TotalMarks = totalMarks
	test.PassingMarks = passingMarks
	test.StartTime = startTime
	test.EndTime = endTime

	if err := s.testRepo.Update(test); err != nil {
		return nil, err
	}

	return test, nil
}

// DeleteTest deletes a test
func (s *TestService) DeleteTest(testID int) error {
	// Check if test exists
	_, err := s.GetTest(testID)
	if err != nil {
		return err
	}

	return s.testRepo.Delete(testID)
}

// ListTests retrieves tests by creator with pagination
func (s *TestService) ListTests(creatorID int, limit, offset int) ([]*models.Test, error) {
	if creatorID == 0 {
		return s.testRepo.List(limit, offset)
	}
	return s.testRepo.GetByCreator(creatorID, limit, offset)
}

// GetAvailableTests retrieves tests available for a user
func (s *TestService) GetAvailableTests(userID int, limit, offset int) ([]*models.Test, error) {
	return s.testRepo.GetAvailableTests(userID, limit, offset)
}

// ActivateTest activates a test
func (s *TestService) ActivateTest(testID int) error {
	test, err := s.GetTest(testID)
	if err != nil {
		return err
	}

	test.IsActive = true
	return s.testRepo.Update(test)
}

// DeactivateTest deactivates a test
func (s *TestService) DeactivateTest(testID int) error {
	test, err := s.GetTest(testID)
	if err != nil {
		return err
	}

	test.IsActive = false
	return s.testRepo.Update(test)
}
