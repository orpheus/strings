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
	RestoreStringByStringId(stringId uuid.UUID) error
	ActivateStringByStringId(stringId uuid.UUID) error
	DeactivateStringByStringId(stringId uuid.UUID) error
	PrivateStringByStringId(stringId uuid.UUID) error
	PublicStringByStringId(stringId uuid.UUID) error
}

func (s StringService) ArchiveString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}

func (s StringService) RestoreString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}

func (s StringService) ActivateString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}

func (s StringService) DeactivateString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}

func (s StringService) PrivateString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}

func (s StringService) PublicString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}

func (s StringService) DeleteString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}
