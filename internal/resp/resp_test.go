package resp_test

import (
	"testing"

	"github.com/eduardpeters/cashew/internal/resp"
)

func TestNewSimpleString(t *testing.T) {
	input := "OK"

	got, err := resp.NewSimpleString(input)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if got.Value != input {
		t.Errorf("want %q - got %q", input, got.Value)
	}
}

func TestInvalidSimpleStrings(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"Not\rgood"},
		{"\rNot good"},
		{"Not good\r"},
		{"Not\nvalid"},
		{"\nNot valid"},
		{"Not valid\n"},
	}

	for _, tt := range tests {
		_, err := resp.NewSimpleString(tt.input)
		t.Log(err)
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
	}
}

func TestMarshalSimpleStrings(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{{"OK", "+OK\r\n"},
		{"Hello, World", "+Hello, World\r\n"}}

	for _, tt := range tests {
		value, err := resp.NewSimpleString(tt.input)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		got := value.Marshal()
		if got != tt.expected {
			t.Errorf("want %q - got %q", tt.expected, got)
		}
	}
}

func TestNewSimpleError(t *testing.T) {
	input := "Error message"

	got, err := resp.NewSimpleError(input)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	if got.Value != input {
		t.Errorf("want %q - got %q", input, got.Value)
	}
}

func TestInvalidSimpleErrors(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"Error\rmessage"},
		{"\rError message"},
		{"Not working\r"},
		{"Not\nworking"},
		{"\nsomething wrong"},
		{"something wrong\n"},
	}

	for _, tt := range tests {
		_, err := resp.NewSimpleError(tt.input)
		t.Log(err)
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
	}
}

func TestMarshalSimpleErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{{"ERROR unknown command 'asdf'", "-ERROR unknown command 'asdf'\r\n"},
		{"Error message", "-Error message\r\n"}}

	for _, tt := range tests {
		value, err := resp.NewSimpleError(tt.input)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		got := value.Marshal()
		if got != tt.expected {
			t.Errorf("want %q - got %q", tt.expected, got)
		}
	}
}
