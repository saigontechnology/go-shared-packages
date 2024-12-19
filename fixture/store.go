package fixture

import (
	"sync"

	"github.com/pkg/errors"
)

var mutex = &sync.Mutex{}

type Store interface {
	HasReference(reference string) bool
	Set(reference string, fixture Model) error
	Get(reference string) (Model, error)
}

type store struct {
	fixtures map[string]Model
}

func NewStore() Store {
	return &store{
		fixtures: make(map[string]Model),
	}
}

func (s *store) Set(reference string, fixture Model) error {
	mutex.Lock()
	defer mutex.Unlock()
	if s.HasReference(reference) {
		return errors.Errorf(`fixture %q was created already`, reference)
	}
	s.fixtures[reference] = fixture

	return nil
}

func (s *store) Get(reference string) (Model, error) {
	fixture, ok := s.fixtures[reference]
	if !ok {
		return nil, errors.Errorf(`fixture %q is not created yet`, reference)
	}

	return fixture, nil
}

func (s *store) HasReference(reference string) bool {
	_, ok := s.fixtures[reference]
	return ok
}
