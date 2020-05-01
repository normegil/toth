package versions

import (
	"database/sql"
	"fmt"
	"github.com/normegil/postgres"
)

type SchemaCreation struct {
	DB *sql.DB
}

func (_ SchemaCreation) Number() int {
	return 1
}

func (v SchemaCreation) Upgrade() error {
	row := v.DB.QueryRow(`SELECT pg_catalog.pg_get_userbyid(d.datdba) as "Owner" FROM pg_catalog.pg_database d WHERE d.datname = current_database();`)
	var owner string
	if err := row.Scan(&owner); nil != err {
		return fmt.Errorf("load database owner: %w", err)
	}

	tableExistence := `SELECT EXISTS ( SELECT 1 FROM information_schema.tables WHERE table_name = '%s');`
	tableSetOwner := `ALTER TABLE %s OWNER TO $1;`

	userTableName := "toth_user"
	err := postgres.CreateTable(v.DB, postgres.TableInfos{
		Queries: map[string]string{
			"Table-Existence": fmt.Sprintf(tableExistence, userTableName),
			"Table-Create": `CREATE TABLE toth_user (
				id uuid primary key,
				name varchar(300),
				surname varchar(300),
				mail varchar(300) unique
				CONSTRAINT valid_mail CHECK (mail ~* '^[A-Za-z0-9._%-]+@[A-Za-z0-9.-]+[.][A-Za-z]+$'),
				hash bytea,
				algorithmID uuid)`,
			"Table-Set-Owner": fmt.Sprintf(tableSetOwner, userTableName),
		},
		Owner: owner,
	})
	if err != nil {
		return fmt.Errorf("creating table '%s': %w", userTableName, err)
	}

	return nil
}

func (v SchemaCreation) Rollback() error {
	dropTableQuery := "DROP TABLE %s"

	userTableName := "toth_user"
	if _, err := v.DB.Exec(fmt.Sprintf(dropTableQuery, userTableName)); nil != err {
		return fmt.Errorf("drop table '%s': %w", userTableName, err)
	}
	return nil
}
