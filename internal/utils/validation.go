package utils

import (
	"errors"
	"html"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername validates username format
func ValidateUsername(username string) bool {
	// Username should be 3-50 characters, alphanumeric and underscores only
	if len(username) < 3 || len(username) > 50 {
		return false
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	return usernameRegex.MatchString(username)
}

// SanitizeString trims whitespace and removes extra spaces
func SanitizeString(s string) string {
	// Trim leading/trailing whitespace
	s = strings.TrimSpace(s)
	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	return spaceRegex.ReplaceAllString(s, " ")
}

// IsEmpty checks if a string is empty after trimming
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// ValidateRequired validates that required fields are not empty
func ValidateRequired(fields map[string]string) []string {
	var errors []string
	for field, value := range fields {
		if IsEmpty(value) {
			errors = append(errors, field+" is required")
		}
	}
	return errors
}

// SanitizeHTML removes HTML tags and escapes HTML entities
func SanitizeHTML(input string) string {
	// Remove HTML tags
	htmlTagRegex := regexp.MustCompile(`<[^>]*>`)
	cleaned := htmlTagRegex.ReplaceAllString(input, "")
	// Escape HTML entities
	return html.EscapeString(cleaned)
}

// ValidateInteger validates and converts string to integer with bounds
func ValidateInteger(value string, min, max int) (int, error) {
	if value == "" {
		return 0, errors.New("value is required")
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	if num < min || num > max {
		return 0, errors.New("value out of range")
	}

	return num, nil
}

// ValidateTextLength validates text length and content
func ValidateTextLength(text string, minLen, maxLen int) bool {
	if len(text) < minLen || len(text) > maxLen {
		return false
	}

	// Check for valid Unicode characters
	for _, r := range text {
		if !unicode.IsPrint(r) && !unicode.IsSpace(r) {
			return false
		}
	}

	return true
}

// ValidateNoSQLInjection checks for common SQL injection patterns
func ValidateNoSQLInjection(input string) bool {
	// Convert to lowercase for pattern matching
	lower := strings.ToLower(input)

	// Common SQL injection patterns
	sqlPatterns := []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
		"union", "select", "insert", "update", "delete",
		"drop", "create", "alter", "exec", "execute",
		"script", "javascript", "vbscript", "onload",
		"onerror", "onclick", "<script", "</script>",
	}

	for _, pattern := range sqlPatterns {
		if strings.Contains(lower, pattern) {
			return false
		}
	}

	return true
}

// ValidateRole validates user role
func ValidateRole(role string) bool {
	validRoles := []string{"admin", "teacher", "student"}
	for _, validRole := range validRoles {
		if role == validRole {
			return true
		}
	}
	return false
}
