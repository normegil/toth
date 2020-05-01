package versions

import (
	"database/sql"
	"github.com/normegil/godatabaseversioner"
)

func Load(db *sql.DB) []godatabaseversioner.Version {
	return []godatabaseversioner.Version{
		godatabaseversioner.PostgresVersioning{DB: db},
		&SchemaCreation{DB: db},
	}
}
