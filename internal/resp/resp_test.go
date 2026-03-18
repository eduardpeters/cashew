package resp_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/eduardpeters/cashew/internal/resp"
)

func TestNewNull(t *testing.T) {
	got, err := resp.NewNull()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
	if got.GetValue() != nil {
		t.Errorf("want nil, got %v", got.GetValue())
	}
}

func TestMarshalNull(t *testing.T) {
	value, err := resp.NewNull()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	got := value.Marshal()
	want := "_\r\n"
	if got != want {
		t.Errorf("want %q - got %q", want, got)
	}
}

func TestNewSimpleString(t *testing.T) {
	input := "OK"

	got, err := resp.NewSimpleString(input)

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
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
				t.Fatalf("Unexpected error %v", err)
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
		t.Fatalf("Unexpected error %v", err)
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
				t.Fatalf("Unexpected error %v", err)
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
				t.Fatalf("Unexpected error %v", err)
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
				t.Fatalf("Unexpected error %v", err)
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
				t.Fatalf("Unexpected error %v", err)
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
				t.Fatalf("Unexpected error %v", err)
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
		{"mixed data types", []resp.CashewValue{
			mustNewSimpleString(t, "numbers"),
			mustNewInteger(t, "1"),
			mustNewInteger(t, "2"),
			mustNewInteger(t, "3"),
			mustNewBulkString(t, "the end"),
		}, "*5\r\n+numbers\r\n:1\r\n:2\r\n:3\r\n$7\r\nthe end\r\n"},
		{"simple nested array", []resp.CashewValue{
			mustNewArray(
				t, []resp.CashewValue{mustNewSimpleString(t, "nested")},
			),
		}, "*1\r\n*1\r\n+nested\r\n"},
		{"paired nested arrays", []resp.CashewValue{
			mustNewArray(
				t, []resp.CashewValue{
					mustNewInteger(t, "1"),
					mustNewInteger(t, "2"),
					mustNewInteger(t, "3"),
				},
			),
			mustNewArray(
				t, []resp.CashewValue{
					mustNewBulkString(t, "Hello"),
					mustNewBulkString(t, "World"),
				},
			),
		}, "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n$5\r\nHello\r\n$5\r\nWorld\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := resp.NewArray(tt.input)

			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			got := value.Marshal()
			if got != tt.expected {
				t.Errorf("want %q - got %q", tt.expected, got)
			}
		})
	}
}

func TestUnmarshalInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"missing identifier", "hello\r\n"},
		{"invalid identifier", "@hello\r\n"},
		{"unterminated", "+hello"},
		{"only LF", "+hello world\n"},
		{"null only LF", "_\n"},
		{"null only CR", "_\r"},
		{"only CR", "+hello world\r"},
		{"invalid CR in simple string", "+hello\rworld\r\n"},
		{"invalid LF in simple string", "+hello\nworld\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			_, err := resp.Unmarshal(bufio.NewReader(r))
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}

func TestUnmarshalNull(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Null array", "*-1\r\n"},
		{"Null bulk string", "$-1\r\n"},
		{"Null value", "_\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			value, err := resp.Unmarshal(bufio.NewReader(r))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			got := value.GetValue()
			if got != nil {
				t.Errorf("want nil, got %v", got)
			}
		})
	}
}

