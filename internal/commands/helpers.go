package commands

import (
	"fmt"

	"github.com/eduardpeters/cashew/internal/resp"
)

// Arguments in commands must always be bulk strings
func ExtractArgument(arg resp.CashewValue) (string, error) {
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

// Generates an OK simple string Result
func ResultOK() Result {
	c, err := resp.NewSimpleString("OK")
	if err != nil {
		return ResultError(err)
	}
	return Result{
		Content: c.Marshal(),
	}
}

// Generates an error Result from an existing error
func ResultError(err error) Result {
	e, err := resp.NewSimpleError(fmt.Sprintf("ERR %s", err.Error()))
	if err != nil {
		// Manual RESP encoding for errors that fail validation
		return Result{Content: fmt.Sprintf("%sERR invalid error details\r\n", resp.IDENTIFIER_SIMPLE_ERROR)}
	}
	return Result{
		Content: e.Marshal(),
	}
}
