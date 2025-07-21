package services

import (
	"database/sql"
	"gocbt/internal/auth"
	"gocbt/internal/models"
)

// TestResultService implements the models.TestResultService interface
type TestResultService struct {
	resultRepo   models.TestResultRepository
	sessionRepo  models.TestSessionRepository
	answerRepo   models.UserAnswerRepository
	testRepo     models.TestRepository
	questionRepo models.QuestionRepository
}

// NewTestResultService creates a new test result service
func NewTestResultService(resultRepo models.TestResultRepository, sessionRepo models.TestSessionRepository, answerRepo models.UserAnswerRepository, testRepo models.TestRepository, questionRepo models.QuestionRepository) models.TestResultService {
	return &TestResultService{
		resultRepo:   resultRepo,
		sessionRepo:  sessionRepo,
		answerRepo:   answerRepo,
		testRepo:     testRepo,
		questionRepo: questionRepo,
	}
}

// CalculateResult calculates and stores the result for a test session
func (s *TestResultService) CalculateResult(sessionID int) (*models.TestResult, error) {
	// Check if result already exists
	existingResult, err := s.resultRepo.GetBySessionID(sessionID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingResult != nil {
		return existingResult, nil
	}

	// Get session
	session, err := s.sessionRepo.GetByID(sessionID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, auth.ErrUserNotFound
	}

	// Get test
	test, err := s.testRepo.GetByID(session.TestID)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, auth.ErrUserNotFound
	}

	// Get all questions for the test
	questions, err := s.questionRepo.GetByTestID(session.TestID)
	if err != nil {
		return nil, err
	}

	// Get all answers for the session
	answers, err := s.answerRepo.GetBySession(sessionID)
	if err != nil {
		return nil, err
	}

	// Calculate statistics
	totalQuestions := len(questions)
	answeredQuestions := len(answers)
	correctAnswers := 0
	marksObtained := 0

	for _, answer := range answers {
		if answer.IsCorrect != nil && *answer.IsCorrect {
			correctAnswers++
		}
		marksObtained += answer.MarksAwarded
	}

	// Calculate percentage
	var percentage float64
	if test.TotalMarks > 0 {
		percentage = (float64(marksObtained) / float64(test.TotalMarks)) * 100
	}

	// Determine if passed
	isPassed := marksObtained >= test.PassingMarks

	// Calculate time taken
	var timeTaken *int
	if session.StartedAt != nil && session.SubmittedAt != nil {
		duration := int(session.SubmittedAt.Sub(*session.StartedAt).Seconds())
		timeTaken = &duration
	}

	// Create result
	result := &models.TestResult{
		SessionID:         sessionID,
		TestID:            session.TestID,
		UserID:            session.UserID,
		TotalQuestions:    totalQuestions,
		AnsweredQuestions: answeredQuestions,
		CorrectAnswers:    correctAnswers,
		TotalMarks:        test.TotalMarks,
		MarksObtained:     marksObtained,
		Percentage:        percentage,
		IsPassed:          isPassed,
		TimeTaken:         timeTaken,
	}

	// Calculate grade
	grade := result.CalculateGrade()
	result.Grade = &grade

	// Save result
	if err := s.resultRepo.Create(result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetResult retrieves a test result by ID
func (s *TestResultService) GetResult(resultID int) (*models.TestResult, error) {
	result, err := s.resultRepo.GetByID(resultID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	if result == nil {
		return nil, auth.ErrUserNotFound
	}

	return result, nil
}

// GetUserResults retrieves test results for a user with pagination
func (s *TestResultService) GetUserResults(userID int, limit, offset int) ([]*models.TestResult, error) {
	return s.resultRepo.GetByUser(userID, limit, offset)
}

// GetTestResults retrieves test results for a test with pagination
func (s *TestResultService) GetTestResults(testID int, limit, offset int) ([]*models.TestResult, error) {
	return s.resultRepo.GetByTest(testID, limit, offset)
}

// GetTestStatistics retrieves statistics for a test
func (s *TestResultService) GetTestStatistics(testID int) (*models.TestStatistics, error) {
	return s.resultRepo.GetTestStatistics(testID)
}

// GetResultBySession retrieves a result by session ID
func (s *TestResultService) GetResultBySession(sessionID int) (*models.TestResult, error) {
	result, err := s.resultRepo.GetBySessionID(sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	if result == nil {
		return nil, auth.ErrUserNotFound
	}

	return result, nil
}

// GetResultByUserAndTest retrieves the latest result for a user and test
func (s *TestResultService) GetResultByUserAndTest(userID, testID int) (*models.TestResult, error) {
	result, err := s.resultRepo.GetByUserAndTest(userID, testID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	if result == nil {
		return nil, auth.ErrUserNotFound
	}

	return result, nil
}
