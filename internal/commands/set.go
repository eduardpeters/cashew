package commands

import (
	"errors"
	"fmt"
	"time"

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
	if len(args) > 2 {
		if len(args) < 4 {
			return Result{}, errors.New("missing arguments")
		}

		opt := args[2]
		_, ok := opt.(resp.BulkString)
		if !ok {
			return Result{}, fmt.Errorf("argument not bulk string: %v", k)
		}
		d := args[3]
		_, ok = d.(resp.BulkString)
		if !ok {
			return Result{}, fmt.Errorf("argument not bulk string: %v", k)
		}

		expiration := time.Now().Add(time.Second * 1)
		err := s.SetWithExpiry(key, value, expiration)
		if err != nil {
			return Result{}, err
		}

		return ResultOK(), nil
	}

	err := s.Set(key, value)
	if err != nil {
		return Result{}, err
	}

	return ResultOK(), nil
}
