package commands

import (
	"errors"
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

const (
	SET = "SET"
)

func HandleSet(s *store.Store, args []resp.CashewValue) (Result, error) {
	if len(args) < 2 {
		return Result{}, errors.New("missing arguments")
	}
	k := args[0]
	key, ok := k.(resp.BulkString)
	if !ok {
		return Result{}, fmt.Errorf("argument not bulk string: %v", k)
	}
	v := args[1]
	value, ok := v.(resp.BulkString)
	if !ok {
		return Result{}, fmt.Errorf("argument not bulk string: %v", k)
	}

	err := s.Set(key, value)
	if err != nil {
		return Result{}, err
	}

	return ResultOK(), nil
}
