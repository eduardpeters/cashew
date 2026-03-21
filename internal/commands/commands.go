package commands

import (
	"bufio"
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
)

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

func HandleCommand(cmd []resp.CashewValue) (resp.CashewValue, error) {
	return nil, nil
}
