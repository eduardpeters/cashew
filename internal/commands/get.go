package commands

import (
	"errors"

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
	key, err := ExtractBulkStringArgument(args[0])
	if err != nil {
		return Result{}, err
	}

	value, err := s.Get(key)
	if err != nil {
		return Result{}, err
	}

	return Result{value.Marshal(), false}, nil
}