func TestUnmarshalInvalidNull(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Not -1 array", "*-2\r\n"},
		{"Not -1 bulk string", "$-2\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			_, err := resp.Unmarshal(bufio.NewReader(r))
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}

func TestUnmarshalSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no whitespace", "+hello\r\n", "hello"},
		{"whitespace", "+hello world\r\n", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			v, err := resp.Unmarshal(bufio.NewReader(r))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			value, ok := v.(resp.SimpleString)
			if !ok {
				t.Fatalf("value not Simple String")
			}

			got := value.GetValue()
			if got != tt.expected {
				t.Errorf("want %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestUnmarshalSimpleErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no whitespace", "-ERROR\r\n", "ERROR"},
		{"whitespace", "-ERR unknown command 'asdf'\r\n", "ERR unknown command 'asdf'"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			v, err := resp.Unmarshal(bufio.NewReader(r))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			value, ok := v.(resp.SimpleError)
			if !ok {
				t.Fatalf("value not Simple Error")
			}

			got := value.GetValue()
			if got != tt.expected {
				t.Errorf("want %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestUnmarshalIntegers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"zero", ":0\r\n", 0},
		{"positive one", ":+1\r\n", 1},
		{"negative one", ":-1\r\n", -1},
		{"positive multiple digits", ":5432\r\n", 5432},
		{"negative multiple digits", ":-42\r\n", -42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			v, err := resp.Unmarshal(bufio.NewReader(r))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			value, ok := v.(resp.Integer)
			if !ok {
				t.Fatalf("value not Integer")
			}

			got := value.GetValue()
			if got != tt.expected {
				t.Errorf("want %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestUnmarshalBulkString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "$0\r\n\r\n", ""},
		{"one character", "$1\r\na\r\n", "a"},
		{"no whitespace", "$5\r\nhello\r\n", "hello"},
		{"with whitespace", "$11\r\nhello world\r\n", "hello world"},
		{"with LF", "$11\r\nhello\nworld\r\n", "hello\nworld"},
		{"with CR", "$11\r\nhello\rworld\r\n", "hello\rworld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			v, err := resp.Unmarshal(bufio.NewReader(r))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			value, ok := v.(resp.BulkString)
			if !ok {
				t.Fatalf("value not BulkString")
			}

			got := value.GetValue()
			if got != tt.expected {
				t.Errorf("want %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestUnmarshalInvalidBulkString(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty without terminator sequence", "$0\r\n"},
		{"without terminator sequence", "$11\r\nhello world"},
		{"longer data than length declared", "$1\r\nabc\r\n"},
		{"shorter data than length declared", "$5\r\nhey\r\n"},
		{"negative length", "$-5\r\nhello\r\n"},
		{"length as string", "$two\r\nhi\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			_, err := resp.Unmarshal(bufio.NewReader(r))
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}

func TestUnmarshalArray(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []resp.CashewValue
	}{
		{"empty", "*0\r\n", []resp.CashewValue{}},
		{"single simple string", "*1\r\n+hello\r\n", []resp.CashewValue{mustNewSimpleString(t, "hello")}},
		{"multiple simple strings", "*2\r\n+hello\r\n+world\r\n", []resp.CashewValue{
			mustNewSimpleString(t, "hello"),
			mustNewSimpleString(t, "world"),
		}},
		{"single integer",
			"*1\r\n:6379\r\n",
			[]resp.CashewValue{
				mustNewInteger(t, "6379"),
			},
		},
		{"multiple integers",
			"*3\r\n:1234\r\n:4321\r\n:56789\r\n",
			[]resp.CashewValue{
				mustNewInteger(t, "1234"),
				mustNewInteger(t, "4321"),
				mustNewInteger(t, "56789"),
			},
		},
		{"single bulk string",
			"*1\r\n$5\r\nhello\r\n",
			[]resp.CashewValue{
				mustNewBulkString(t, "hello"),
			},
		},
		{"multiple bulk strings",
			"*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$5\r\nvalue\r\n",
			[]resp.CashewValue{
				mustNewBulkString(t, "set"),
				mustNewBulkString(t, "key"),
				mustNewBulkString(t, "value"),
			},
		},
		{"longer data than length declared only parses declared length",
			"*1\r\n+abc\r\n+xyz\r\n",
			[]resp.CashewValue{
				mustNewSimpleString(t, "abc"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			v, err := resp.Unmarshal(bufio.NewReader(r))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			value, ok := v.(resp.Array)
			if !ok {
				t.Fatal("value not Array")
			}

			arrayValue := value.GetValue()
			elements, ok := arrayValue.([]resp.CashewValue)
			if !ok {
				t.Fatal("internal values not CashewValue")
			}
			if len(elements) != len(tt.expected) {
				t.Fatalf("incorrect length for array: want %d - got %d", len(tt.expected), len(elements))
			}
			for i, e := range elements {
				want := tt.expected[i].GetValue()
				got := e.GetValue()
				if got != want {
					t.Errorf("want %q, got %q", want, got)
				}
			}
		})
	}
}

func TestUnmarshalNestedArrays(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []resp.CashewValue
	}{
		{"simple nested array",
			"*1\r\n*1\r\n+nested\r\n",
			[]resp.CashewValue{
				mustNewArray(
					t, []resp.CashewValue{mustNewSimpleString(t, "nested")},
				),
			},
		},
		{"paired nested arrays",
			"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n$5\r\nHello\r\n$5\r\nWorld\r\n",
			[]resp.CashewValue{
				mustNewArray(
					t, []resp.CashewValue{
						mustNewInteger(t, "1"),
						mustNewInteger(t, "2"),
						mustNewInteger(t, "3"),
					},
				),
				mustNewArray(
					t, []resp.CashewValue{
						mustNewBulkString(t, "Hello"),
						mustNewBulkString(t, "World"),
					},
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			v, err := resp.Unmarshal(bufio.NewReader(r))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			value, ok := v.(resp.Array)
			if !ok {
				t.Fatal("value not Array")
			}

			arrayValue := value.GetValue()
			topLevelElements, ok := arrayValue.([]resp.CashewValue)
			if !ok {
				t.Fatal("internal values not CashewValue")
			}
			if len(topLevelElements) != len(tt.expected) {
				t.Fatalf("incorrect length for array: want %d - got %d", len(tt.expected), len(topLevelElements))
			}
			for i, tle := range topLevelElements {
				got := tle.GetValue()
				nestedElements, ok := got.([]resp.CashewValue)
				if !ok {
					t.Fatal("nested values not CashewValue")
				}
				want := tt.expected[i].GetValue()
				nestedExpected, ok := want.([]resp.CashewValue)
				if !ok {
					t.Fatal("nested expected values not CashewValue")
				}
				if len(nestedElements) != len(nestedExpected) {
					t.Fatalf("incorrect length for nested array: want %d - got %d", len(nestedExpected), len(nestedElements))
				}
				for i, ne := range nestedElements {
					expected := nestedExpected[i].GetValue()
					nested := ne.GetValue()
					if nested != expected {
						t.Errorf("want %q, got %q", expected, nested)
					}
				}
			}
		})
	}
}

func TestUnmarshalInvalidArray(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty without terminator sequence", "*0"},
		{"data without terminator sequence", "*1\r\n+hello world"},
		{"shorter data than length declared", "*2\r\n+hey\r\n"},
		{"missing length", "*\r\n+hello\r\n"},
		{"negative length", "*-5\r\n+hello\r\n"},
		{"length as string", "*two\r\n+hi\r\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			_, err := resp.Unmarshal(bufio.NewReader(r))
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

func mustNewArray(t testing.TB, values []resp.CashewValue) resp.CashewValue {
	t.Helper()
	v, err := resp.NewArray(values)
	if err != nil {
		t.Fatalf("NewArray(%v): %v", values, err)
	}
	return v
}
