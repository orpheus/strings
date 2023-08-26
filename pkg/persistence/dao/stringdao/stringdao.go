package stringdao

import (
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

type StringRecord struct {
	Id uuid.UUID
}

type StringDao struct {
	*sqldb.Store
}
