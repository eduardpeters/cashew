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
		name  string
		input string
	}{
		{"with CR", "Not\rgood"},
		{"leading CR", "\rNot good"},
		{"trailing CR", "Not good\r"},
		{"with LF", "Not\nvalid"},
		{"leading LF", "\nNot valid"},
		{"trailing LF", "Not valid\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resp.NewSimpleString(tt.input)
			if err == nil {
				t.Errorf("Expected error, got %v", err)
			}
		})
	}
}

func TestMarshalSimpleStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no whitespace", "OK", "+OK\r\n"},
		{"with whitespace", "Hello, World", "+Hello, World\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resp.NewSimpleString(tt.input)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
			got := value.Marshal()
			if got != tt.expected {
				t.Errorf("want %q - got %q", tt.expected, got)
			}
		})
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
		name  string
		input string
	}{
		{"with CR", "Error\rmessage"},
		{"leading CR", "\rError message"},
		{"trailing CR", "Not working\r"},
		{"with LF", "Not\nworking"},
		{"leading LF", "\nsomething wrong"},
		{"trailing LF", "something wrong\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resp.NewSimpleError(tt.input)
			if err == nil {
				t.Errorf("Expected error, got %v", err)
			}
		})
	}
}

func TestMarshalSimpleErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"longer message", "ERROR unknown command 'asdf'", "-ERROR unknown command 'asdf'\r\n"},
		{"shorter message", "Error message", "-Error message\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resp.NewSimpleError(tt.input)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
			got := value.Marshal()
			if got != tt.expected {
				t.Errorf("want %q - got %q", tt.expected, got)
			}
		})
	}
}

func TestNewInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"zero", "0", 0},
		{"negative one", "-1", -1},
		{"one", "1", 1},
		{"with leading plus", "+42", 42},
		{"plus zero", "+0", 0},
		{"negative zero", "-0", 0},
		{"largest int64", "9223372036854775807", 9223372036854775807},
		{"lowest int64", "-9223372036854775808", -9223372036854775808},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resp.NewInteger(tt.input)

			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			if got.GetValue() != tt.expected {
				t.Errorf("want %d - got %d", tt.expected, got.GetValue())
			}
		})
	}
}

func TestInvalidIntegers(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"floating point", "3.14"},
		{"trailing letter", "1234s"},
		{"leading letter", "O9128"},
		{"only letters", "abcdef"},
		{"whitespace", "-123 7"},
		{"empty", ""},
		{"largest int64 + 1", "9223372036854775808"},
		{"lowest int64 -1", "-9223372036854775809"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := resp.NewInteger(tt.input)
			if err == nil {
				t.Errorf("Expected error, got %v", err)
			}
		})
	}
}

func TestMarshalIntegers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"zero", "0", ":0\r\n"},
		{"negative one", "-1", ":-1\r\n"},
		{"one", "1", ":1\r\n"},
		{"positive number", "9302", ":9302\r\n"},
		{"negative number", "-1234", ":-1234\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resp.NewInteger(tt.input)
			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}
			got := value.Marshal()
			if got != tt.expected {
				t.Errorf("want %q - got %q", tt.expected, got)
			}
		})
	}
}

func TestNewBulkString(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"single string", "hello"},
		{"whitespace", "hello world"},
		{"with newlines", "this\nis\nallowed"},
		{"with CRLF", "this\r\nis also allowed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resp.NewBulkString(tt.input)

			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			if got.GetValue() != tt.input {
				t.Errorf("want %q - got %q", tt.input, got.GetValue())
			}
		})
	}
}

func TestMarshalBulkString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", "$0\r\n\r\n"},
		{"single string", "hello", "$5\r\nhello\r\n"},
		{"whitespace", "hello world", "$11\r\nhello world\r\n"},
		{"with newlines", "this\nis\nallowed", "$15\r\nthis\nis\nallowed\r\n"},
		{"with CRLF", "this\r\nis also allowed", "$21\r\nthis\r\nis also allowed\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resp.NewBulkString(tt.input)

			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			got := value.Marshal()
			if got != tt.expected {
				t.Errorf("want %q - got %q", tt.expected, got)
			}
		})
	}
}

func TestNewArray(t *testing.T) {
	tests := []struct {
		name  string
		input []resp.CashewValue
	}{
		{"empty", []resp.CashewValue{}},
		{"single simple string", []resp.CashewValue{
			mustNewSimpleString(t, "hello"),
		}},
		{"multiple simple string", []resp.CashewValue{
			mustNewSimpleString(t, "hello"),
			mustNewSimpleString(t, "world"),
		}},
		{"single integer", []resp.CashewValue{
			mustNewInteger(t, "6379"),
		}},
		{"multiple integers", []resp.CashewValue{
			mustNewInteger(t, "1234"),
			mustNewInteger(t, "4321"),
			mustNewInteger(t, "56789"),
		}},
		{"single bulk string", []resp.CashewValue{
			mustNewBulkString(t, "hello"),
		}},
		{"multiple bulk strings", []resp.CashewValue{
			mustNewBulkString(t, "set"),
			mustNewBulkString(t, "key"),
			mustNewBulkString(t, "value"),
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resp.NewArray(tt.input)

			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			gotValuesRaw := got.GetValue()
			gotValues, ok := gotValuesRaw.([]resp.CashewValue)
			if !ok {
				t.Fatalf("array value not []CashewValue")
			}

			gotLength := len(gotValues)
			wantLength := len(tt.input)
			if gotLength != wantLength {
				t.Errorf("wrong array length: want %d - got %d", wantLength, gotLength)
			}

			for i := range gotValues {
				if gotValues[i] != tt.input[i] {
					t.Errorf("wrong array element: want %v - got %v", tt.input[i], gotValues[i])
				}
			}
		})
	}
}

func TestMarshalArray(t *testing.T) {
	tests := []struct {
		name     string
		input    []resp.CashewValue
		expected string
	}{
		{"empty", []resp.CashewValue{}, "*0\r\n"},
		{"single simple string", []resp.CashewValue{
			mustNewSimpleString(t, "hello"),
		}, "*1\r\n+hello\r\n"},
		{"multiple simple string", []resp.CashewValue{
			mustNewSimpleString(t, "hello"),
			mustNewSimpleString(t, "world"),
		}, "*2\r\n+hello\r\n+world\r\n"},
		{"single integer", []resp.CashewValue{
			mustNewInteger(t, "6379"),
		}, "*1\r\n:6379\r\n"},
		{"multiple integers", []resp.CashewValue{
			mustNewInteger(t, "1234"),
			mustNewInteger(t, "4321"),
			mustNewInteger(t, "56789"),
		}, "*3\r\n:1234\r\n:4321\r\n:56789\r\n"},
		{"single bulk string", []resp.CashewValue{
			mustNewBulkString(t, "hello"),
		}, "*1\r\n$5\r\nhello\r\n"},
		{"multiple bulk strings", []resp.CashewValue{
			mustNewBulkString(t, "set"),
			mustNewBulkString(t, "key"),
			mustNewBulkString(t, "value"),
		}, "*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resp.NewArray(tt.input)

			if err != nil {
				t.Errorf("Unexpected error %v", err)
			}

			got := value.Marshal()
			if got != tt.expected {
				t.Errorf("want %q - got %q", tt.expected, got)
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
