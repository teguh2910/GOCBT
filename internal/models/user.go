package models

import (
	"database/sql"
	"time"
)

// UserRole represents user roles in the system
type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Never expose password hash in JSON
	FirstName    string    `json:"first_name" db:"first_name"`
	LastName     string    `json:"last_name" db:"last_name"`
	Role         UserRole  `json:"role" db:"role"`
	IsActive     bool      `json:"is_active" db:"is_active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id int) error
	List(limit, offset int) ([]*User, error)
	GetByRole(role UserRole, limit, offset int) ([]*User, error)
}

// UserService defines the interface for user business logic
type UserService interface {
	Register(username, email, password, firstName, lastName string, role UserRole) (*User, error)
	Login(username, password string) (*User, error)
	GetProfile(userID int) (*User, error)
	UpdateProfile(userID int, firstName, lastName, email string) (*User, error)
	ChangePassword(userID int, oldPassword, newPassword string) error
	ListUsers(role UserRole, limit, offset int) ([]*User, error)
	ActivateUser(userID int) error
	DeactivateUser(userID int) error
}

// FullName returns the user's full name
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsValidRole checks if the role is valid
func (r UserRole) IsValid() bool {
	switch r {
	case RoleStudent, RoleTeacher, RoleAdmin:
		return true
	default:
		return false
	}
}

// CanManageTests checks if user can manage tests
func (u *User) CanManageTests() bool {
	return u.Role == RoleTeacher || u.Role == RoleAdmin
}

// CanManageUsers checks if user can manage other users
func (u *User) CanManageUsers() bool {
	return u.Role == RoleAdmin
}

// ScanUser scans database row into User struct
func ScanUser(row interface {
	Scan(dest ...interface{}) error
}) (*User, error) {
	user := &User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
