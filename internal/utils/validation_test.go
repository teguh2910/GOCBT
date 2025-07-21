package utils

import (
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email    string
		expected bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"invalid-email", false},
		{"@domain.com", false},
		{"user@", false},
		{"", false},
		{"user@domain", false},
	}

	for _, test := range tests {
		result := ValidateEmail(test.email)
		if result != test.expected {
			t.Errorf("ValidateEmail(%s) = %v, expected %v", test.email, result, test.expected)
		}
	}
}

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		username string
		expected bool
	}{
		{"validuser", true},
		{"user123", true},
		{"user_name", true},
		{"ab", false},        // too short
		{"", false},          // empty
		{"user@name", false}, // invalid character
		{"user name", false}, // space
		{"user-name", false}, // hyphen
	}

	for _, test := range tests {
		result := ValidateUsername(test.username)
		if result != test.expected {
			t.Errorf("ValidateUsername(%s) = %v, expected %v", test.username, result, test.expected)
		}
	}
}

func TestValidateNoSQLInjection(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"normal text", true},
		{"user123", true},
		{"SELECT * FROM users", false},
		{"'; DROP TABLE users; --", false},
		{"<script>alert('xss')</script>", false},
		{"javascript:alert(1)", false},
		{"onload=alert(1)", false},
		{"UNION SELECT", false},
		{"/*comment*/", false},
	}

	for _, test := range tests {
		result := ValidateNoSQLInjection(test.input)
		if result != test.expected {
			t.Errorf("ValidateNoSQLInjection(%s) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestSanitizeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<script>alert('xss')</script>", "alert(&#39;xss&#39;)"},
		{"<b>bold</b> text", "bold text"},
		{"normal text", "normal text"},
		{"<img src=x onerror=alert(1)>", ""},
		{"<p>paragraph</p>", "paragraph"},
	}

	for _, test := range tests {
		result := SanitizeHTML(test.input)
		if result != test.expected {
			t.Errorf("SanitizeHTML(%s) = %s, expected %s", test.input, result, test.expected)
		}
	}
}

func TestValidateTextLength(t *testing.T) {
	tests := []struct {
		text     string
		minLen   int
		maxLen   int
		expected bool
	}{
		{"hello", 1, 10, true},
		{"", 1, 10, false},            // too short
		{"verylongtext", 1, 5, false}, // too long
		{"ok", 2, 5, true},
		{"  spaced  ", 5, 15, true}, // should trim
	}

	for _, test := range tests {
		result := ValidateTextLength(test.text, test.minLen, test.maxLen)
		if result != test.expected {
			t.Errorf("ValidateTextLength(%s, %d, %d) = %v, expected %v",
				test.text, test.minLen, test.maxLen, result, test.expected)
		}
	}
}

func TestValidateRole(t *testing.T) {
	tests := []struct {
		role     string
		expected bool
	}{
		{"admin", true},
		{"teacher", true},
		{"student", true},
		{"invalid", false},
		{"", false},
		{"ADMIN", false}, // case sensitive
	}

	for _, test := range tests {
		result := ValidateRole(test.role)
		if result != test.expected {
			t.Errorf("ValidateRole(%s) = %v, expected %v", test.role, result, test.expected)
		}
	}
}

func TestValidateInteger(t *testing.T) {
	tests := []struct {
		value    string
		min      int
		max      int
		expected bool
	}{
		{"5", 1, 10, true},
		{"0", 1, 10, false},   // below min
		{"15", 1, 10, false},  // above max
		{"abc", 1, 10, false}, // not a number
		{"", 1, 10, false},    // empty
	}

	for _, test := range tests {
		_, err := ValidateInteger(test.value, test.min, test.max)
		result := err == nil
		if result != test.expected {
			t.Errorf("ValidateInteger(%s, %d, %d) error = %v, expected success = %v",
				test.value, test.min, test.max, err, test.expected)
		}
	}
}

func BenchmarkValidateEmail(b *testing.B) {
	email := "test@example.com"
	for i := 0; i < b.N; i++ {
		ValidateEmail(email)
	}
}

func BenchmarkValidateNoSQLInjection(b *testing.B) {
	input := "normal user input text"
	for i := 0; i < b.N; i++ {
		ValidateNoSQLInjection(input)
	}
}

func BenchmarkSanitizeHTML(b *testing.B) {
	input := "<p>Some HTML content with <b>bold</b> text</p>"
	for i := 0; i < b.N; i++ {
		SanitizeHTML(input)
	}
}
