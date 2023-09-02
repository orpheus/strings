package stringrepo

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/core"
	"github.com/orpheus/strings/pkg/persistence/dao/stringdao"
	"time"
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
	FindByStringId(stringId uuid.UUID) (*stringdao.VersionedStringRecord, error)
	FindAllByThreadId(threadId uuid.UUID) ([]*stringdao.VersionedStringRecord, error)
	//FindAllInThreadByStringId(stringId uuid.UUID) ([]*stringdao.VersionedStringRecord, error)
}

type StringRepository struct {
	StringDao
	VersionedStringDao
}

// CreateNewString creates a new string record and a new versioned string record with version 1.
func (s *StringRepository) CreateNewString(string *core.String) (*core.String, error) {
	if string == nil {
		return nil, fmt.Errorf("failed to create new string, missing `string`")
	}

	if string.StringId == uuid.Nil {
		string.StringId = uuid.New()
	}

	newStringRecord, err := s.StringDao.Save(&stringdao.StringRecord{
		Id: string.StringId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save new string record: %s", err)
	}

	newVersionedRecord, err := s.VersionedStringDao.Save(&stringdao.VersionedStringRecord{
		Id:          uuid.New(),
		Name:        string.Name,
		Version:     1,
		StringId:    newStringRecord.Id,
		ThreadId:    string.ThreadId,
		Order:       string.Order,
		Active:      false,
		Archived:    false,
		Deleted:     false,
		DateCreated: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save new versioned string record: %s", err)
	}

	return &core.String{
		Id:          newVersionedRecord.Id,
		Name:        newVersionedRecord.Name,
		Version:     newVersionedRecord.Version,
		StringId:    newVersionedRecord.StringId,
		ThreadId:    newVersionedRecord.ThreadId,
		Order:       newVersionedRecord.Order,
		Active:      newVersionedRecord.Active,
		Archived:    newVersionedRecord.Archived,
		Deleted:     newVersionedRecord.Deleted,
		DateCreated: newVersionedRecord.DateCreated,
	}, nil
}

// CreateVersionedString creates a new versioned string record.
func (s *StringRepository) CreateVersionedString(string *core.String) (*core.String, error) {
	if string == nil {
		return nil, fmt.Errorf("string is nil")
	}

	if string.ThreadId == uuid.Nil {
		return nil, fmt.Errorf("missing thread id`")
	}

	if string.StringId == uuid.Nil {
		return nil, fmt.Errorf("missing string id")
	}

	serverString, err := s.VersionedStringDao.FindByStringId(string.StringId)
	if err != nil {
		return nil, fmt.Errorf("failed to find string by string id: %s", err)
	}

	if serverString == nil {
		return nil, fmt.Errorf("cannot update string, string not found for id %s", string.Id)
	}

	newVersionedStringRecord, err := s.VersionedStringDao.Save(&stringdao.VersionedStringRecord{
		Id:          uuid.New(),
		Name:        string.Name,
		Version:     serverString.Version + 1,
		StringId:    serverString.StringId,
		ThreadId:    serverString.ThreadId,
		Order:       string.Order,
		Active:      string.Active,
		Archived:    string.Archived,
		Deleted:     string.Deleted,
		DateCreated: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to save new versioned string record: %s", err)
	}

	return newVersionedStringRecord.ToString(), nil
}

func (s *StringRepository) FindAllByThreadId(threadId uuid.UUID) ([]*core.String, error) {
	versionedStrings, err := s.VersionedStringDao.FindAllByThreadId(threadId)
	if err != nil {
		return nil, fmt.Errorf("failed to find strings by thread id: %s", err)
	}

	return convertVersionedStringsToCoreStrings(versionedStrings), nil
}

var ErrStringNotFound = fmt.Errorf("string not found")
var ErrStringAlreadyDeleted = fmt.Errorf("string already deleted")

func (s *StringRepository) DeleteStringByStringId(stringId uuid.UUID) error {
	serverString, err := s.VersionedStringDao.FindByStringId(stringId)
	if err != nil {
		return err
	}

	if serverString == nil {
		return ErrStringNotFound
	}

	if serverString.Deleted {
		return ErrStringAlreadyDeleted
	}

	versionedStrings, err := s.VersionedStringDao.FindAllByThreadId(serverString.ThreadId)
	if err != nil {
		return err
	}

	deletedString := false
	for _, versionedString := range versionedStrings {

		// delete string
		if versionedString.StringId == stringId {
			_, err := s.VersionedStringDao.Save(&stringdao.VersionedStringRecord{
				Id:          uuid.New(),
				Name:        versionedString.Name,
				Version:     versionedString.Version + 1,
				StringId:    versionedString.StringId,
				ThreadId:    versionedString.ThreadId,
				Order:       0,
				Active:      versionedString.Active,
				Archived:    versionedString.Archived,
				Deleted:     true,
				DateCreated: time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to delete string: %s", err)
			}
			deletedString = true
		} else if deletedString {

			// re-order proceeding strings
			_, err := s.VersionedStringDao.Save(&stringdao.VersionedStringRecord{
				Id:          uuid.New(),
				Name:        versionedString.Name,
				Version:     versionedString.Version + 1,
				StringId:    versionedString.StringId,
				ThreadId:    versionedString.ThreadId,
				Order:       versionedString.Order - 1,
				Active:      versionedString.Active,
				Archived:    versionedString.Archived,
				Deleted:     versionedString.Deleted,
				DateCreated: time.Now(),
			})
			if err != nil {
				return fmt.Errorf("failed to re-order string after delete: %s", err)
			}
		}
	}

	return nil
}

func convertVersionedStringsToCoreStrings(versionedStrings []*stringdao.VersionedStringRecord) []*core.String {
	coreStrings := make([]*core.String, len(versionedStrings))

	for i, versionedString := range versionedStrings {
		coreStrings[i] = versionedString.ToString()
	}

	return coreStrings
}
