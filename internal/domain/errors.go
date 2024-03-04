package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrUserNotFound        = errors.New("USER_NOT_FOUND")
	ErrWrongPassword       = errors.New("PASSWORD_WRONG")
	ErrInvalidToken        = errors.New("TOKEN_INVALID")
	ErrInvalidVideoName    = errors.New("INVALID_VIDEO_NAME")
)
