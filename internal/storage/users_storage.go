package storage // nolint: dupl

import (
	"io"

	"github.com/euvsvirus-banan/backend/users/rpc/userspb"
)

type UsersStorage struct {
	wr   io.WriteSeeker
	data map[string]*userspb.User
}

func NewUsersStorage(wr io.WriteSeeker, data map[string]*userspb.User) *UsersStorage {
	return &UsersStorage{
		wr:   wr,
		data: data,
	}
}

func (s *UsersStorage) Add(id string, user *userspb.User) error {
	_, ok := s.data[id]
	if ok {
		return ErrDuplicate
	}
	s.data[id] = user
	return dump(s.wr, s.data)
}

func (s *UsersStorage) Update(id string, element *userspb.User) error {
	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	s.data[id] = element
	return dump(s.wr, s.data)
}

func (s *UsersStorage) Delete(id string) error {
	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.data, id)
	return dump(s.wr, s.data)
}

func (s *UsersStorage) Get(id string) (*userspb.User, error) {
	e, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *UsersStorage) All() map[string]*userspb.User {
	return s.data
}
