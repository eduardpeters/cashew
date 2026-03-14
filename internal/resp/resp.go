package resp

import (
	"errors"
	"fmt"
	"strconv"
)

// VALUE TYPE IDENTIFIERS
const (
	IDENTIFER_SIMPLE_ERROR  = "-"
	IDENTIFER_INTEGER       = ":"
	IDENTIFER_SIMPLE_STRING = "+"
	TERMINATOR              = "\r\n"
)

// ERRORS
var (
	INVALID_CHARACTERS_ERROR = errors.New("illegal character present: '\\r' or '\\n'")
	INVALID_INTEGER_ERROR    = errors.New("input is not an integer")
)

type CashewValue interface {
	Marshal() string
}

type SimpleString struct {
	Value string
}

func NewSimpleString(s string) (*SimpleString, error) {
	for _, rune := range s {
		if rune == '\r' || rune == '\n' {
			return nil, INVALID_CHARACTERS_ERROR
		}
	}
	return &SimpleString{s}, nil
}

func (s SimpleString) Marshal() string {
	return fmt.Sprintf("%s%s%s", IDENTIFER_SIMPLE_STRING, s.Value, TERMINATOR)
}

type SimpleError struct {
	Value string
}

func NewSimpleError(s string) (*SimpleError, error) {
	if hasInvalidCharacters(s) {
		return nil, INVALID_CHARACTERS_ERROR
	}

	return &SimpleError{s}, nil
}

func (s SimpleError) Marshal() string {
	return fmt.Sprintf("%s%s%s", IDENTIFER_SIMPLE_ERROR, s.Value, TERMINATOR)
}

type Integer struct {
	Value int64
}

func NewInteger(s string) (*Integer, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, INVALID_INTEGER_ERROR
	}
	return &Integer{n}, nil
}

func (i Integer) Marshal() string {
	return fmt.Sprintf("%s%d%s", IDENTIFER_INTEGER, i.Value, TERMINATOR)
}

func hasInvalidCharacters(s string) bool {
	for _, rune := range s {
		if rune == '\r' || rune == '\n' {
			return true
		}
	}

	return false
}
