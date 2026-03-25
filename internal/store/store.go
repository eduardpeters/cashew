package store

import (
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
)

type StoredValue struct {
	value resp.CashewValue
}

type Store struct {
	store map[string]StoredValue
}

func NewStore() *Store {
	return &Store{store: map[string]StoredValue{}}
}

func (s *Store) Set(key, value resp.BulkString) error {
	v := key.GetValue()
	k, ok := v.(string)
	if !ok {
		return fmt.Errorf("key value is not string: %v", v)
	}

	s.store[k] = StoredValue{value: value}

	return nil
}

func (s *Store) Get(key resp.BulkString) (resp.CashewValue, error) {
	v := key.GetValue()
	k, ok := v.(string)
	if !ok {
		return resp.BulkString{}, fmt.Errorf("key value is not string: %v", v)
	}
	stored, ok := s.store[k]
	if !ok {
		return resp.NewNull()
	}

	return stored.value, nil
}
