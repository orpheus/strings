package threads

import (
	"github.com/orpheus/strings/pkg/infra/postgres"
)

type Repository struct {
	DB postgres.PgxConn
}
