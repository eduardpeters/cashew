package commands_test

import (
	"testing"

	"github.com/eduardpeters/cashew/internal/commands"
	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

func TestExtractBulkStringArgument(t *testing.T) {
	tests := []struct {
		name     string
		input    resp.CashewValue
		expected string
	}{
		{"Returns string from BulkString",
			mustNewBulkString(t, "PING"),
			"PING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg, err := commands.ExtractBulkStringArgument(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if arg.GetValue() != tt.expected {
				t.Errorf("incorrect argument. want %q, got %q", tt.expected, arg)
			}
		})
	}
}

func TestExtractArgumentAsString(t *testing.T) {
	tests := []struct {
		name     string
		input    resp.CashewValue
		expected string
	}{
		{"Returns string from BulkString",
			mustNewBulkString(t, "PING"),
			"PING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg, err := commands.ExtractArgumentString(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if arg != tt.expected {
				t.Errorf("incorrect argument. want %q, got %q", tt.expected, arg)
			}
		})
	}
}

func TestExtractInvalidStringArgument(t *testing.T) {
	tests := []struct {
		name  string
		input resp.CashewValue
	}{
		{"Cannot extract from SimpleString",
			mustNewSimpleString(t, "PING"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := commands.ExtractArgumentString(tt.input)
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}

func TestExtractArgumentAsInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    resp.CashewValue
		expected int64
	}{
		{"Extracts 0",
			mustNewBulkString(t, "0"),
			0,
		},
		{"Extracts 1",
			mustNewBulkString(t, "1"),
			1,
		},
		{"Extracts -1",
			mustNewBulkString(t, "-1"),
			-1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arg, err := commands.ExtractArgumentInteger(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if arg != tt.expected {
				t.Errorf("incorrect argument. want %d, got %d", tt.expected, arg)
			}
		})
	}
}

func TestExtractInvalidIntegerArgument(t *testing.T) {
	tests := []struct {
		name  string
		input resp.CashewValue
	}{
		{"Cannot extract string as integer",
			mustNewBulkString(t, "uno"),
		},
		{"Cannot extract integer from simple string",
			mustNewSimpleString(t, "123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := commands.ExtractArgumentInteger(tt.input)
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}

func mustNewSimpleString(t testing.TB, s string) resp.CashewValue {
	t.Helper()
	v, err := resp.NewSimpleString(s)
	if err != nil {
		t.Fatalf("NewSimpleString(%q): %v", s, err)
	}
	return v
}

func mustNewInteger(t testing.TB, s string) resp.CashewValue {
	t.Helper()
	v, err := resp.NewInteger(s)
	if err != nil {
		t.Fatalf("NewInteger(%q): %v", s, err)
	}
	return v
}

func mustNewBulkString(t testing.TB, s string) resp.CashewValue {
	t.Helper()
	v, err := resp.NewBulkString(s)
	if err != nil {
		t.Fatalf("NewBulkString(%q): %v", s, err)
	}
	return v
}

func mustNewNull(t testing.TB) resp.CashewValue {
	t.Helper()
	v, err := resp.NewNull()
	if err != nil {
		t.Fatalf("NewNull(): %v", err)
	}
	return v
}

type KV struct {
	key   resp.CashewValue
	value resp.CashewValue
}

func storeValues(t testing.TB, s *store.Store, keyValuePairs []KV) {
	t.Helper()
	for _, kv := range keyValuePairs {
		key := kv.key.(resp.BulkString)
		value := kv.value.(resp.BulkString)

		err := s.Set(key, value)
		if err != nil {
			t.Fatalf("unexpected error %v", err)
		}
	}
}
