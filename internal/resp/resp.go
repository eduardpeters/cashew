package resp

import (
	"errors"
	"fmt"
)

const (
	IDENTIFER_SIMPLE_STRING = "+"
	TERMINATOR              = "\r\n"
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
			return nil, errors.New("illegal character present: '\\r' or '\\n'")
		}
	}
	return &SimpleString{s}, nil
}

func (s SimpleString) Marshal() string {
	return fmt.Sprintf("%s%s%s", IDENTIFER_SIMPLE_STRING, s.Value, TERMINATOR)
}
