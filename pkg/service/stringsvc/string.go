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
}

func (s StringService) ArchiveString(stringId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s StringService) RestoreString(stringId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s StringService) ActivateString(stringId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s StringService) DeactivateString(stringId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s StringService) DeleteString(stringId uuid.UUID) error {
	return s.StringRepository.DeleteStringByStringId(stringId)
}
