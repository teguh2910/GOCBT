package services

import (
	"database/sql"
	"gocbt/internal/auth"
	"gocbt/internal/models"
	"strings"
)

// QuestionService implements the models.QuestionService interface
type QuestionService struct {
	questionRepo models.QuestionRepository
}

// NewQuestionService creates a new question service
func NewQuestionService(questionRepo models.QuestionRepository) models.QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
	}
}

// CreateQuestion creates a new question
func (s *QuestionService) CreateQuestion(testID int, questionText string, questionType models.QuestionType, marks, orderIndex int) (*models.Question, error) {
	// Validate input
	if strings.TrimSpace(questionText) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if !questionType.IsValid() {
		return nil, auth.ErrInvalidCredentials
	}
	if marks <= 0 {
		return nil, auth.ErrInvalidCredentials
	}

	question := &models.Question{
		TestID:       testID,
		QuestionText: strings.TrimSpace(questionText),
		QuestionType: questionType,
		Marks:        marks,
		OrderIndex:   orderIndex,
	}

	if err := s.questionRepo.Create(question); err != nil {
		return nil, err
	}

	return question, nil
}

// GetQuestion retrieves a question by ID with its options and correct answers
func (s *QuestionService) GetQuestion(questionID int) (*models.Question, error) {
	question, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	if question == nil {
		return nil, auth.ErrUserNotFound
	}

	// Load options for multiple choice questions
	if question.QuestionType == models.QuestionTypeMultipleChoice || question.QuestionType == models.QuestionTypeTrueFalse {
		options, err := s.questionRepo.GetOptionsByQuestionID(questionID)
		if err != nil {
			return nil, err
		}
		question.Options = options
	}

	// Load correct answers for non-multiple choice questions
	if question.QuestionType == models.QuestionTypeShortAnswer {
		answers, err := s.questionRepo.GetCorrectAnswersByQuestionID(questionID)
		if err != nil {
			return nil, err
		}
		question.CorrectAnswers = answers
	}

	return question, nil
}

// GetTestQuestions retrieves all questions for a test with their options and correct answers
func (s *QuestionService) GetTestQuestions(testID int) ([]*models.Question, error) {
	questions, err := s.questionRepo.GetByTestID(testID)
	if err != nil {
		return nil, err
	}

	// Load options and correct answers for each question
	for _, question := range questions {
		if question.QuestionType == models.QuestionTypeMultipleChoice || question.QuestionType == models.QuestionTypeTrueFalse {
			options, err := s.questionRepo.GetOptionsByQuestionID(question.ID)
			if err != nil {
				return nil, err
			}
			question.Options = options
		}

		if question.QuestionType == models.QuestionTypeShortAnswer {
			answers, err := s.questionRepo.GetCorrectAnswersByQuestionID(question.ID)
			if err != nil {
				return nil, err
			}
			question.CorrectAnswers = answers
		}
	}

	return questions, nil
}

// UpdateQuestion updates a question
func (s *QuestionService) UpdateQuestion(questionID int, questionText string, marks, orderIndex int) (*models.Question, error) {
	// Get existing question
	question, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	if question == nil {
		return nil, auth.ErrUserNotFound
	}

	// Validate input
	if strings.TrimSpace(questionText) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if marks <= 0 {
		return nil, auth.ErrInvalidCredentials
	}

	// Update question fields
	question.QuestionText = strings.TrimSpace(questionText)
	question.Marks = marks
	question.OrderIndex = orderIndex

	if err := s.questionRepo.Update(question); err != nil {
		return nil, err
	}

	return question, nil
}

// DeleteQuestion deletes a question
func (s *QuestionService) DeleteQuestion(questionID int) error {
	// Check if question exists
	_, err := s.questionRepo.GetByID(questionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return auth.ErrUserNotFound
		}
		return err
	}

	return s.questionRepo.Delete(questionID)
}

// AddOption adds an option to a question
func (s *QuestionService) AddOption(questionID int, optionText string, isCorrect bool, orderIndex int) (*models.QuestionOption, error) {
	// Validate input
	if strings.TrimSpace(optionText) == "" {
		return nil, auth.ErrInvalidCredentials
	}

	option := &models.QuestionOption{
		QuestionID: questionID,
		OptionText: strings.TrimSpace(optionText),
		IsCorrect:  isCorrect,
		OrderIndex: orderIndex,
	}

	if err := s.questionRepo.CreateOption(option); err != nil {
		return nil, err
	}

	return option, nil
}

// UpdateOption updates a question option
func (s *QuestionService) UpdateOption(optionID int, optionText string, isCorrect bool, orderIndex int) (*models.QuestionOption, error) {
	// Validate input
	if strings.TrimSpace(optionText) == "" {
		return nil, auth.ErrInvalidCredentials
	}

	option := &models.QuestionOption{
		ID:         optionID,
		OptionText: strings.TrimSpace(optionText),
		IsCorrect:  isCorrect,
		OrderIndex: orderIndex,
	}

	if err := s.questionRepo.UpdateOption(option); err != nil {
		return nil, err
	}

	return option, nil
}

// DeleteOption deletes a question option
func (s *QuestionService) DeleteOption(optionID int) error {
	return s.questionRepo.DeleteOption(optionID)
}

// AddCorrectAnswer adds a correct answer to a question
func (s *QuestionService) AddCorrectAnswer(questionID int, answerText string, isCaseSensitive bool) (*models.CorrectAnswer, error) {
	// Validate input
	if strings.TrimSpace(answerText) == "" {
		return nil, auth.ErrInvalidCredentials
	}

	answer := &models.CorrectAnswer{
		QuestionID:      questionID,
		AnswerText:      strings.TrimSpace(answerText),
		IsCaseSensitive: isCaseSensitive,
	}

	if err := s.questionRepo.CreateCorrectAnswer(answer); err != nil {
		return nil, err
	}

	return answer, nil
}

// UpdateCorrectAnswer updates a correct answer
func (s *QuestionService) UpdateCorrectAnswer(answerID int, answerText string, isCaseSensitive bool) (*models.CorrectAnswer, error) {
	// Validate input
	if strings.TrimSpace(answerText) == "" {
		return nil, auth.ErrInvalidCredentials
	}

	answer := &models.CorrectAnswer{
		ID:              answerID,
		AnswerText:      strings.TrimSpace(answerText),
		IsCaseSensitive: isCaseSensitive,
	}

	if err := s.questionRepo.UpdateCorrectAnswer(answer); err != nil {
		return nil, err
	}

	return answer, nil
}

// DeleteCorrectAnswer deletes a correct answer
func (s *QuestionService) DeleteCorrectAnswer(answerID int) error {
	return s.questionRepo.DeleteCorrectAnswer(answerID)
}
