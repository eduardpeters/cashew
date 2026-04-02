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
	key, err := ExtractBulkStringArgument(args[0])
	if err != nil {
		return Result{}, err
	}
	value, err := ExtractBulkStringArgument(args[1])
	if err != nil {
		return Result{}, err
	}
	if len(args) > 2 {
		if len(args) < 4 {
			return Result{}, errors.New("missing arguments")
		}

		option, err := ExtractArgumentString(args[2])
		if err != nil {
			return Result{}, err
		}
		if option != "EX" {
			return Result{}, fmt.Errorf("Unknown option for SET: %q", option)
		}
		expiry, err := ExtractArgumentInteger(args[3])
		if err != nil {
			return Result{}, err
		}
		if expiry < 0 {
			return Result{}, fmt.Errorf("Expiration must be positive integer: %d", expiry)
		}

		expiration := time.Now().Add(time.Second * 1)
		err = s.SetWithExpiry(key, value, expiration)
		if err != nil {
			return Result{}, err
		}

		return ResultOK(), nil
	}

	err = s.Set(key, value)
	if err != nil {
		return Result{}, err
	}

	return ResultOK(), nil
}
