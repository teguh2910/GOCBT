package auth

import (
	"errors"
	"fmt"
	"gocbt/internal/config"
	"gocbt/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	UserID   int             `json:"user_id"`
	Username string          `json:"username"`
	Role     models.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey  string
	expiration time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(cfg *config.JWTConfig) *JWTManager {
	return &JWTManager{
		secretKey:  cfg.Secret,
		expiration: cfg.Expiration,
	}
}

// GenerateToken generates a JWT token for a user
func (j *JWTManager) GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gocbt",
			Subject:   fmt.Sprintf("user:%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secretKey))
}

// ValidateToken validates a JWT token and returns the claims
func (j *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	// Basic token format validation
	if tokenString == "" {
		return nil, errors.New("empty token")
	}

	// Check token length to prevent extremely long tokens
	if len(tokenString) > 2048 {
		return nil, errors.New("token too long")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method - only allow HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Ensure it's specifically HS256
		if token.Method.Alg() != "HS256" {
			return nil, fmt.Errorf("unexpected algorithm: %v", token.Method.Alg())
		}

		return []byte(j.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Additional validation
		now := time.Now()

		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(now) {
			return nil, errors.New("token expired")
		}

		// Check if token is used before valid time
		if claims.NotBefore != nil && claims.NotBefore.Time.After(now) {
			return nil, errors.New("token not valid yet")
		}

		// Validate issuer
		if claims.Issuer != "gocbt" {
			return nil, errors.New("invalid issuer")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new token with extended expiration
func (j *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Check if token is close to expiration (within 1 hour)
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", errors.New("token is not close to expiration")
	}

	// Create new claims with extended expiration
	newClaims := &Claims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gocbt",
			Subject:   claims.Subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)
	return token.SignedString([]byte(j.secretKey))
}

// ExtractUserID extracts user ID from token string
func (j *JWTManager) ExtractUserID(tokenString string) (int, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ExtractUserRole extracts user role from token string
func (j *JWTManager) ExtractUserRole(tokenString string) (models.UserRole, error) {
	claims, err := j.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}
