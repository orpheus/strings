package system

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/core"
	"github.com/orpheus/strings/infrastructure/log"
)

type ThreadInteractor struct {
	Repo   ThreadRepository
	Logger log.Logger
}

type ThreadRepository interface {
	FindAll() ([]core.Thread, error)
	CreateOne(thread core.Thread) (core.Thread, error)
	DeleteById(id uuid.UUID) error
}

func (t *ThreadInteractor) FindAll() ([]core.Thread, error) {
	return t.Repo.FindAll()
}

func (t *ThreadInteractor) CreateOne(thread core.Thread) (core.Thread, error) {
	return t.Repo.CreateOne(thread)
}

func (t *ThreadInteractor) DeleteById(id uuid.UUID) error {
	return t.Repo.DeleteById(id)
}
