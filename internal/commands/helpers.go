package commands

import (
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
)

// Arguments in commands must always be bulk strings
func extractArgument(arg resp.CashewValue) (string, error) {
	bulkString, ok := arg.(resp.BulkString)
	if !ok {
		return "", fmt.Errorf("argument not bulk string: %v", arg)
	}
	v, ok := bulkString.GetValue().(string)
	if !ok {
		return "", fmt.Errorf("bulk string value not string: %v", bulkString.GetValue())
	}
	return v, nil
}
