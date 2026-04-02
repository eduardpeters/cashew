package commands

import (
	"errors"

	"github.com/eduardpeters/cashew/internal/resp"
)

const (
	ECHO = "ECHO"
)

func HandleEcho(args []resp.CashewValue) (Result, error) {
	if len(args) == 0 {
		return Result{}, errors.New("missing argument")
	}
	arg, err := ExtractBulkStringArgument(args[0])
	if err != nil {
		return Result{}, err
	}

	return Result{arg.Marshal(), false}, nil
}
