package store

import (
	"fmt"
	"strconv"
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
	k, err := extractKeyString(key)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.store[k] = StoredValue{value: value, expires: false}
	s.mu.Unlock()

	return nil
}

func (s *Store) SetWithExpiry(key, value resp.BulkString, expiry time.Time) error {
	k, err := extractKeyString(key)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.store[k] = StoredValue{value: value, expiry: expiry, expires: true}
	s.mu.Unlock()

	return nil
}

func (s *Store) Get(key resp.BulkString) (resp.CashewValue, error) {
	k, err := extractKeyString(key)
	if err != nil {
		return resp.BulkString{}, err
	}

	s.mu.RLock()
	stored, ok := s.store[k]
	s.mu.RUnlock()

	if !ok {
		return resp.NewNull()
	}

	if stored.expires && stored.expiry.Before(time.Now()) {
		// Pass the the observed expired value to check we delete this same value, not another assignment
		observedExpiry := stored.expiry
		go func(key string, exp time.Time) {
			s.mu.Lock()
			defer s.mu.Unlock()

			current, ok := s.store[key]
			if ok && current.expires && current.expiry.Equal(exp) {
				s.deleteKey(k)
			}
		}(k, observedExpiry)

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

	s.mu.Lock()
	s.deleteKey(k)
	s.mu.Unlock()

	return nil
}

func (s *Store) Add(key resp.BulkString, qty int64) (resp.Integer, error) {
	stored, err := s.Get(key)
	if err != nil {
		return resp.Integer{}, err
	}

	current, ok := stored.(resp.BulkString)
	if !ok {
		defaultValue, err := resp.NewBulkString("0")
		if err != nil {
			return resp.Integer{}, err
		}
		err = s.Set(key, defaultValue)
		if err != nil {
			return resp.Integer{}, err
		}
		return s.Add(key, qty)
	}

	currentInteger, err := parseBulkStringInteger(current)
	if err != nil {
		return resp.Integer{}, err
	}

	nextValueString := strconv.Itoa(int(currentInteger + qty))

	k, err := extractKeyString(key)
	newValue, err := resp.NewBulkString(nextValueString)
	if err != nil {
		return resp.Integer{}, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	currentStored := s.store[k]
	if currentStored.value.GetValue() != stored.GetValue() {
		currentInteger, err := parseBulkStringInteger(current)
		if err != nil {
			return resp.Integer{}, err
		}
		nextValueString := strconv.Itoa(int(currentInteger + qty))
		newValue, err = resp.NewBulkString(nextValueString)
		if err != nil {
			return resp.Integer{}, err
		}
	}
	s.store[k] = StoredValue{value: newValue, expires: currentStored.expires, expiry: currentStored.expiry}

	return resp.NewInteger(nextValueString)
}

func (s *Store) Prepend(key resp.BulkString, elements ...resp.CashewValue) (resp.Integer, error) {
	k, err := extractKeyString(key)
	if err != nil {
		return resp.Integer{}, err
	}

	stored, err := s.Get(key)
	if err != nil {
		return resp.Integer{}, err
	}

	current, ok := stored.(resp.Array)
	if !ok {
		_, isEmptyKey := stored.(resp.Null)
		if !isEmptyKey {
			return resp.Integer{}, fmt.Errorf("stored value is not empty: %v", stored)
		}
	}

	currentValue := current.GetValue()
	storedArray, ok := currentValue.([]resp.CashewValue)
	if !ok {
		return resp.Integer{}, fmt.Errorf("inner value is not array: %v", currentValue)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	newArray := make([]resp.CashewValue, len(elements)+len(storedArray))
	i := 0
	for j := len(elements) - 1; j >= 0; j-- {
		newArray[i] = elements[j]
		i++
	}
	for _, se := range storedArray {
		newArray[i] = se
		i++
	}

	a, err := resp.NewArray(newArray)
	if err != nil {
		return resp.Integer{}, err
	}

	s.store[k] = StoredValue{value: a, expires: false}

	finalLength := strconv.Itoa(int(len(newArray)))

	return resp.NewInteger(finalLength)

}

// For use only with lock in calling context
func (s *Store) deleteKey(key string) {
	delete(s.store, key)
}

func extractKeyString(key resp.BulkString) (string, error) {
	v := key.GetValue()
	k, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("key value is not string: %v", v)
	}
	return k, nil
}

func parseBulkStringInteger(value resp.CashewValue) (int64, error) {
	v := value.GetValue()
	s, ok := v.(string)
	if !ok {
		return 0, fmt.Errorf("inner value is not string: %v", v)
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("value not integer representation: %s", s)
	}
	return n, nil
}
