package commands_test

import (
	"testing"
	"testing/synctest"
	"time"

	"github.com/eduardpeters/cashew/internal/commands"
	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

func TestHandleSet(t *testing.T) {
	tests := []struct {
		name     string
		input    []resp.CashewValue
		expected commands.Result
	}{
		{"replies with simple string OK",
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan")},
			commands.Result{mustNewSimpleString(t, "OK").Marshal(), false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()

			r, err := commands.HandleSet(s, tt.input)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if r.Content != tt.expected.Content {
				t.Errorf("incorrect result content want %q, got %q", tt.expected.Content, r.Content)
			}
			if r.CloseConn != tt.expected.CloseConn {
				t.Errorf("incorrect close connection value want %v, got %v", tt.expected.CloseConn, r.CloseConn)
			}

			key := tt.input[0].(resp.BulkString)
			stored, err := s.Get(key)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if stored.GetValue() != tt.input[1].GetValue() {
				t.Errorf("incorrect stored value want %q, got %q", tt.input[1].GetValue(), stored.GetValue())
			}
		})
	}
}

func TestHandleSetInvalidArguments(t *testing.T) {
	tests := []struct {
		name  string
		input []resp.CashewValue
	}{
		{"Fails on empty argument",
			[]resp.CashewValue{},
		},
		{"Fails on single argument",
			[]resp.CashewValue{mustNewBulkString(t, "key")},
		},
		{"Fails on non bulk string arguments - first simple string",
			[]resp.CashewValue{mustNewSimpleString(t, "key"), mustNewBulkString(t, "value")},
		},
		{"Fails on non bulk string arguments - second simple string",
			[]resp.CashewValue{mustNewBulkString(t, "key"), mustNewSimpleString(t, "value")},
		},
		{"Fails on non bulk string arguments - all simple string",
			[]resp.CashewValue{mustNewSimpleString(t, "key"), mustNewSimpleString(t, "value")},
		},
		{"Fails on missing expiry seconds",
			[]resp.CashewValue{mustNewBulkString(t, "key"), mustNewBulkString(t, "value"), mustNewBulkString(t, "EX")},
		},
		{"Fails on negative expiry seconds",
			[]resp.CashewValue{mustNewBulkString(t, "key"), mustNewBulkString(t, "value"), mustNewBulkString(t, "EX"), mustNewBulkString(t, "-1")},
		},
		{"Fails on unknown option",
			[]resp.CashewValue{mustNewBulkString(t, "key"), mustNewBulkString(t, "value"), mustNewBulkString(t, "expire"), mustNewBulkString(t, "1")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()

			_, err := commands.HandleSet(s, tt.input)
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}

func TestHandleSetEx(t *testing.T) {
	tests := []struct {
		name          string
		input         []resp.CashewValue
		expireSeconds int
	}{
		{"sets an expiry for 1 second",
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan"), mustNewBulkString(t, "EX"), mustNewBulkString(t, "1")},
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := store.NewStore()

				_, err := commands.HandleSet(s, tt.input)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				key := tt.input[0].(resp.BulkString)
				stored, err := s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != tt.input[1].GetValue() {
					t.Fatalf("incorrect stored value want %q, got %q", tt.input[1].GetValue(), stored.GetValue())
				}

				time.Sleep(time.Second*time.Duration(tt.expireSeconds) + time.Millisecond)

				stored, err = s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != nil {
					t.Errorf("incorrect value after expiration want %v, got %q", nil, stored.GetValue())
				}
			})
		})
	}
}
