package models

import "errors"

var (
	ErrModelNotFound        = errors.New("ai model not found")
	ErrInvalidRequestStatus = errors.New("invalid request status")
)
