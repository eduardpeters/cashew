package commands

import (
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
)

const (
	PING = "PING"
	PONG = "PONG"
)

func HandlePing(args []resp.CashewValue) (Result, error) {
	var content resp.CashewValue
	var err error

	if len(args) == 0 {
		content, err = resp.NewSimpleString(PONG)
	} else {
		arg := args[0]
		if v, ok := arg.(resp.BulkString); !ok {
			err = fmt.Errorf("argument not bulk string: %v", arg)
		} else {
			content = v
		}
	}

	if err != nil {
		return Result{}, err
	}
	return Result{content.Marshal(), false}, nil
}
