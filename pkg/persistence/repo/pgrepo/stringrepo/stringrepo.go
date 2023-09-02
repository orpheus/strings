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

	return newVersionedRecord.ToString(), nil
}

// CreateVersionedString creates a new versioned string record.
func (s *StringRepository) CreateVersionedString(clientString *core.String) (*core.String, error) {
	if clientString == nil {
		return nil, fmt.Errorf("string is nil")
	}

	if clientString.ThreadId == uuid.Nil {
		return nil, fmt.Errorf("missing thread id`")
	}

	if clientString.StringId == uuid.Nil {
		return nil, fmt.Errorf("missing string id")
	}

	serverString, err := s.VersionedStringDao.FindByStringId(clientString.StringId)
	if err != nil {
		return nil, fmt.Errorf("failed to find string by string id: %s", err)
	}

	if serverString == nil {
		return nil, fmt.Errorf("cannot update string, string not found for id %s", clientString.Id)
	}

	savedVersionedStringRecord, err := s.VersionedStringDao.Save(newVersionedStringRecord(serverString, func(versionedString *stringdao.VersionedStringRecord) {
		versionedString.Name = clientString.Name
		versionedString.Order = clientString.Order
	}))

	if err != nil {
		return nil, fmt.Errorf("failed to save new versioned string record: %s", err)
	}

	return savedVersionedStringRecord.ToString(), nil
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
var ErrStringAlreadyArchived = fmt.Errorf("string already archived")
var ErrStringAlreadyRestored = fmt.Errorf("string already restored")
var ErrStringAlreadyActive = fmt.Errorf("string already active")
var ErrStringAlreadyDeactivated = fmt.Errorf("string already de-active")
var ErrStringAlreadyPrivate = fmt.Errorf("string already private")
var ErrStringAlreadyPublic = fmt.Errorf("string already public")

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
			_, err = s.VersionedStringDao.Save(newVersionedStringRecord(versionedString, func(versionedString *stringdao.VersionedStringRecord) {
				versionedString.Deleted = true
				versionedString.Order = -1
			}))
			if err != nil {
				return fmt.Errorf("failed to delete string: %s", err)
			}
			deletedString = true
		} else if deletedString {

			// re-order proceeding strings
			_, err = s.VersionedStringDao.Save(newVersionedStringRecord(versionedString, func(versionedString *stringdao.VersionedStringRecord) {
				versionedString.Order = versionedString.Order - 1
			}))
			if err != nil {
				return fmt.Errorf("failed to re-order string after delete: %s", err)
			}
		}
	}

	return nil
}

func (s *StringRepository) ArchiveStringByStringId(stringId uuid.UUID) error {
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

	if serverString.Archived {
		return ErrStringAlreadyArchived
	}

	versionedStrings, err := s.VersionedStringDao.FindAllByThreadId(serverString.ThreadId)
	if err != nil {
		return err
	}

	archivedString := false
	for _, versionedString := range versionedStrings {

		// archive string
		if versionedString.StringId == stringId {
			_, err = s.VersionedStringDao.Save(newVersionedStringRecord(versionedString, func(versionedString *stringdao.VersionedStringRecord) {
				versionedString.Archived = true
				versionedString.Order = -1
			}))
			if err != nil {
				return fmt.Errorf("failed to delete string: %s", err)
			}
			archivedString = true
		} else if archivedString {

			// re-order proceeding strings
			_, err = s.VersionedStringDao.Save(newVersionedStringRecord(versionedString, func(versionedString *stringdao.VersionedStringRecord) {
				versionedString.Order = versionedString.Order - 1
			}))
			if err != nil {
				return fmt.Errorf("failed to re-order string after delete: %s", err)
			}
		}
	}

	return nil
}

func (s *StringRepository) RestoreStringByStringId(stringId uuid.UUID) error {
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

	if !serverString.Archived {
		return ErrStringAlreadyRestored
	}

	versionedStrings, err := s.VersionedStringDao.FindAllByThreadId(serverString.ThreadId)
	if err != nil {
		return err
	}

	_, err = s.VersionedStringDao.Save(newVersionedStringRecord(serverString, func(versionedString *stringdao.VersionedStringRecord) {
		versionedString.Archived = false
		versionedString.Order = len(versionedStrings) + 1
	}))
	if err != nil {
		return fmt.Errorf("failed to restore string: %s", err)
	}

	return nil
}

func (s *StringRepository) ActivateStringByStringId(stringId uuid.UUID) error {
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

	if serverString.Active {
		return ErrStringAlreadyActive
	}

	_, err = s.VersionedStringDao.Save(newVersionedStringRecord(serverString, func(versionedString *stringdao.VersionedStringRecord) {
		versionedString.Active = true
	}))
	if err != nil {
		return fmt.Errorf("failed to activate string: %s", err)
	}

	return nil
}

func (s *StringRepository) DeactivateStringByStringId(stringId uuid.UUID) error {
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

	if !serverString.Active {
		return ErrStringAlreadyDeactivated
	}

	_, err = s.VersionedStringDao.Save(newVersionedStringRecord(serverString, func(versionedString *stringdao.VersionedStringRecord) {
		versionedString.Active = false
	}))
	if err != nil {
		return fmt.Errorf("failed to deactivate string: %s", err)
	}

	return nil
}

func (s *StringRepository) PrivateStringByStringId(stringId uuid.UUID) error {
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

	if serverString.Private {
		return ErrStringAlreadyPrivate
	}

	_, err = s.VersionedStringDao.Save(newVersionedStringRecord(serverString, func(versionedString *stringdao.VersionedStringRecord) {
		versionedString.Private = true
	}))
	if err != nil {
		return fmt.Errorf("failed to make string private: %s", err)
	}

	return nil
}

func (s *StringRepository) PublicStringByStringId(stringId uuid.UUID) error {
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

	if !serverString.Private {
		return ErrStringAlreadyPublic
	}

	_, err = s.VersionedStringDao.Save(newVersionedStringRecord(serverString, func(versionedString *stringdao.VersionedStringRecord) {
		versionedString.Private = false
	}))
	if err != nil {
		return fmt.Errorf("failed to make string public: %s", err)
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

// helper to create a new versioned string record from a server string to allow for easy updates
// when the VersionedStringRecord changes.
func newVersionedStringRecord(sourceString *stringdao.VersionedStringRecord, updater func(versionedString *stringdao.VersionedStringRecord)) *stringdao.VersionedStringRecord {
	newVersionedString := &stringdao.VersionedStringRecord{
		Id:          uuid.New(),
		Name:        sourceString.Name,
		Version:     sourceString.Version + 1,
		StringId:    sourceString.StringId,
		ThreadId:    sourceString.ThreadId,
		Order:       sourceString.Order,
		Active:      sourceString.Active,
		Archived:    sourceString.Archived,
		Private:     sourceString.Private,
		Deleted:     sourceString.Deleted,
		DateCreated: time.Now(),
	}

	updater(newVersionedString)
	return newVersionedString
}
