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

func (d UserDAO) LoadAll() ([]internal.User, error) {
	rows, err := d.Querier.Query(`SELECT id, name, mail, hash, algorithmID FROM toth_user`)
	if err != nil {
		return nil, fmt.Errorf("querying all users: %w", err)
	}
	users := make([]internal.User, 0)
	for rows.Next() {
		var idStr string
		var name string
		var mailStr string
		var hash []byte
		var algorithmIDStr string
		if err := rows.Scan(&idStr, &name, &mailStr, &hash, &algorithmIDStr); nil != err {
			return nil, fmt.Errorf("scanning rows: %w", err)
		}

		id, err := uuid.Parse(idStr)
		if err != nil {
			return nil, fmt.Errorf("user '%s' has wrong id %s: %w", name, idStr, err)
		}

		algorithmID, err := uuid.Parse(algorithmIDStr)
		if err != nil {
			return nil, fmt.Errorf("user '%s' has wrong algorithm id %s: %w", id.String(), algorithmIDStr, err)
		}

		mail, err := internal.NewMail(mailStr)
		if err != nil {
			return nil, fmt.Errorf("invalid mail '%s': %w", mailStr, err)
		}

		users = append(users, internal.User{
			ID:              id,
			Name:            name,
			Mail:            mail,
			PasswordHash:    hash,
			HashAlgorithmID: algorithmID,
		})
	}
	return users, nil
}

func (d UserDAO) Load(id uuid.UUID) (*internal.User, error) {
	row := d.Querier.QueryRow(`SELECT id, name, mail, hash, algorithmID FROM toth_user WHERE id=$1`, id.String())
	return d.toUser(id.String(), row)
}

func (d UserDAO) LoadByMail(mail internal.Mail) (*internal.User, error) {
	row := d.Querier.QueryRow(`SELECT id, name, mail, hash, algorithmID FROM toth_user WHERE mail=$1`, string(mail))
	return d.toUser(string(mail), row)
}

func (d UserDAO) toUser(queryID string, row *sql.Row) (*internal.User, error) {
	var idStr string
	var name string
	var mail string
	var hash []byte
	var algorithmIDStr string

	if err := row.Scan(&idStr, &name, &mail, &hash, &algorithmIDStr); nil != err {
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
		return nil, fmt.Errorf("user %s has wrong id %s: %w", queryID, idStr, err)
	}

	algorithmID, err := uuid.Parse(algorithmIDStr)
	if err != nil {
		return nil, fmt.Errorf("user %s has wrong algorithm id %s: %w", queryID, algorithmIDStr, err)
	}

	return &internal.User{
		ID:              id,
		Name:            name,
		Mail:            internal.Mail(mail),
		PasswordHash:    hash,
		HashAlgorithmID: algorithmID,
	}, nil
}

func (d UserDAO) Insert(user internal.User) error {
	if _, err := d.Querier.Exec(`INSERT INTO toth_user (id, name, mail, hash, algorithmID) VALUES (gen_random_uuid(), $1, $2, $3::bytea, $4);`, user.Name, string(user.Mail), user.PasswordHash, user.HashAlgorithmID); err != nil {
		return fmt.Errorf("inserting %s: %w", user.Mail, err)
	}
	return nil
}

func (d UserDAO) IsNotFoundError(err error) bool {
	return IsNotFoundError(err)
}
