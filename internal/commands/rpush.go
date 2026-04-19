package commands

import (
	"errors"

	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

const (
	RPUSH = "RPUSH"
)

func HandleRPush(s *store.Store, args []resp.CashewValue) (Result, error) {
	if len(args) < 1 {
		return Result{}, errors.New("missing arguments")
	}
	key, err := ExtractBulkStringArgument(args[0])
	if err != nil {
		return Result{}, err
	}

	count, err := s.Append(key, args[1:]...)
	if err != nil {
		return Result{}, err
	}

	return Result{count.Marshal(), false}, nil
}
