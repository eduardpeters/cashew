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

func TestDeleteMissingValue(t *testing.T) {
	s := store.NewStore()
	key := "name"

	err := s.Delete(mustNewBulkString(t, key))
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}

func TestDeleteStoredValue(t *testing.T) {
	s := store.NewStore()
	key := "name"
	value := "john"
	mustSetInStore(t, s, key, value)

	bulkStringKey := mustNewBulkString(t, key)
	exists, err := s.Exists(bulkStringKey)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	if !exists {
		t.Errorf("Value should be found, got %v", exists)
	}

	s.Delete(bulkStringKey)

	exists, err = s.Exists(bulkStringKey)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	if exists {
		t.Errorf("Value should not be found, got %v", exists)
	}
}

func TestAddToStoredValue(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		value  string
		addQty int64
		want   int64
	}{
		{"Increments 1 to 2", "counter", "1", 1, 2},
		{"Increments 1 to 3", "counter", "1", 2, 3},
		{"Decrements 1 to 0", "counter", "1", -1, 0},
		{"Decrements 1 to -1", "counter", "1", -2, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()
			mustSetInStore(t, s, tt.key, tt.value)

			bulkStringKey := mustNewBulkString(t, tt.key)
			newValue, err := s.Add(bulkStringKey, tt.addQty)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if newValue.GetValue() != tt.want {
				t.Errorf("incorrect added value, got %d want %d", newValue.GetValue(), tt.want)
			}
		})
	}
}

func TestAddToMissingValue(t *testing.T) {
	key := "counter"

	s := store.NewStore()

	bulkStringKey := mustNewBulkString(t, key)
	newValue, err := s.Add(bulkStringKey, 1)
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	if newValue.GetValue() != int64(1) {
		t.Errorf("incorrect result from adding to missing value, got %d want %d", newValue.GetValue(), 1)
	}
}

func TestPrependToNonArrayValueFails(t *testing.T) {
	s := store.NewStore()
	mustSetInStore(t, s, "notlist", "notlist")

	_, err := s.Prepend(mustNewBulkString(t, "notlist"), mustNewBulkString(t, "notlist"))
	if err == nil {
		t.Errorf("Expected error, got %v", err)
	}
}

func TestPrependElementsToEmptyKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		elements []string
	}{
		{"Prepends nothing to an empty key", "list", []string{}},
		{"Prepends single value to an empty key", "list", []string{"a"}},
		{"Prepends multiple values to an empty key", "list", []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()

			k := mustNewBulkString(t, tt.key)

			elementsToAdd := make([]resp.CashewValue, len(tt.elements))
			for i, e := range tt.elements {
				elementsToAdd[i] = mustNewBulkString(t, e)
			}

			count, err := s.Prepend(k, elementsToAdd...)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if count.GetValue() != int64(len(tt.elements)) {
				t.Errorf("incorrect added count, got %d want %d", count.GetValue(), len(tt.elements))
			}

			stored, err := s.Get(k)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			array, ok := stored.(resp.Array)
			if !ok {
				t.Fatalf("Stored value not array %t", stored)
			}
			values, ok := array.GetValue().([]resp.CashewValue)
			if !ok {
				t.Fatalf("Array values not cashew values %t", stored)
			}
			if len(values) != len(tt.elements) {
				t.Errorf("incorrect stored array length, got %d want %d", len(values), len(tt.elements))
			}

			// Elements are preprended one after another
			for offset := 0; offset < len(tt.elements); offset++ {
				preprended := values[offset].GetValue()
				wanted := tt.elements[len(tt.elements)-offset-1]
				if preprended != wanted {
					t.Errorf("incorrect prepend order, got %s want %s", preprended, wanted)
				}
			}
		})
	}
}

func TestPrependElementsToExistingKey(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		elements []string
	}{
		{"Prepends nothing to an existing key", "list", []string{}},
		{"Prepends single value to an existing key", "list", []string{"a"}},
		{"Prepends multiple values to an existing key", "list", []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()

			k := mustNewBulkString(t, tt.key)
			_, err := s.Prepend(k, mustNewBulkString(t, "0"))
			if err != nil {
				t.Fatalf("Unexpected error populating store %v", err)
			}

			elementsToAdd := make([]resp.CashewValue, len(tt.elements))
			for i, e := range tt.elements {
				elementsToAdd[i] = mustNewBulkString(t, e)
			}

			count, err := s.Prepend(k, elementsToAdd...)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if count.GetValue() != int64(len(tt.elements)+1) {
				t.Errorf("incorrect added count, got %d want %d", count.GetValue(), len(tt.elements)+1)
			}

			stored, err := s.Get(k)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			array, ok := stored.(resp.Array)
			if !ok {
				t.Fatalf("Stored value not array %t", stored)
			}
			values, ok := array.GetValue().([]resp.CashewValue)
			if !ok {
				t.Fatalf("Array values not cashew values %t", stored)
			}
			if len(values) != len(tt.elements)+1 {
				t.Errorf("incorrect stored array length, got %d want %d", len(values), len(tt.elements)+1)
			}

			// Elements are preprended one after another
			for offset := 0; offset < len(tt.elements); offset++ {
				preprended := values[offset].GetValue()
				wanted := tt.elements[len(tt.elements)-offset-1]
				if preprended != wanted {
					t.Errorf("incorrect prepend order, got %s want %s", preprended, wanted)
				}
			}
		})
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
