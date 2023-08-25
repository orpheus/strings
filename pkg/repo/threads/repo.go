package threads

import (
	"github.com/google/uuid"
	"github.com/orpheus/strings/pkg/infra/sqldb"
)

type Repository struct {
	*sqldb.Store

	isAtomic bool
}

func NewThreadRepository(store *sqldb.Store) *Repository {
	return &Repository{
		Store: store,
	}
}

func (r *Repository) FindByThreadId(threadId uuid.UUID) (*Thread, error) {
	_, err := r.Store.GetExecutor(r.isAtomic)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Repository) CreateThread(thread *Thread) (*Thread, error) {
	return nil, nil
}
