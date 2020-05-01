package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/normegil/toth/server/internal"
	"strings"
)

type UserDAO struct {
	Querier Querier
}

func (d UserDAO) Load(id uuid.UUID) (*internal.User, error) {
	row := d.Querier.QueryRow(`SELECT name, surname, mail, hash, algorithmID FROM "user" WHERE id=$1`, id.String())
	return d.toUser(id.String(), row)
}

func (d UserDAO) LoadByMail(mail internal.Mail) (*internal.User, error) {
	row := d.Querier.QueryRow(`SELECT id, name, surname, mail, hash, algorithmID FROM "user" WHERE mail=$1`, string(mail))
	return d.toUser(string(mail), row)
}

func (d UserDAO) toUser(queryID string, row *sql.Row) (*internal.User, error) {
	var idStr string
	var name string
	var surname string
	var mail string
	var hash []byte
	var algorithmID uuid.UUID

	if err := row.Scan(&idStr, &name, &surname, &mail, &hash, &algorithmID); nil != err {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, NotFoundError{
				ID:       queryID,
				Original: err,
			}
		}
		return nil, fmt.Errorf("loading user '%s': %w", queryID, err)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("User %s has wrong id %s: %w", queryID, idStr, err)
	}

	return &internal.User{
		ID:              id,
		Name:            name,
		Surname:         surname,
		Mail:            internal.Mail(mail),
		PasswordHash:    hash,
		HashAlgorithmID: algorithmID,
	}, nil
}

func (d UserDAO) Insert(user internal.User) error {
	if _, err := d.Querier.Exec(`INSERT INTO "user" (id, name, surname, mail, hash, algorithmID) VALUES (gen_random_uuid(), $1, $2, $3, $4::bytea, $5);`, user.Name, user.Surname, string(user.Mail), user.PasswordHash, user.HashAlgorithmID); err != nil {
		return fmt.Errorf("inserting %s %s (%s): %w", user.Name, user.Surname, user.Mail, err)
	}
	return nil
}

func (d UserDAO) IsNotFoundError(err error) bool {
	return IsNotFoundError(err)
}
