package strings

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/pkg/infra/postgres"
)

type StringRepository struct {
	DB postgres.PgxConn
}

func (s *StringRepository) FindAll() ([]String, error) {
	return nil, nil
}

func (s *StringRepository) FindAllByThread(threadId uuid.UUID) ([]String, error) {
	return nil, nil
}

func (s *StringRepository) CreateOne(coreString String) (*String, error) {
	return nil, nil
}

func (s *StringRepository) DeleteById(id uuid.UUID) error {
	return nil
}

func (s *StringRepository) DeleteAllByThread(threadId uuid.UUID) error {
	return nil
}

func (s *StringRepository) UpdateName(stringId uuid.UUID, name string) error {
	return nil
}

func (s *StringRepository) UpdateOrder(stringOrders []struct{}) error {
	return nil
}
