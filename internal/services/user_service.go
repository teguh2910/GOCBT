package services

import (
	"database/sql"
	"gocbt/internal/auth"
	"gocbt/internal/models"
	"strings"
)

// UserService implements the models.UserService interface
type UserService struct {
	userRepo        models.UserRepository
	passwordManager *auth.PasswordManager
}

// NewUserService creates a new user service
func NewUserService(userRepo models.UserRepository, passwordManager *auth.PasswordManager) models.UserService {
	return &UserService{
		userRepo:        userRepo,
		passwordManager: passwordManager,
	}
}

// Register creates a new user account
func (s *UserService) Register(username, email, password, firstName, lastName string, role models.UserRole) (*models.User, error) {
	// Validate input
	if strings.TrimSpace(username) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if strings.TrimSpace(email) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if strings.TrimSpace(firstName) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if strings.TrimSpace(lastName) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if !role.IsValid() {
		return nil, auth.ErrInvalidCredentials
	}

	// Validate password strength
	if err := s.passwordManager.ValidatePasswordStrength(password); err != nil {
		return nil, err
	}

	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, auth.ErrUsernameExists
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, auth.ErrEmailExists
	}

	// Hash password
	hashedPassword, err := s.passwordManager.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Username:     strings.TrimSpace(username),
		Email:        strings.TrimSpace(email),
		PasswordHash: hashedPassword,
		FirstName:    strings.TrimSpace(firstName),
		LastName:     strings.TrimSpace(lastName),
		Role:         role,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns user data
func (s *UserService) Login(username, password string) (*models.User, error) {
	// Get user by username
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrInvalidCredentials
		}
		return nil, err
	}

	if user == nil {
		return nil, auth.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, auth.ErrUserNotActive
	}

	// Verify password
	if !s.passwordManager.VerifyPassword(password, user.PasswordHash) {
		return nil, auth.ErrInvalidCredentials
	}

	return user, nil
}

// GetProfile retrieves user profile by ID
func (s *UserService) GetProfile(userID int) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, auth.ErrUserNotFound
		}
		return nil, err
	}

	if user == nil {
		return nil, auth.ErrUserNotFound
	}

	return user, nil
}

// UpdateProfile updates user profile information
func (s *UserService) UpdateProfile(userID int, firstName, lastName, email string) (*models.User, error) {
	// Get existing user
	user, err := s.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	// Validate input
	if strings.TrimSpace(firstName) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if strings.TrimSpace(lastName) == "" {
		return nil, auth.ErrInvalidCredentials
	}
	if strings.TrimSpace(email) == "" {
		return nil, auth.ErrInvalidCredentials
	}

	// Check if email is being changed and if new email already exists
	if user.Email != email {
		existingUser, err := s.userRepo.GetByEmail(email)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if existingUser != nil {
			return nil, auth.ErrEmailExists
		}
	}

	// Update user fields
	user.FirstName = strings.TrimSpace(firstName)
	user.LastName = strings.TrimSpace(lastName)
	user.Email = strings.TrimSpace(email)

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(userID int, oldPassword, newPassword string) error {
	// Get user
	user, err := s.GetProfile(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !s.passwordManager.VerifyPassword(oldPassword, user.PasswordHash) {
		return auth.ErrInvalidCredentials
	}

	// Validate new password strength
	if err := s.passwordManager.ValidatePasswordStrength(newPassword); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := s.passwordManager.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

// ListUsers retrieves a list of users by role
func (s *UserService) ListUsers(role models.UserRole, limit, offset int) ([]*models.User, error) {
	if role == "" {
		return s.userRepo.List(limit, offset)
	}
	return s.userRepo.GetByRole(role, limit, offset)
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(userID int) error {
	user, err := s.GetProfile(userID)
	if err != nil {
		return err
	}

	user.IsActive = true
	return s.userRepo.Update(user)
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(userID int) error {
	user, err := s.GetProfile(userID)
	if err != nil {
		return err
	}

	user.IsActive = false
	return s.userRepo.Update(user)
}
