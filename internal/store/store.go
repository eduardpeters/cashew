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

	v := key.GetValue()
	k, ok := v.(string)
	if !ok {
		return fmt.Errorf("key value is not string: %v", v)
	}

	s.store[k] = StoredValue{value: value, expires: false}

	return nil
}

func (s *Store) SetWithExpiryMillis(key, value resp.BulkString, expiryMillis int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	v := key.GetValue()
	k, ok := v.(string)
	if !ok {
		return fmt.Errorf("key value is not string: %v", v)
	}

	expiry := time.Now().Add(time.Millisecond * time.Duration(expiryMillis))

	s.store[k] = StoredValue{value: value, expiry: expiry, expires: true}

	return nil
}

func (s *Store) Get(key resp.BulkString) (resp.CashewValue, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v := key.GetValue()
	k, ok := v.(string)
	if !ok {
		return resp.BulkString{}, fmt.Errorf("key value is not string: %v", v)
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
