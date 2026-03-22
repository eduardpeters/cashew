package commands

import (
	"errors"
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
)

const (
	ECHO = "ECHO"
)

func HandleEcho(args []resp.CashewValue) (Result, error) {
	if len(args) == 0 {
		return Result{}, errors.New("missing argument")
	}
	arg := args[0]
	v, ok := arg.(resp.BulkString)
	if !ok {
		return Result{}, fmt.Errorf("argument not bulk string: %v", arg)
	}

	return Result{v.Marshal(), false}, nil
}
