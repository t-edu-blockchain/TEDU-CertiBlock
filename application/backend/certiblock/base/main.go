package base

import (
	"certiblock/configurations"
	"database/sql"
)

type ApplicationContext struct {
	DB     *sql.DB
	Config *configurations.Configurations
}
