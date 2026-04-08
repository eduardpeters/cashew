package commands_test

import (
	"testing"

	"github.com/eduardpeters/cashew/internal/commands"
	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

func TestHandleIncr(t *testing.T) {
	tests := []struct {
		name          string
		keyValuePairs []KV
		input         []resp.CashewValue
		expected      commands.Result
	}{
		{"replies 1 if store empty",
			[]KV{},
			[]resp.CashewValue{mustNewBulkString(t, "counter")},
			commands.Result{mustNewInteger(t, "1").Marshal(), false},
		},
		{"replies with 2 if key is set with value 1",
			[]KV{{mustNewBulkString(t, "counter"), mustNewBulkString(t, "1")}},
			[]resp.CashewValue{mustNewBulkString(t, "counter")},
			commands.Result{mustNewInteger(t, "2").Marshal(), false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()
			storeValues(t, s, tt.keyValuePairs)

			r, err := commands.HandleIncr(s, tt.input)
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

func TestHandleInvalidIncr(t *testing.T) {
	tests := []struct {
		name          string
		keyValuePairs []KV
		input         []resp.CashewValue
	}{
		{"returns error if no key is provided",
			[]KV{{mustNewBulkString(t, "counter"), mustNewBulkString(t, "un0")}},
			[]resp.CashewValue{},
		},
		{"returns error if stored value is not parseable as int",
			[]KV{{mustNewBulkString(t, "counter"), mustNewBulkString(t, "un0")}},
			[]resp.CashewValue{mustNewBulkString(t, "counter")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()
			storeValues(t, s, tt.keyValuePairs)

			_, err := commands.HandleIncr(s, tt.input)
			if err == nil {
				t.Fatalf("Expected error, got %v", err)
			}
		})
	}
}
