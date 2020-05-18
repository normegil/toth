package internal

import (
	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID
	Name            string
	Mail            Mail
	PasswordHash    []byte
	HashAlgorithmID uuid.UUID
}

func UserAnonymous() User {
	return User{
		ID:              uuid.Nil,
		Name:            "anonymous",
		Mail:            "anonymous@toth.org",
		PasswordHash:    []byte(""),
		HashAlgorithmID: uuid.Nil,
	}
}
