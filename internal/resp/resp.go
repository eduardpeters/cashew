package resp

import (
	"errors"
	"strconv"
	"strings"
)

// VALUE TYPE IDENTIFIERS
const (
	IDENTIFER_SIMPLE_ERROR  = "-"
	IDENTIFER_INTEGER       = ":"
	IDENTIFER_SIMPLE_STRING = "+"
	IDENTIFER_BULK_STRING   = "$"
	TERMINATOR              = "\r\n"
)

// ERRORS
var (
	INVALID_CHARACTERS_ERROR = errors.New("illegal character present: '\\r' or '\\n'")
	INVALID_INTEGER_ERROR    = errors.New("input is not an integer")
)

type CashewValue interface {
	GetValue() any
	Marshal() string
}

type SimpleString struct {
	value string
}

func NewSimpleString(s string) (*SimpleString, error) {
	for _, rune := range s {
		if rune == '\r' || rune == '\n' {
			return nil, INVALID_CHARACTERS_ERROR
		}
	}
	return &SimpleString{s}, nil
}

func (s SimpleString) GetValue() any {
	return s.value
}

func (s SimpleString) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFER_SIMPLE_STRING)
	b.WriteString(s.value)
	b.WriteString(TERMINATOR)
	return b.String()
}

type SimpleError struct {
	value string
}

func NewSimpleError(s string) (*SimpleError, error) {
	if hasInvalidCharacters(s) {
		return nil, INVALID_CHARACTERS_ERROR
	}

	return &SimpleError{s}, nil
}

func (s SimpleError) GetValue() any {
	return s.value
}

func (s SimpleError) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFER_SIMPLE_ERROR)
	b.WriteString(s.value)
	b.WriteString(TERMINATOR)
	return b.String()
}

type Integer struct {
	value int64
}

func (s Integer) GetValue() any {
	return s.value
}

func NewInteger(s string) (*Integer, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, INVALID_INTEGER_ERROR
	}
	return &Integer{n}, nil
}

func (i Integer) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFER_INTEGER)
	b.WriteString(strconv.FormatInt(i.value, 10))
	b.WriteString(TERMINATOR)
	return b.String()
}

type BulkString struct {
	value string
}

func NewBulkString(s string) (*BulkString, error) {
	return &BulkString{s}, nil
}

func (s BulkString) GetValue() any {
	return s.value
}

func (s BulkString) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFER_BULK_STRING)
	b.WriteString(strconv.FormatInt(int64(len(s.value)), 10))
	b.WriteString(TERMINATOR)
	b.WriteString(s.value)
	b.WriteString(TERMINATOR)
	return b.String()
}

func hasInvalidCharacters(s string) bool {
	for _, rune := range s {
		if rune == '\r' || rune == '\n' {
			return true
		}
	}

	return false
}

type Array struct {
	values []CashewValue
}

func NewArray(values []CashewValue) (*Array, error) {
	return &Array{values}, nil
}

func (a Array) GetValue() any {
	return a.values
}
