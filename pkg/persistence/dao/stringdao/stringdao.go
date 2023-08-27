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

func (s StringDao) Save(record *StringRecord) (*StringRecord, error) {
	//TODO implement me
	panic("implement me")
}
