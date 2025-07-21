package auth

import (
	"context"
	"gocbt/internal/models"
	"net/http"
	"strings"
)

// ContextKey represents a key for context values
type ContextKey string

const (
	// UserContextKey is the key for user in context
	UserContextKey ContextKey = "user"
	// ClaimsContextKey is the key for JWT claims in context
	ClaimsContextKey ContextKey = "claims"
)

// Middleware provides authentication middleware
type Middleware struct {
	jwtManager *JWTManager
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(jwtManager *JWTManager) *Middleware {
	return &Middleware{
		jwtManager: jwtManager,
	}
}

// Authenticate is a middleware that validates JWT tokens
func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			http.Error(w, "Token required", http.StatusUnauthorized)
			return
		}

		// Validate token
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole is a middleware that requires specific user roles
func (m *Middleware) RequireRole(roles ...models.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get claims from context
			claims, ok := r.Context().Value(ClaimsContextKey).(*Claims)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if claims.Role == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireTeacherOrAdmin is a convenience middleware for teacher/admin access
func (m *Middleware) RequireTeacherOrAdmin(next http.Handler) http.Handler {
	return m.RequireRole(models.RoleTeacher, models.RoleAdmin)(next)
}

// RequireAdmin is a convenience middleware for admin-only access
func (m *Middleware) RequireAdmin(next http.Handler) http.Handler {
	return m.RequireRole(models.RoleAdmin)(next)
}

// OptionalAuth is a middleware that optionally validates JWT tokens
func (m *Middleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString != "" {
				// Validate token if present
				if claims, err := m.jwtManager.ValidateToken(tokenString); err == nil {
					// Add claims to context if valid
					ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)
					r = r.WithContext(ctx)
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts user claims from request context
func GetUserFromContext(r *http.Request) (*Claims, bool) {
	claims, ok := r.Context().Value(ClaimsContextKey).(*Claims)
	return claims, ok
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(r *http.Request) (int, bool) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		return 0, false
	}
	return claims.UserID, true
}

// GetUserRoleFromContext extracts user role from request context
func GetUserRoleFromContext(r *http.Request) (models.UserRole, bool) {
	claims, ok := GetUserFromContext(r)
	if !ok {
		return "", false
	}
	return claims.Role, true
}
