package base

import (
	"CertiBlock/application/backend/certiblock/configurations"
	"database/sql"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

type ApplicationContext struct {
	DB       *sql.DB
	Config   *configurations.Configurations
	Contract *client.Contract
}
