package strings

import (
	"github.com/orpheus/strings/pkg/infra/postgres"
)

type StringRepository struct {
	DB postgres.PgxConn
}
