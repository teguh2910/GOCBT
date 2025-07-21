package api

import (
	"gocbt/internal/auth"
	"gocbt/internal/models"
	"gocbt/internal/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ResultHandler handles test result-related requests
type ResultHandler struct {
	resultService models.TestResultService
}

// NewResultHandler creates a new result handler
func NewResultHandler(resultService models.TestResultService) *ResultHandler {
	return &ResultHandler{
		resultService: resultService,
	}
}

// GetResult handles getting a test result by ID
func (h *ResultHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	resultID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid result ID", http.StatusBadRequest)
		return
	}

	result, err := h.resultService.GetResult(resultID)
	if err != nil {
		utils.WriteErrorResponse(w, "Result not found", http.StatusNotFound)
		return
	}

	// Check if user can access this result
	userID, ok := auth.GetUserIDFromContext(r)
	if !ok {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if result.UserID != userID {
		// Allow teachers/admins to view any result
		userRole, ok := auth.GetUserRoleFromContext(r)
		if !ok || (userRole != models.RoleTeacher && userRole != models.RoleAdmin) {
			utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	utils.WriteSuccessResponse(w, result)
}

// GetUserResults handles getting results for a user
func (h *ResultHandler) GetUserResults(w http.ResponseWriter, r *http.Request) {
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

	results, err := h.resultService.GetUserResults(userID, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to get user results", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, results)
}

// GetTestResults handles getting results for a test
func (h *ResultHandler) GetTestResults(w http.ResponseWriter, r *http.Request) {
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

	// Only teachers and admins can view test results
	userRole, ok := auth.GetUserRoleFromContext(r)
	if !ok || (userRole != models.RoleTeacher && userRole != models.RoleAdmin) {
		utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
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

	results, err := h.resultService.GetTestResults(testID, limit, offset)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to get test results", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, results)
}

// GetTestStatistics handles getting statistics for a test
func (h *ResultHandler) GetTestStatistics(w http.ResponseWriter, r *http.Request) {
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

	// Only teachers and admins can view test statistics
	userRole, ok := auth.GetUserRoleFromContext(r)
	if !ok || (userRole != models.RoleTeacher && userRole != models.RoleAdmin) {
		utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
		return
	}

	stats, err := h.resultService.GetTestStatistics(testID)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to get test statistics", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, stats)
}

// CalculateResult handles calculating result for a session
func (h *ResultHandler) CalculateResult(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionID, err := strconv.Atoi(vars["sessionId"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	// Only teachers and admins can manually calculate results
	userRole, ok := auth.GetUserRoleFromContext(r)
	if !ok || (userRole != models.RoleTeacher && userRole != models.RoleAdmin) {
		utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
		return
	}

	result, err := h.resultService.CalculateResult(sessionID)
	if err != nil {
		utils.WriteErrorResponse(w, "Failed to calculate result", http.StatusInternalServerError)
		return
	}

	utils.WriteSuccessResponse(w, result)
}

// GetResultBySession handles getting result by session ID
func (h *ResultHandler) GetResultBySession(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteErrorResponse(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	sessionID, err := strconv.Atoi(vars["sessionId"])
	if err != nil {
		utils.WriteErrorResponse(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	result, err := h.resultService.GetResultBySession(sessionID)
	if err != nil {
		utils.WriteErrorResponse(w, "Result not found", http.StatusNotFound)
		return
	}

	// Check if user can access this result
	userID, ok := auth.GetUserIDFromContext(r)
	if !ok {
		utils.WriteErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if result.UserID != userID {
		// Allow teachers/admins to view any result
		userRole, ok := auth.GetUserRoleFromContext(r)
		if !ok || (userRole != models.RoleTeacher && userRole != models.RoleAdmin) {
			utils.WriteErrorResponse(w, "Forbidden", http.StatusForbidden)
			return
		}
	}

	utils.WriteSuccessResponse(w, result)
}
