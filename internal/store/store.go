package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/eduardpeters/cashew/internal/resp"
)

type StoredValue struct {
	value   resp.CashewValue
	expiry  time.Time
	expires bool
}

type Store struct {
	mu    sync.RWMutex
	store map[string]StoredValue
}

func NewStore() *Store {
	return &Store{store: map[string]StoredValue{}}
}

func (s *Store) Set(key, value resp.BulkString) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k, err := extractKeyString(key)
	if err != nil {
		return err
	}

	s.store[k] = StoredValue{value: value, expires: false}

	return nil
}

func (s *Store) SetWithExpiry(key, value resp.BulkString, expiry time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	k, err := extractKeyString(key)
	if err != nil {
		return err
	}

	s.store[k] = StoredValue{value: value, expiry: expiry, expires: true}

	return nil
}

func (s *Store) Get(key resp.BulkString) (resp.CashewValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	k, err := extractKeyString(key)
	if err != nil {
		return resp.BulkString{}, err
	}

	stored, ok := s.store[k]
	if !ok {
		return resp.NewNull()
	}

	if stored.expires && stored.expiry.Before(time.Now()) {
		delete(s.store, k)
		return resp.NewNull()
	}

	return stored.value, nil
}

func (s *Store) Exists(key resp.BulkString) (bool, error) {
	stored, err := s.Get(key)
	if err != nil {
		return false, err
	}
	return stored.GetValue() != nil, nil
}

func (s *Store) Delete(key resp.BulkString) error {
	k, err := extractKeyString(key)
	if err != nil {
		return err
	}
	delete(s.store, k)
	return nil
}

func extractKeyString(key resp.BulkString) (string, error) {
	v := key.GetValue()
	k, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("key value is not string: %v", v)
	}
	return k, nil
}
