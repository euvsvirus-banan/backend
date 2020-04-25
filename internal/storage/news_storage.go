package storage // nolint: dupl

import (
	"io"

	"github.com/euvsvirus-banan/backend/news/rpc/newspb"
)

type NewsStorage struct {
	wr   io.WriteSeeker
	data map[string]*newspb.News
}

func NewNewsStorage(wr io.WriteSeeker, data map[string]*newspb.News) *NewsStorage {
	return &NewsStorage{
		wr:   wr,
		data: data,
	}
}

func (s *NewsStorage) Add(id string, new *newspb.News) error {
	_, ok := s.data[id]
	if ok {
		return ErrDuplicate
	}
	s.data[id] = new
	return dump(s.wr, s.data)
}

func (s *NewsStorage) Update(id string, element *newspb.News) error {
	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	s.data[id] = element
	return dump(s.wr, s.data)
}

func (s *NewsStorage) Delete(id string) error {
	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.data, id)
	return dump(s.wr, s.data)
}

func (s *NewsStorage) Get(id string) (*newspb.News, error) {
	e, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *NewsStorage) All() map[string]*newspb.News {
	return s.data
}
