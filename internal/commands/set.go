package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

const (
	SET       = "SET"
	OPTION_EX = "EX"
	OPTION_PX = "PX"
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
	if len(args) == 2 {
		return handlePlainSet(s, key, value)
	}

	return handleSetWithOptions(s, key, value, args[2:])
}

func handlePlainSet(s *store.Store, key, value resp.BulkString) (Result, error) {
	err := s.Set(key, value)
	if err != nil {
		return Result{}, err
	}

	return ResultOK(), nil
}

func handleSetWithOptions(s *store.Store, key, value resp.BulkString, options []resp.CashewValue) (Result, error) {
	if len(options) < 2 {
		return Result{}, errors.New("missing arguments")
	}
	option, err := ExtractArgumentString(options[0])
	if err != nil {
		return Result{}, err
	}
	if option != OPTION_EX && option != OPTION_PX {
		return Result{}, fmt.Errorf("Unknown option for SET: %q", option)
	}
	expiry, err := ExtractArgumentInteger(options[1])
	if err != nil {
		return Result{}, err
	}
	if expiry < 0 {
		return Result{}, fmt.Errorf("Expiration must be positive integer: %d", expiry)
	}

	var expiration time.Time
	switch option {
	case OPTION_EX:
		expiration = time.Now().Add(time.Second * time.Duration(expiry))
	case OPTION_PX:
		expiration = time.Now().Add(time.Millisecond * time.Duration(expiry))
	default:
		return Result{}, fmt.Errorf("Unknown option for SET: %q", option)
	}

	err = s.SetWithExpiry(key, value, expiration)
	if err != nil {
		return Result{}, err
	}

	return ResultOK(), nil
}
