package commands_test

import (
	"testing"

	"github.com/eduardpeters/cashew/internal/commands"
	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

func TestHandleExists(t *testing.T) {
	tests := []struct {
		name          string
		keyValuePairs []KV
		input         []resp.CashewValue
		expected      commands.Result
	}{
		{"Replies with zero count if nothing stored",
			[]KV{},
			[]resp.CashewValue{mustNewBulkString(t, "key")},
			commands.Result{mustNewInteger(t, "0").Marshal(), false},
		},
		{"Replies with zero count if key does not match",
			[]KV{
				{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan")},
			},
			[]resp.CashewValue{mustNewBulkString(t, "nam3")},
			commands.Result{mustNewInteger(t, "0").Marshal(), false},
		},
		{"Replies with count if single key is stored",
			[]KV{
				{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan")},
			},
			[]resp.CashewValue{mustNewBulkString(t, "name")},
			commands.Result{mustNewInteger(t, "1").Marshal(), false},
		},
		{"Replies with count if multiple keys are stored",
			[]KV{
				{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan")},
				{mustNewBulkString(t, "nombre"), mustNewBulkString(t, "john")},
			},
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "nombre")},
			commands.Result{mustNewInteger(t, "2").Marshal(), false},
		},
		{"Replies with stored count if some keys are not stored",
			[]KV{
				{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan")},
				{mustNewBulkString(t, "nombre"), mustNewBulkString(t, "john")},
			},
			[]resp.CashewValue{mustNewBulkString(t, "nombre"), mustNewBulkString(t, "key"), mustNewBulkString(t, "name")},
			commands.Result{mustNewInteger(t, "2").Marshal(), false},
		},
		{"Replies with double count if same found key is requested twice",
			[]KV{
				{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan")},
			},
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "name")},
			commands.Result{mustNewInteger(t, "2").Marshal(), false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()
			storeValues(t, s, tt.keyValuePairs)

			r, err := commands.HandleExists(s, tt.input)
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

func TestHandleExistsInvalidArguments(t *testing.T) {
	tests := []struct {
		name  string
		input []resp.CashewValue
	}{
		{"Fails on empty arguments",
			[]resp.CashewValue{},
		},
		{"Fails on non bulk string argument",
			[]resp.CashewValue{mustNewSimpleString(t, "key")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()

			_, err := commands.HandleExists(s, tt.input)
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}
