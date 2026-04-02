package commands

import (
	"errors"
	"strconv"

	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

const (
	EXISTS = "EXISTS"
)

func HandleExists(s *store.Store, args []resp.CashewValue) (Result, error) {
	if len(args) < 1 {
		return Result{}, errors.New("missing argument")
	}
	count := 0
	for _, arg := range args {
		key, err := ExtractBulkStringArgument(arg)
		if err != nil {
			return Result{}, err
		}

		found, err := s.Exists(key)
		if err != nil {
			return Result{}, err
		}

		if found {
			count++
		}
	}

	value, err := resp.NewInteger(strconv.Itoa(count))
	if err != nil {
		return Result{}, err
	}

	return Result{value.Marshal(), false}, nil
}
