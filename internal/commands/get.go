package commands

import (
	"errors"
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

const (
	GET = "GET"
)

func HandleGet(s *store.Store, args []resp.CashewValue) (Result, error) {
	if len(args) < 1 {
		return Result{}, errors.New("missing argument")
	}
	k := args[0]
	key, ok := k.(resp.BulkString)
	if !ok {
		return Result{}, fmt.Errorf("argument not bulk string: %v", k)
	}

	value, err := s.Get(key)
	if err != nil {
		return Result{}, err
	}

	return Result{value.Marshal(), false}, nil
}
