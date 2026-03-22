package commands_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/eduardpeters/cashew/internal/commands"
	"github.com/eduardpeters/cashew/internal/resp"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"Reads simple command",
			"*1\r\n$4\r\nPING\r\n",
			1,
		},
		{"Reads longer commands",
			"*2\r\n$4\r\nECHO\r\n$11\r\nHello World\r\n",
			2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			cmd, err := commands.ParseCommand(bufio.NewReader(r))
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if len(cmd) != tt.expected {
				t.Errorf("incorrect parsed length. want %d, got %d", tt.expected, cmd)
			}
		})
	}
}

func TestParseInvalidCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{

		{"Fails on simple string",
			"+PING\r\n",
		},
		{"Fails on simple error",
			"-ERROR\r\n",
		},
		{"Fails on integers",
			":1234\r\n",
		},
		{"Fails on bulk strings",
			"$4\r\nPING\r\n",
		},
		{"Fails on null",
			"_\r\n",
		},
		{"Fails on bulk string null",
			"$-1\r\n",
		},
		{"Fails on array null",
			"*-1\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			_, err := commands.ParseCommand(bufio.NewReader(r))
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}

func TestHandleCommand(t *testing.T) {
	tests := []struct {
		name     string
		input    []resp.CashewValue
		expected commands.Result
	}{

		{"handles PING with no args",
			[]resp.CashewValue{mustNewBulkString(t, "PING")},
			commands.Result{"+PONG\r\n", false},
		},
		{"handles PING with single arg",
			[]resp.CashewValue{
				mustNewBulkString(t, "PING"),
				mustNewBulkString(t, "hello")},
			commands.Result{"$5\r\nhello\r\n", false},
		},
		{"handles ECHO with single arg",
			[]resp.CashewValue{
				mustNewBulkString(t, "ECHO"),
				mustNewBulkString(t, "hello")},
			commands.Result{"$5\r\nhello\r\n", false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := commands.HandleCommand(tt.input)
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

func TestHandleInvalidCommand(t *testing.T) {
	tests := []struct {
		name  string
		input []resp.CashewValue
	}{

		{"Fails on unknown command",
			[]resp.CashewValue{mustNewBulkString(t, "PONG")},
		},
		{"Fails on non bulk string command",
			[]resp.CashewValue{mustNewSimpleString(t, "PONG")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := commands.HandleCommand(tt.input)
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}
