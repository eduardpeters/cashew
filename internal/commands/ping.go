package commands

import (
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
		v, err := ExtractBulkStringArgument(args[0])
		if err != nil {
			return Result{}, err
		}
		content = v
	}

	if err != nil {
		return Result{}, err
	}
	return Result{content.Marshal(), false}, nil
}
