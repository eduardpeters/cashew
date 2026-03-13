package resp

import (
	"errors"
	"fmt"
)

const (
	IDENTIFER_SIMPLE_ERROR  = "-"
	IDENTIFER_SIMPLE_STRING = "+"
	TERMINATOR              = "\r\n"
)

var INVALID_CHARACTERS_ERROR = errors.New("illegal character present: '\\r' or '\\n'")

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

func hasInvalidCharacters(s string) bool {
	for _, rune := range s {
		if rune == '\r' || rune == '\n' {
			return true
		}
	}

	return false
}
