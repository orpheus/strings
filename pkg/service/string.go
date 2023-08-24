package service

import (
	"github.com/orpheus/strings/pkg/repo/strings"
)

type StringService struct {
	StringRepository StringRepository
}

type StringRepository interface {
}

func (s *StringService) FindAll() ([]strings.String, error) {
	return nil, nil
}
