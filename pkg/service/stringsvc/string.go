package stringsvc

import (
	"github.com/google/uuid"
)

func NewStringService(stringRepository StringRepository) *StringService {
	return &StringService{
		StringRepository: stringRepository,
	}
}

type StringService struct {
	StringRepository StringRepository
}

type StringRepository interface {
	DeleteStringByStringId(stringId uuid.UUID) error
	ArchiveStringByStringId(stringId uuid.UUID) error
	RestoreStringByStringId(stringId uuid.UUID) error
	ActivateStringByStringId(stringId uuid.UUID) error
	DeactivateStringByStringId(stringId uuid.UUID) error
	PrivateStringByStringId(stringId uuid.UUID) error
	PublicStringByStringId(stringId uuid.UUID) error
}

func (s StringService) ArchiveString(stringId uuid.UUID) error {
	return s.StringRepository.ArchiveStringByStringId(stringId)
}

func (s StringService) RestoreString(stringId uuid.UUID) error {
	return s.StringRepository.RestoreStringByStringId(stringId)
}

func (s StringService) ActivateString(stringId uuid.UUID) error {
	return s.StringRepository.ActivateStringByStringId(stringId)
}

func (s StringService) DeactivateString(stringId uuid.UUID) error {
	return s.StringRepository.DeactivateStringByStringId(stringId)
}

func (s StringService) PrivateString(stringId uuid.UUID) error {
	return s.StringRepository.PrivateStringByStringId(stringId)
}

func (s StringService) PublicString(stringId uuid.UUID) error {
	return s.StringRepository.PublicStringByStringId(stringId)
}

func (s StringService) DeleteString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}
