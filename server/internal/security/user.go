package security

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/toth/server/internal"
)

func NewUser(name string, mail internal.Mail, password string) (*internal.User, error) {
	hashAlgorithm := HashAlgorithmBcrypt14()
	hash, err := hashAlgorithm.HashAndSalt(password)
	if err != nil {
		return nil, fmt.Errorf("hashing with '%s': %w", hashAlgorithm.ID().String(), err)
	}

	return &internal.User{
		ID:              uuid.New(),
		Name:            name,
		Mail:            mail,
		PasswordHash:    hash,
		HashAlgorithmID: hashAlgorithm.ID(),
	}, nil
}
