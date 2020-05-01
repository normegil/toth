package security

import (
	"fmt"
	"github.com/normegil/toth/server/internal"
)

type Authenticator struct {
	DAO UserDAO
}

func (a Authenticator) Authenticate(mail internal.Mail, password string) (*internal.User, error) {
	user, err := a.DAO.LoadByMail(mail)
	if err != nil {
		if a.DAO.IsNotFoundError(err) {
			return nil, userNotExistError{
				UserID:   string(mail),
				Original: err,
			}
		}
		return nil, fmt.Errorf("loading user '%s': %w", string(mail), err)
	}

	if err := AllHashAlgorithms().FindByID(user.HashAlgorithmID).Validate(user.PasswordHash, password); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}
	return user, nil
}

type UserDAO interface {
	LoadByMail(mail internal.Mail) (*internal.User, error)
	IsNotFoundError(error) bool
}
