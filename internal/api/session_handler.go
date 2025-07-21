package api

import (
	"encoding/json"
	"fmt"
	"gocbt/internal/auth"
	"gocbt/internal/models"
	"gocbt/internal/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// SessionHandler handles test session-related requests
type SessionHandler struct {
	sessionService models.TestSessionService
}

// NewSessionHandler creates a new session handler
func NewSessionHandler(sessionService models.TestSessionService) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
	}
}

// StartSessionRequest represents a session start request
type StartSessionRequest struct {
	TestID int `json:"test_id"`
}

// SubmitAnswerRequest represents an answer submission request
type SubmitAnswerRequest struct {
	QuestionID       int     `json:"question_id"`
	AnswerText       *string `json:"answer_text,omitempty"`
	SelectedOptionID *int    `json:"selected_option_id,omitempty"`
}

// UpdateProgressRequest represents a progress update request
type UpdateProgressRequest struct {
	CurrentQuestionIndex int `json:"current_question_index"`
}

// SessionResponse represents a session response with additional info
type SessionResponse struct {
	*models.TestSession
	RemainingTime int `json:"remaining_time_seconds"`
}

// StartSession handles starting a new test session
func (h *SessionHandler) StartSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r)
	if !ok {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req StartSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, err := h.sessionService.StartSession(userID, req.TestID)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to start session", http.StatusInternalServerError)
		return
	}

	response := &SessionResponse{
		TestSession:   session,
		RemainingTime: session.GetRemainingTime(),
	}

	utils.WriteCreatedResponse(w, response)
}

// GetSession handles getting session information
func (h *SessionHandler) GetSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionToken := vars["token"]

	session, err := h.sessionService.GetSession(sessionToken)
	if err != nil {
		utils.WriteErrorResponse(w, "Session not found", http.StatusNotFound)
		return
	}

	// Check if user owns this session
	userID, ok := auth.GetUserIDFromContext(r)
	if !ok {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if session.UserID != userID {
		// Allow teachers/admins to view any session
		userRole, ok := auth.GetUserRoleFromContext(r)
		if !ok || (userRole != models.RoleTeacher && userRole != models.RoleAdmin) {
			utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	response := &SessionResponse{
		TestSession:   session,
		RemainingTime: session.GetRemainingTime(),
	}

	utils.WriteSuccessResponse(w, response)
}

// SubmitAnswer handles answer submission
func (h *SessionHandler) SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionToken := vars["token"]

	// Validate session token format
	if utils.IsEmpty(sessionToken) || !utils.ValidateNoSQLInjection(sessionToken) {
		utils.WriteErrorResponse(w, "Invalid session token", http.StatusBadRequest)
		return
	}

	var req SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate question ID
	if req.QuestionID <= 0 {
		utils.WriteErrorResponse(w, "Invalid question ID", http.StatusBadRequest)
		return
	}

	// Sanitize answer text if provided
	if req.AnswerText != nil {
		sanitized := utils.SanitizeHTML(*req.AnswerText)
		req.AnswerText = &sanitized

		// Validate answer text length
		if !utils.ValidateTextLength(*req.AnswerText, 0, 5000) {
			utils.WriteErrorResponse(w, "Answer text too long", http.StatusBadRequest)
			return
		}
	}

	// Validate selected option ID if provided
	if req.SelectedOptionID != nil && *req.SelectedOptionID <= 0 {
		utils.WriteErrorResponse(w, "Invalid option ID", http.StatusBadRequest)
		return
	}

	// Verify user owns this session
	session, err := h.sessionService.GetSession(sessionToken)
	if err != nil {
		utils.WriteErrorResponse(w, "Session not found", http.StatusNotFound)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r)
	if !ok || session.UserID != userID {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	answer, err := h.sessionService.SubmitAnswer(sessionToken, req.QuestionID, req.AnswerText, req.SelectedOptionID)
	if err != nil {
		utils.WriteErrorResponse(w, fmt.Sprintf("Failed to submit answer: %v", err), http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, answer)
}

// GetSessionAnswers handles getting all answers for a session
func (h *SessionHandler) GetSessionAnswers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionToken := vars["token"]

	// Verify user owns this session or is teacher/admin
	session, err := h.sessionService.GetSession(sessionToken)
	if err != nil {
		utils.WriteErrorResponse(w, "Session not found", http.StatusNotFound)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r)
	if !ok {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if session.UserID != userID {
		// Allow teachers/admins to view any session answers
		userRole, ok := auth.GetUserRoleFromContext(r)
		if !ok || (userRole != models.RoleTeacher && userRole != models.RoleAdmin) {
			utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	answers, err := h.sessionService.GetSessionAnswers(sessionToken)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to get session answers", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, answers)
}

// SubmitSession handles session submission
func (h *SessionHandler) SubmitSession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionToken := vars["token"]

	// Verify user owns this session
	session, err := h.sessionService.GetSession(sessionToken)
	if err != nil {
		utils.WriteErrorResponse(w, "Session not found", http.StatusNotFound)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r)
	if !ok || session.UserID != userID {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	submittedSession, err := h.sessionService.SubmitSession(sessionToken)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to submit session", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, submittedSession)
}

// UpdateProgress handles updating session progress
func (h *SessionHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionToken := vars["token"]

	var req UpdateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify user owns this session
	session, err := h.sessionService.GetSession(sessionToken)
	if err != nil {
		utils.WriteErrorResponse(w, "Session not found", http.StatusNotFound)
		return
	}

	userID, ok := auth.GetUserIDFromContext(r)
	if !ok || session.UserID != userID {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.sessionService.UpdateSessionProgress(sessionToken, req.CurrentQuestionIndex); err != nil {
		utils.WriteErrorResponse(w, "Failed to update progress", http.StatusInternalServerError)
		return
	}

	utils.WriteNoContentResponse(w)
}

// GetUserSessions handles getting sessions for a user
func (h *SessionHandler) GetUserSessions(w http.ResponseWriter, r *http.Request) {
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

	sessions, err := h.sessionService.GetUserSessions(userID, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to get user sessions", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, sessions)
}
