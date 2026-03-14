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

	if got.GetValue() != input {
		t.Errorf("want %q - got %q", input, got.GetValue())
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

	if got.GetValue() != input {
		t.Errorf("want %q - got %q", input, got.GetValue())
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
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
	}
}

func TestMarshalSimpleErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"ERROR unknown command 'asdf'", "-ERROR unknown command 'asdf'\r\n"},
		{"Error message", "-Error message\r\n"},
	}

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

func TestNewInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"0", 0},
		{"-1", -1},
		{"1", 1},
		{"+42", 42},
		{"+0", 0},
		{"-0", 0},
		{"9223372036854775807", 9223372036854775807},
		{"-9223372036854775808", -9223372036854775808},
	}

	for _, tt := range tests {
		got, err := resp.NewInteger(tt.input)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if got.GetValue() != tt.expected {
			t.Errorf("want %d - got %d", tt.expected, got.GetValue())
		}
	}
}

func TestInvalidIntegers(t *testing.T) {
	tests := []struct {
		input string
	}{
		{"3.14"},
		{"1234s"},
		{"O9128"},
		{"abcdef"},
		{"-123 7"},
		{""},
		{"9223372036854775808"},
		{"-9223372036854775809"},
	}

	for _, tt := range tests {
		_, err := resp.NewInteger(tt.input)
		if err == nil {
			t.Errorf("Expected error, got %v", err)
		}
	}
}

func TestMarshalIntegers(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", ":0\r\n"},
		{"-1", ":-1\r\n"},
		{"1", ":1\r\n"},
		{"9302", ":9302\r\n"},
		{"-1234", ":-1234\r\n"},
	}

	for _, tt := range tests {
		value, err := resp.NewInteger(tt.input)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}
		got := value.Marshal()
		if got != tt.expected {
			t.Errorf("want %q - got %q", tt.expected, got)
		}
	}
}

func TestNewBulkString(t *testing.T) {
	tests := []struct {
		input string
	}{
		{""},
		{"hello"},
		{"hello world"},
		{"this\nis\nallowed"},
		{"this\r\nis also allowed"},
	}

	for _, tt := range tests {
		got, err := resp.NewBulkString(tt.input)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		if got.GetValue() != tt.input {
			t.Errorf("want %q - got %q", tt.input, got.GetValue())
		}
	}
}

func TestMarshalBulkString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "$0\r\n\r\n"},
		{"hello", "$5\r\nhello\r\n"},
		{"hello world", "$11\r\nhello world\r\n"},
		{"this\nis\nallowed", "$15\r\nthis\nis\nallowed\r\n"},
		{"this\r\nis also allowed", "$21\r\nthis\r\nis also allowed\r\n"},
	}

	for _, tt := range tests {
		value, err := resp.NewBulkString(tt.input)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		got := value.Marshal()
		if got != tt.expected {
			t.Errorf("want %q - got %q", tt.expected, got)
		}

	}
}

func TestNewArray(t *testing.T) {
	tests := []struct {
		input []resp.CashewValue
	}{
		{[]resp.CashewValue{}},
	}

	for _, tt := range tests {
		got, err := resp.NewArray(tt.input)

		if err != nil {
			t.Errorf("Unexpected error %v", err)
		}

		gotLength := len(got.Values)
		wantLength := len(tt.input)
		if gotLength != wantLength {
			t.Errorf("wrong array length: want %d - got %d", wantLength, gotLength)
		}

	}
}
