package auth

import "errors"

// Authentication errors
var (
	ErrInvalidCredentials       = errors.New("invalid credentials")
	ErrUserNotFound             = errors.New("user not found")
	ErrUserNotActive            = errors.New("user account is not active")
	ErrInvalidToken             = errors.New("invalid token")
	ErrTokenExpired             = errors.New("token has expired")
	ErrUnauthorized             = errors.New("unauthorized")
	ErrForbidden                = errors.New("forbidden")
	ErrPasswordTooShort         = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong          = errors.New("password must be no more than 128 characters long")
	ErrPasswordMissingUppercase = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingLowercase = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingDigit     = errors.New("password must contain at least one digit")
	ErrPasswordMissingSpecial   = errors.New("password must contain at least one special character")
	ErrUsernameExists           = errors.New("username already exists")
	ErrEmailExists              = errors.New("email already exists")
)
