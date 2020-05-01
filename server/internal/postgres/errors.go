package postgres

import (
	"errors"
	"fmt"
)

type NotFoundError struct {
	ID       string
	Original error
}

func (e NotFoundError) Error() string {
	if nil != e.Original {
		return fmt.Errorf("data '%s' doesn't exist: %w", e.ID, e.Original).Error()
	}
	return fmt.Errorf("data '%s' doesn't exist", e.ID).Error()
}

func IsNotFoundError(err error) bool {
	return errors.As(err, &NotFoundError{})
}
