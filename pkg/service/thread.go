package service

import (
	"github.com/orpheus/strings/pkg/repo/threads"
)

type ThreadInteractor struct {
	Repo ThreadRepository
}

type ThreadRepository interface {
}

func (t *ThreadInteractor) FindAll() ([]threads.Thread, error) {
	return nil, nil
}
