package api

import (
	"encoding/json"
	"gocbt/internal/models"
	"gocbt/internal/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// QuestionHandler handles question-related requests
type QuestionHandler struct {
	questionService models.QuestionService
}

// NewQuestionHandler creates a new question handler
func NewQuestionHandler(questionService models.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

// CreateQuestionRequest represents a question creation request
type CreateQuestionRequest struct {
	TestID       int                    `json:"test_id"`
	QuestionText string                 `json:"question_text"`
	QuestionType models.QuestionType    `json:"question_type"`
	Marks        int                    `json:"marks"`
	OrderIndex   int                    `json:"order_index"`
	Options      []CreateOptionRequest  `json:"options,omitempty"`
	Answers      []CreateAnswerRequest  `json:"answers,omitempty"`
}

// CreateOptionRequest represents an option creation request
type CreateOptionRequest struct {
	OptionText string `json:"option_text"`
	IsCorrect  bool   `json:"is_correct"`
	OrderIndex int    `json:"order_index"`
}

// CreateAnswerRequest represents a correct answer creation request
type CreateAnswerRequest struct {
	AnswerText      string `json:"answer_text"`
	IsCaseSensitive bool   `json:"is_case_sensitive"`
}

// CreateQuestion handles question creation
func (h *QuestionHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateQuestionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.CreateQuestion(req.TestID, req.QuestionText,
		req.QuestionType, req.Marks, req.OrderIndex)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to create question", http.StatusInternalServerError)
		return
	}

	// Add options for multiple choice questions
	if req.QuestionType == models.QuestionTypeMultipleChoice || req.QuestionType == models.QuestionTypeTrueFalse {
		for _, optionReq := range req.Options {
			option, err := h.questionService.AddOption(question.ID, optionReq.OptionText,
				optionReq.IsCorrect, optionReq.OrderIndex)
			if err != nil {
				utils.WriteErrorResponse(w, "Failed to create question option", http.StatusInternalServerError)
				return
			}
			question.Options = append(question.Options, option)
		}
	}

	// Add correct answers for short answer questions
	if req.QuestionType == models.QuestionTypeShortAnswer {
		for _, answerReq := range req.Answers {
			answer, err := h.questionService.AddCorrectAnswer(question.ID, answerReq.AnswerText,
				answerReq.IsCaseSensitive)
			if err != nil {
				utils.WriteErrorResponse(w, "Failed to create correct answer", http.StatusInternalServerError)
				return
			}
			question.CorrectAnswers = append(question.CorrectAnswers, answer)
		}
	}

	utils.WriteCreatedResponse(w, question)
}

// GetQuestion handles getting a question by ID
func (h *QuestionHandler) GetQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	questionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.GetQuestion(questionID)
	if err != nil {
		utils.WriteErrorResponse(w, "Question not found", http.StatusNotFound)
		return
	}

	utils.WriteSuccessResponse(w, question)
}

// UpdateQuestion handles question updates
func (h *QuestionHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	questionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req struct {
		QuestionText string `json:"question_text"`
		Marks        int    `json:"marks"`
		OrderIndex   int    `json:"order_index"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	question, err := h.questionService.UpdateQuestion(questionID, req.QuestionText,
		req.Marks, req.OrderIndex)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to update question", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, question)
}

// DeleteQuestion handles question deletion
func (h *QuestionHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	questionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	if err := h.questionService.DeleteQuestion(questionID); err != nil {
		utils.WriteErrorResponse(w, "Failed to delete question", http.StatusInternalServerError)
		return
	}

	utils.WriteNoContentResponse(w)
}

// AddOption handles adding an option to a question
func (h *QuestionHandler) AddOption(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	questionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req CreateOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	option, err := h.questionService.AddOption(questionID, req.OptionText,
		req.IsCorrect, req.OrderIndex)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to add option", http.StatusInternalServerError)
		return
	}

	utils.WriteCreatedResponse(w, option)
}

// UpdateOption handles updating a question option
func (h *QuestionHandler) UpdateOption(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	optionID, err := strconv.Atoi(vars["optionId"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid option ID", http.StatusBadRequest)
		return
	}

	var req CreateOptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	option, err := h.questionService.UpdateOption(optionID, req.OptionText,
		req.IsCorrect, req.OrderIndex)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to update option", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, option)
}

// DeleteOption handles deleting a question option
func (h *QuestionHandler) DeleteOption(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	optionID, err := strconv.Atoi(vars["optionId"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid option ID", http.StatusBadRequest)
		return
	}

	if err := h.questionService.DeleteOption(optionID); err != nil {
		utils.WriteErrorResponse(w, "Failed to delete option", http.StatusInternalServerError)
		return
	}

	utils.WriteNoContentResponse(w)
}

// AddCorrectAnswer handles adding a correct answer to a question
func (h *QuestionHandler) AddCorrectAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	questionID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	var req CreateAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	answer, err := h.questionService.AddCorrectAnswer(questionID, req.AnswerText,
		req.IsCaseSensitive)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to add correct answer", http.StatusInternalServerError)
		return
	}

	utils.WriteCreatedResponse(w, answer)
}
