package services

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"gocbt/internal/auth"
	"gocbt/internal/models"
	"strings"
	"time"
)

// TestSessionService implements the models.TestSessionService interface
type TestSessionService struct {
	sessionRepo   models.TestSessionRepository
	answerRepo    models.UserAnswerRepository
	testRepo      models.TestRepository
	questionRepo  models.QuestionRepository
	resultService models.TestResultService
}

// NewTestSessionService creates a new test session service
func NewTestSessionService(sessionRepo models.TestSessionRepository, answerRepo models.UserAnswerRepository, testRepo models.TestRepository, questionRepo models.QuestionRepository, resultService models.TestResultService) models.TestSessionService {
	return &TestSessionService{
		sessionRepo:   sessionRepo,
		answerRepo:    answerRepo,
		testRepo:      testRepo,
		questionRepo:  questionRepo,
		resultService: resultService,
	}
}

// StartSession starts a new test session for a user
func (s *TestSessionService) StartSession(userID, testID int) (*models.TestSession, error) {
	// Check if test exists and is available
	test, err := s.testRepo.GetByID(testID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	if test == nil || !test.IsAvailable() {
		return nil, fmt.Errorf("test is not available")
	}

	// Check if user already has a session for this test
	existingSession, err := s.sessionRepo.GetByUserAndTest(userID, testID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingSession != nil {
		// If session exists and is not expired, return it
		if !existingSession.IsExpired() && existingSession.Status != models.SessionStatusSubmitted {
			return existingSession, nil
		}
	}

	// Generate session token
	token, err := s.generateSessionToken()
	if err != nil {
		return nil, err
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(test.GetDuration())

	// Create new session
	session := &models.TestSession{
		TestID:               testID,
		UserID:               userID,
		SessionToken:         token,
		Status:               models.SessionStatusNotStarted,
		ExpiresAt:            expiresAt,
		TimeRemaining:        &[]int{int(test.GetDuration().Seconds())}[0],
		CurrentQuestionIndex: 0,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession retrieves a session by token
func (s *TestSessionService) GetSession(sessionToken string) (*models.TestSession, error) {
	session, err := s.sessionRepo.GetByToken(sessionToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	if session == nil {
		return nil, auth.ErrUserNotFound
	}

	// Check if session has expired
	if session.IsExpired() && session.Status == models.SessionStatusInProgress {
		session.Status = models.SessionStatusExpired
		s.sessionRepo.Update(session)
	}

	return session, nil
}

// SubmitAnswer submits an answer for a question in a session
func (s *TestSessionService) SubmitAnswer(sessionToken string, questionID int, answerText *string, selectedOptionID *int) (*models.UserAnswer, error) {
	// Get session
	session, err := s.GetSession(sessionToken)
	if err != nil {
		return nil, err
	}

	// Check if session is available for answers (not_started or in_progress and not expired)
	if session.IsExpired() || (session.Status != models.SessionStatusNotStarted && session.Status != models.SessionStatusInProgress) {
		return nil, fmt.Errorf("session is not available for answers")
	}

	// Start session if not started
	if session.Status == models.SessionStatusNotStarted {
		now := time.Now()
		session.Status = models.SessionStatusInProgress
		session.StartedAt = &now
		if err := s.sessionRepo.Update(session); err != nil {
			return nil, err
		}
	}

	// Get question
	question, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		return nil, err
	}

	if question == nil || question.TestID != session.TestID {
		return nil, fmt.Errorf("invalid question for this test")
	}

	// Check if answer already exists
	existingAnswer, err := s.answerRepo.GetBySessionAndQuestion(session.ID, questionID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Validate and score the answer
	isCorrect, marksAwarded := s.scoreAnswer(question, answerText, selectedOptionID)

	if existingAnswer != nil {
		// Update existing answer
		existingAnswer.AnswerText = answerText
		existingAnswer.SelectedOptionID = selectedOptionID
		existingAnswer.IsCorrect = &isCorrect
		existingAnswer.MarksAwarded = marksAwarded
		if err := s.answerRepo.Update(existingAnswer); err != nil {
			return nil, err
		}
		return existingAnswer, nil
	}

	// Create new answer
	answer := &models.UserAnswer{
		SessionID:        session.ID,
		QuestionID:       questionID,
		AnswerText:       answerText,
		SelectedOptionID: selectedOptionID,
		IsCorrect:        &isCorrect,
		MarksAwarded:     marksAwarded,
	}

	if err := s.answerRepo.Create(answer); err != nil {
		return nil, err
	}

	return answer, nil
}

// GetSessionAnswers retrieves all answers for a session
func (s *TestSessionService) GetSessionAnswers(sessionToken string) ([]*models.UserAnswer, error) {
	session, err := s.GetSession(sessionToken)
	if err != nil {
		return nil, err
	}

	return s.answerRepo.GetBySession(session.ID)
}

// SubmitSession submits a test session and calculates results
func (s *TestSessionService) SubmitSession(sessionToken string) (*models.TestSession, error) {
	session, err := s.GetSession(sessionToken)
	if err != nil {
		return nil, err
	}

	if session.Status == models.SessionStatusSubmitted {
		return session, nil
	}

	// Mark session as submitted
	now := time.Now()
	session.Status = models.SessionStatusSubmitted
	session.SubmittedAt = &now

	if err := s.sessionRepo.Update(session); err != nil {
		return nil, err
	}

	// Automatically calculate results if result service is available
	if s.resultService != nil {
		_, err := s.resultService.CalculateResult(session.ID)
		if err != nil {
			// Log error but don't fail the submission
			// In a production system, you might want to queue this for retry
			fmt.Printf("Warning: Failed to calculate result for session %d: %v\n", session.ID, err)
		}
	}

	return session, nil
}

// GetUserSessions retrieves sessions for a user
func (s *TestSessionService) GetUserSessions(userID int, limit, offset int) ([]*models.TestSession, error) {
	return s.sessionRepo.GetUserSessions(userID, limit, offset)
}

// UpdateSessionProgress updates the current question index for a session
func (s *TestSessionService) UpdateSessionProgress(sessionToken string, currentQuestionIndex int) error {
	session, err := s.GetSession(sessionToken)
	if err != nil {
		return err
	}

	if !session.IsActive() {
		return fmt.Errorf("session is not active")
	}

	session.CurrentQuestionIndex = currentQuestionIndex
	return s.sessionRepo.Update(session)
}

// generateSessionToken generates a random session token
func (s *TestSessionService) generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// scoreAnswer scores an answer based on the question type and correct answers
func (s *TestSessionService) scoreAnswer(question *models.Question, answerText *string, selectedOptionID *int) (bool, int) {
	switch question.QuestionType {
	case models.QuestionTypeMultipleChoice, models.QuestionTypeTrueFalse:
		return s.scoreMultipleChoiceAnswer(question, selectedOptionID)
	case models.QuestionTypeShortAnswer:
		return s.scoreShortAnswer(question, answerText)
	default:
		return false, 0
	}
}

// scoreMultipleChoiceAnswer scores a multiple choice answer
func (s *TestSessionService) scoreMultipleChoiceAnswer(question *models.Question, selectedOptionID *int) (bool, int) {
	if selectedOptionID == nil {
		return false, 0
	}

	// Get options for the question
	options, err := s.questionRepo.GetOptionsByQuestionID(question.ID)
	if err != nil {
		return false, 0
	}

	for _, option := range options {
		if option.ID == *selectedOptionID && option.IsCorrect {
			return true, question.Marks
		}
	}

	return false, 0
}

// scoreShortAnswer scores a short answer
func (s *TestSessionService) scoreShortAnswer(question *models.Question, answerText *string) (bool, int) {
	if answerText == nil || strings.TrimSpace(*answerText) == "" {
		return false, 0
	}

	// Get correct answers for the question
	correctAnswers, err := s.questionRepo.GetCorrectAnswersByQuestionID(question.ID)
	if err != nil {
		return false, 0
	}

	userAnswer := strings.TrimSpace(*answerText)

	for _, correctAnswer := range correctAnswers {
		expectedAnswer := correctAnswer.AnswerText
		if !correctAnswer.IsCaseSensitive {
			userAnswer = strings.ToLower(userAnswer)
			expectedAnswer = strings.ToLower(expectedAnswer)
		}

		if userAnswer == expectedAnswer {
			return true, question.Marks
		}
	}

	return false, 0
}
