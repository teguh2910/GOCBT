package api

import (
	"encoding/json"
	"gocbt/internal/auth"
	"gocbt/internal/models"
	"gocbt/internal/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// TestHandler handles test-related requests
type TestHandler struct {
	testService     models.TestService
	questionService models.QuestionService
}

// NewTestHandler creates a new test handler
func NewTestHandler(testService models.TestService, questionService models.QuestionService) *TestHandler {
	return &TestHandler{
		testService:     testService,
		questionService: questionService,
	}
}

// CreateTestRequest represents a test creation request
type CreateTestRequest struct {
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	Instructions    string     `json:"instructions"`
	DurationMinutes int        `json:"duration_minutes"`
	TotalMarks      int        `json:"total_marks"`
	PassingMarks    int        `json:"passing_marks"`
	StartTime       *time.Time `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
}

// CreateTest handles test creation
func (h *TestHandler) CreateTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r)
	if !ok {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate and sanitize input
	req.Title = utils.SanitizeHTML(utils.SanitizeString(req.Title))
	req.Description = utils.SanitizeHTML(utils.SanitizeString(req.Description))
	req.Instructions = utils.SanitizeHTML(utils.SanitizeString(req.Instructions))

	// Validate required fields
	if utils.IsEmpty(req.Title) {
		utils.WriteErrorResponse(w, "Test title is required", http.StatusBadRequest)
		return
	}

	// Validate text lengths
	if !utils.ValidateTextLength(req.Title, 1, 200) {
		utils.WriteErrorResponse(w, "Test title must be 1-200 characters", http.StatusBadRequest)
		return
	}

	if !utils.ValidateTextLength(req.Description, 0, 1000) {
		utils.WriteErrorResponse(w, "Test description must be 0-1000 characters", http.StatusBadRequest)
		return
	}

	if !utils.ValidateTextLength(req.Instructions, 0, 2000) {
		utils.WriteErrorResponse(w, "Test instructions must be 0-2000 characters", http.StatusBadRequest)
		return
	}

	// Validate duration
	if req.DurationMinutes <= 0 || req.DurationMinutes > 480 { // Max 8 hours
		utils.WriteErrorResponse(w, "Test duration must be between 1 and 480 minutes", http.StatusBadRequest)
		return
	}

	// Validate marks
	if req.TotalMarks <= 0 || req.TotalMarks > 1000 {
		utils.WriteErrorResponse(w, "Total marks must be between 1 and 1000", http.StatusBadRequest)
		return
	}

	if req.PassingMarks < 0 || req.PassingMarks > req.TotalMarks {
		utils.WriteErrorResponse(w, "Passing marks must be between 0 and total marks", http.StatusBadRequest)
		return
	}

	// Check for SQL injection patterns
	if !utils.ValidateNoSQLInjection(req.Title) || !utils.ValidateNoSQLInjection(req.Description) ||
		!utils.ValidateNoSQLInjection(req.Instructions) {
		utils.WriteErrorResponse(w, "Invalid characters detected", http.StatusBadRequest)
		return
	}

	test, err := h.testService.CreateTest(userID, req.Title, req.Description,
		req.Instructions, req.DurationMinutes, req.TotalMarks, req.PassingMarks,
		req.StartTime, req.EndTime)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to create test", http.StatusInternalServerError)
		return
	}

	utils.WriteCreatedResponse(w, test)
}

// GetTest handles getting a test by ID
func (h *TestHandler) GetTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	testID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	test, err := h.testService.GetTest(testID)
	if err != nil {
		utils.WriteErrorResponse(w, "Test not found", http.StatusNotFound)
		return
	}

	utils.WriteSuccessResponse(w, test)
}

// UpdateTest handles test updates
func (h *TestHandler) UpdateTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	testID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	var req CreateTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	test, err := h.testService.UpdateTest(testID, req.Title, req.Description,
		req.Instructions, req.DurationMinutes, req.TotalMarks, req.PassingMarks,
		req.StartTime, req.EndTime)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to update test", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, test)
}

// DeleteTest handles test deletion
func (h *TestHandler) DeleteTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	testID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	if err := h.testService.DeleteTest(testID); err != nil {
		utils.WriteErrorResponse(w, "Failed to delete test", http.StatusInternalServerError)
		return
	}

	utils.WriteNoContentResponse(w)
}

// ListTests handles listing tests
func (h *TestHandler) ListTests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get pagination parameters
	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	// Get creator filter
	creatorID := 0
	if c := r.URL.Query().Get("creator"); c != "" {
		if parsed, err := strconv.Atoi(c); err == nil {
			creatorID = parsed
		}
	}

	tests, err := h.testService.ListTests(creatorID, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to list tests", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, tests)
}

// GetAvailableTests handles getting available tests for students
func (h *TestHandler) GetAvailableTests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r)
	if !ok {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get pagination parameters
	limit := 20
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	tests, err := h.testService.GetAvailableTests(userID, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to get available tests", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, tests)
}

// GetTestQuestions handles getting questions for a test
func (h *TestHandler) GetTestQuestions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	testID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid test ID", http.StatusBadRequest)
		return
	}

	questions, err := h.questionService.GetTestQuestions(testID)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to get test questions", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, questions)
}
