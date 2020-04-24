package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/euvsvirus-banan/backend/users/rpc/userspb"
)

var (
	ErrDuplicate = errors.New("already exists")
	ErrNotFound  = errors.New("not found error")
)

type UserStorage struct {
	wr   io.WriteSeeker
	data map[string]*userspb.User
}

func NewUserStorage(wr io.WriteSeeker, data map[string]*userspb.User) *UserStorage {
	return &UserStorage{
		wr:   wr,
		data: data,
	}
}

func dump(wr io.WriteSeeker, data map[string]*userspb.User) error {
	if _, err := wr.Seek(0, 0); err != nil {
		return fmt.Errorf("problem rewinding file: %w", err)
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("problem marshaling data: %w", err)
	}
	if _, err = wr.Write(b); err != nil {
		return fmt.Errorf("problem saving data: %w", err)
	}
	return nil
}

func (s *UserStorage) Add(id string, user *userspb.User) error {
	_, ok := s.data[id]
	if ok {
		return ErrDuplicate
	}
	s.data[id] = user
	return dump(s.wr, s.data)
}

func (s *UserStorage) Update(id string, element *userspb.User) error {
	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	s.data[id] = element
	return dump(s.wr, s.data)
}

func (s *UserStorage) Delete(id string) error {
	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.data, id)
	return dump(s.wr, s.data)
}

func (s *UserStorage) Get(id string) (*userspb.User, error) {
	e, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *UserStorage) All() map[string]*userspb.User {
	return s.data
}
