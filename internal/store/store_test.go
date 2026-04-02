package store_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

func TestSetNewValue(t *testing.T) {
	key := mustNewBulkString(t, "name")
	value := mustNewBulkString(t, "john")

	s := store.NewStore()

	err := s.Set(key, value)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	stored, err := s.Get(key)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if stored.GetValue() != value.GetValue() {
		t.Errorf("incorrect stored value want %q, got %q", value.GetValue(), stored.GetValue())
	}
}

func TestSetWithExpiry(t *testing.T) {
	key := mustNewBulkString(t, "name")
	value := mustNewBulkString(t, "john")
	expiry := time.Now().Add(time.Millisecond * 5)

	s := store.NewStore()

	err := s.SetWithExpiry(key, value, expiry)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	stored, err := s.Get(key)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if stored.GetValue() != value.GetValue() {
		t.Errorf("incorrect stored value want %q, got %q", value.GetValue(), stored.GetValue())
	}
}

func TestGetMissingValue(t *testing.T) {
	s := store.NewStore()

	stored, err := s.Get(mustNewBulkString(t, "not:there"))
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if stored.GetValue() != nil {
		t.Errorf("incorrect stored value want %v, got %v", nil, stored.GetValue())
	}
}

func TestGetValue(t *testing.T) {
	s := store.NewStore()
	key := "name"
	value := "john"
	mustSetInStore(t, s, key, value)

	stored, err := s.Get(mustNewBulkString(t, key))
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got := stored.GetValue()
	if got != value {
		t.Errorf("incorrect stored value want %q, got %q", value, got)
	}
}

func TestGetValueWithinExpiry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {

		s := store.NewStore()
		key := "name"
		value := "john"
		mustSetInStoreWithExpiryMillis(t, s, key, value, 5000)

		time.Sleep(time.Second * 2)

		stored, err := s.Get(mustNewBulkString(t, key))
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		got := stored.GetValue()
		if got != value {
			t.Errorf("incorrect stored value want %q, got %q", value, got)
		}
	})
}

func TestGetValueAfterExpiry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := store.NewStore()
		key := "name"
		value := "john"
		mustSetInStoreWithExpiryMillis(t, s, key, value, 1000)

		time.Sleep(time.Second * 2)

		stored, err := s.Get(mustNewBulkString(t, key))
		if err != nil {
			t.Fatalf("Unexpected error %v", err)
		}

		if stored.GetValue() != nil {
			t.Errorf("incorrect stored value want %v, got %v", nil, stored.GetValue())
		}
	})
}

func TestExistsMissingValue(t *testing.T) {
	s := store.NewStore()
	key := "name"

	exists, err := s.Exists(mustNewBulkString(t, key))
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if exists {
		t.Errorf("Value should not be found, got %v", exists)
	}
}

func TestExistsStoredValue(t *testing.T) {
	s := store.NewStore()
	key := "name"
	value := "john"
	mustSetInStore(t, s, key, value)

	exists, err := s.Exists(mustNewBulkString(t, key))
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if !exists {
		t.Errorf("Value should be found, got %v", exists)
	}
}

func mustNewBulkString(t testing.TB, s string) resp.BulkString {
	t.Helper()
	v, err := resp.NewBulkString(s)
	if err != nil {
		t.Fatalf("NewBulkString(%q): %v", s, err)
	}
	return v
}

func mustSetInStore(t testing.TB, s *store.Store, key, value string) {
	t.Helper()
	err := s.Set(mustNewBulkString(t, key), mustNewBulkString(t, value))
	if err != nil {
		t.Fatalf("store.Set(%q,%q): %v", key, value, err)
	}
}

func mustSetInStoreWithExpiryMillis(t testing.TB, s *store.Store, key, value string, expiryMillis int) {
	t.Helper()
	expiry := time.Now().Add(time.Millisecond * time.Duration(expiryMillis))
	err := s.SetWithExpiry(mustNewBulkString(t, key), mustNewBulkString(t, value), expiry)
	if err != nil {
		t.Fatalf("store.Set(%q,%q): %v", key, value, err)
	}
}
