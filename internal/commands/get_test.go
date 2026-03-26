package commands_test

import (
	"testing"

	"github.com/eduardpeters/cashew/internal/commands"
	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

type keyValuePair struct {
	key   string
	value string
}

func TestHandleGet(t *testing.T) {
	tests := []struct {
		name     string
		pairs    []keyValuePair
		input    []resp.CashewValue
		expected commands.Result
	}{
		{"replies null if store empty",
			[]keyValuePair{},
			[]resp.CashewValue{mustNewBulkString(t, "name")},
			commands.Result{mustNewNull(t).Marshal(), false},
		},
		{"replies with null if key is not found",
			[]keyValuePair{{"name", "juan"}},
			[]resp.CashewValue{mustNewBulkString(t, "not:there")},
			commands.Result{mustNewNull(t).Marshal(), false},
		},
		{"replies with stored value if key is found",
			[]keyValuePair{{"name", "juan"}},
			[]resp.CashewValue{mustNewBulkString(t, "name")},
			commands.Result{mustNewBulkString(t, "juan").Marshal(), false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()
			for _, pair := range tt.pairs {
				k := mustNewBulkString(t, pair.key)
				v := mustNewBulkString(t, pair.value)
				key, ok := k.(resp.BulkString)
				if !ok {
					t.Fatalf("Key not BulkString %v", k)
				}
				value, ok := v.(resp.BulkString)
				if !ok {
					t.Fatalf("Value not BulkString %v", v)
				}
				s.Set(key, value)
			}

			r, err := commands.HandleGet(s, tt.input)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if r.Content != tt.expected.Content {
				t.Errorf("incorrect result content want %q, got %q", tt.expected.Content, r.Content)
			}
			if r.CloseConn != tt.expected.CloseConn {
				t.Errorf("incorrect close connection value want %v, got %v", tt.expected.CloseConn, r.CloseConn)
			}
		})
	}
}

func TestHandleGetInvalidArguments(t *testing.T) {
	tests := []struct {
		name  string
		input []resp.CashewValue
	}{
		{"Fails on empty argument",
			[]resp.CashewValue{},
		},
		{"Fails on non bulk string argument",
			[]resp.CashewValue{mustNewSimpleString(t, "key")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()

			_, err := commands.HandleGet(s, tt.input)
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}
