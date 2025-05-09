// cli-robo/connections/database_interface.go
package connections

import "database/sql"

type Database interface {
	DB() (*sql.DB, error)
}
