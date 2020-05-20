package internal

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type User struct {
	ID              uuid.UUID
	Name            string
	Mail            Mail
	PasswordHash    []byte
	HashAlgorithmID uuid.UUID
}

func (u User) toParseableUser() parseableUser {
	return parseableUser{
		ID:              u.ID.String(),
		HashAlgorithmID: u.HashAlgorithmID.String(),
		Name:            u.Name,
		Mail:            string(u.Mail),
		PasswordHash:    base64.StdEncoding.EncodeToString(u.PasswordHash),
	}
}

func (u User) MarshalYAML() (interface{}, error) {
	return u.toParseableUser(), nil
}

func (u *User) UnmarshalYAML(value *yaml.Node) error {
	decoded := parseableUser{}
	if err := value.Decode(&decoded); nil != err {
		return err
	}

	user, err := decoded.ToUser()
	if err != nil {
		return err
	}

	u.ID = user.ID
	u.Name = user.Name
	u.Mail = user.Mail
	u.PasswordHash = user.PasswordHash
	u.HashAlgorithmID = user.HashAlgorithmID

	return nil
}

type parseableUser struct {
	ID              string `yaml:"id,omitempty"`
	Name            string `yaml:"name,omitempty"`
	Mail            string `yaml:"mail,omitempty"`
	PasswordHash    string `yaml:"passwordhash,omitempty"`
	HashAlgorithmID string `yaml:"algorithmID,omitempty"`
}

func (u parseableUser) ToUser() (*User, error) {
	parsedID, err := uuid.Parse(u.ID)
	if err != nil {
		return nil, fmt.Errorf("could not parse ID '%s': %w", u.ID, err)
	}

	parsedAlgorithmID, err := uuid.Parse(u.HashAlgorithmID)
	if err != nil {
		return nil, fmt.Errorf("could not parse Algorithm ID '%s': %w", u.HashAlgorithmID, err)
	}

	mail, err := NewMail(u.Mail)
	if err != nil {
		return nil, err
	}

	passwordHash, err := base64.StdEncoding.DecodeString(u.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:              parsedID,
		Name:            u.Name,
		Mail:            mail,
		PasswordHash:    passwordHash,
		HashAlgorithmID: parsedAlgorithmID,
	}, nil
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
