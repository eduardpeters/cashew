package commands

import (
	"testing"

	"github.com/eduardpeters/cashew/internal/resp"
)

func TestExtractArgument(t *testing.T) {
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
			arg, err := extractArgument(tt.input)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if arg != tt.expected {
				t.Errorf("incorrect argument. want %q, got %q", tt.expected, arg)
			}
		})
	}
}

func TestExtractInvalidArgument(t *testing.T) {
	tests := []struct {
		name  string
		input resp.CashewValue
	}{
		{"Returns string from BulkString",
			mustNewSimpleString(t, "PING"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := extractArgument(tt.input)
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

func mustNewBulkString(t testing.TB, s string) resp.CashewValue {
	t.Helper()
	v, err := resp.NewBulkString(s)
	if err != nil {
		t.Fatalf("NewBulkString(%q): %v", s, err)
	}
	return v
}
