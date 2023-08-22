package service

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/core"
	"github.com/orpheus/strings/infrastructure/logging"
)

type StringInteractor struct {
	StringRepository StringRepository
	Logger           logging.Logger
}

type StringRepository interface {
	FindAll() ([]core.String, error)
	FindAllByThread(threadId uuid.UUID) ([]core.String, error)
	CreateOne(core.String) (core.String, error)
	DeleteById(id uuid.UUID) error
	UpdateName(stringId uuid.UUID, name string) error
	UpdateOrder(stringOrders []core.StringOrder) error
}

func (s *StringInteractor) FindAll() ([]core.String, error) {
	return s.StringRepository.FindAll()
}

func (s *StringInteractor) FindAllByThread(threadId uuid.UUID) ([]core.String, error) {
	return s.StringRepository.FindAllByThread(threadId)
}

func (s *StringInteractor) CreateOne(string core.String) (core.String, error) {
	return s.StringRepository.CreateOne(string)
}

func (s *StringInteractor) DeleteById(id uuid.UUID) error {
	return s.StringRepository.DeleteById(id)
}

func (s *StringInteractor) UpdateName(stringId uuid.UUID, name string) error {
	return s.StringRepository.UpdateName(stringId, name)
}

func (s *StringInteractor) UpdateOrder(stringOrders []core.StringOrder) error {
	return s.StringRepository.UpdateOrder(stringOrders)
}
