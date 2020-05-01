package security

import (
	"errors"
	"fmt"
)

type invalidPasswordError struct {
	Original error
}

func (e invalidPasswordError) Error() string {
	return fmt.Errorf("invalid password: %w", e.Original).Error()
}

func IsInvalidPassword(err error) bool {
	return errors.As(err, &invalidPasswordError{})
}

type userNotExistError struct {
	UserID   string
	Original error
}

func (e userNotExistError) Error() string {
	return fmt.Errorf("user '%s' doesn't exist: %w", e.UserID, e.Original).Error()
}

func IsUserNotExistError(err error) bool {
	return errors.As(err, &userNotExistError{})
}

func IsInvalidAuthentication(err error) bool {
	return IsInvalidPassword(err) || IsUserNotExistError(err)
}
