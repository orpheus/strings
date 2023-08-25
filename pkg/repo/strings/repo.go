package strings

import "github.com/orpheus/strings/pkg/infra/sqldb"

type StringRepository struct {
	*sqldb.Store
}
