package models

import "errors"

var (
	ErrInternalServer       = errors.New("internal server error")
	ErrNotAdmin             = errors.New("admin access required")
	ErrUnauthenticated      = errors.New("user not authenticated")
	ErrUnauthorized         = errors.New("user not authorized")
	ErrUserNotFound         = errors.New("user not found")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrInvalidEmail         = errors.New("invalid email format")
	ErrUserInactive         = errors.New("user is inactive")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrMaxLoginAttempts     = errors.New("maximum login attempts exceeded")
	ErrAPIKeyNotFound       = errors.New("api key not found")
	ErrAPIKeyInactive       = errors.New("api key is inactive")
	ErrRateLimitExceeded    = errors.New("rate limit exceeded")
	ErrModelNotFound        = errors.New("ai model not found")
	ErrModelInactive        = errors.New("ai model is inactive")
	ErrInvalidRequestStatus = errors.New("invalid request status")
)
