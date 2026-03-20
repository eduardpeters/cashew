package commands_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/eduardpeters/cashew/internal/commands"
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
