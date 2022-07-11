package system

import (
	"github.com/gofrs/uuid"
	"github.com/orpheus/strings/core"
	"github.com/orpheus/strings/infrastructure/logging"
)

type ThreadInteractor struct {
	Repo          ThreadRepository
	StringDeleter StringDeleter
	Logger        logging.Logger
}

type ThreadRepository interface {
	FindAll() ([]core.Thread, error)
	CreateOne(thread core.Thread) (core.Thread, error)
	DeleteById(id uuid.UUID) error
}

type StringDeleter interface {
	DeleteAllByThread(threadId uuid.UUID) error
}

func (t *ThreadInteractor) FindAll() ([]core.Thread, error) {
	return t.Repo.FindAll()
}

func (t *ThreadInteractor) CreateOne(thread core.Thread) (core.Thread, error) {
	return t.Repo.CreateOne(thread)
}

// DeleteById deletes all the strings associated with a thread and then deletes
// the thread itself.
func (t *ThreadInteractor) DeleteById(id uuid.UUID) error {
	err := t.StringDeleter.DeleteAllByThread(id)
	if err != nil {
		return err
	}
	return t.Repo.DeleteById(id)
}
