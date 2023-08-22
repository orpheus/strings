package threads

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/pkg/infra/postgres"
)

type Repository struct {
	DB postgres.PgxConn
}

func (r *Repository) FindAll() ([]Thread, error) {
	return nil, nil
}

func (r *Repository) CreateOne(thread Thread) (*Thread, error) {
	return nil, nil
}

func (r *Repository) DeleteById(id uuid.UUID) error {
	return nil
}
