package storage // nolint: dupl

import (
	"io"

	"github.com/euvsvirus-banan/backend/requests/rpc/requestspb"
)

type RequestsStorage struct {
	wr   io.WriteSeeker
	data map[string]*requestspb.Request
}

func NewRequestsStorage(wr io.WriteSeeker, data map[string]*requestspb.Request) *RequestsStorage {
	return &RequestsStorage{
		wr:   wr,
		data: data,
	}
}

func (s *RequestsStorage) Add(id string, request *requestspb.Request) error {
	_, ok := s.data[id]
	if ok {
		return ErrDuplicate
	}
	s.data[id] = request
	return dump(s.wr, s.data)
}

func (s *RequestsStorage) Update(id string, element *requestspb.Request) error {
	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	s.data[id] = element
	return dump(s.wr, s.data)
}

func (s *RequestsStorage) Delete(id string) error {
	_, ok := s.data[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.data, id)
	return dump(s.wr, s.data)
}

func (s *RequestsStorage) Get(id string) (*requestspb.Request, error) {
	e, ok := s.data[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *RequestsStorage) All() map[string]*requestspb.Request {
	return s.data
}
