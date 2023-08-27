package stringrepo

import (
	"github.com/orpheus/strings/pkg/core"
	"github.com/orpheus/strings/pkg/persistence/dao/stringdao"
)

func NewStringRepository(stringDao StringDao, versionedStringDao VersionedStringDao) *StringRepository {
	return &StringRepository{
		StringDao:          stringDao,
		VersionedStringDao: versionedStringDao,
	}
}

type StringDao interface {
	Save(record *stringdao.StringRecord) (*stringdao.StringRecord, error)
}

type VersionedStringDao interface {
	Save(record *stringdao.VersionedStringRecord) (*stringdao.VersionedStringRecord, error)
}

type StringRepository struct {
	StringDao
	VersionedStringDao
}

// CreateString creates a new string record and a version 1 string record
func (s *StringRepository) CreateString(string *core.String) (*core.String, error) {
	//TODO implement me
	panic("implement me")
}

// UpdateString creates a new versioned string record.
func (s *StringRepository) UpdateString(string *core.String) (*core.String, error) {
	//TODO implement me
	panic("implement me")
}
