package commands

import (
	"bufio"
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
)

type Result struct {
	Content   string
	CloseConn bool
}

func ParseCommand(b *bufio.Reader) ([]resp.CashewValue, error) {
	parsed, err := resp.Unmarshal(b)
	if err != nil {
		return nil, err
	}

	switch cmd := parsed.GetValue().(type) {
	case []resp.CashewValue:
		return cmd, nil
	default:
		return nil, fmt.Errorf("command must be array, got %v", cmd)
	}
}

func HandleCommand(cmd []resp.CashewValue) (Result, error) {
	if len(cmd) == 0 {
		return Result{}, fmt.Errorf("empty command")
	}

	verb, err := ExtractArgument(cmd[0])
	if err != nil {
		return Result{}, fmt.Errorf("argument parsing error: %w", err)
	}

	switch verb {
	case PING:
		return HandlePing(cmd[1:])
	case ECHO:
		return HandleEcho(cmd[1:])
	// Commands outside the scope to support clients
	case "CLIENT":
		return ResultOK(), nil
	default:
		return Result{}, fmt.Errorf("unknown command %q", verb)
	}
}
