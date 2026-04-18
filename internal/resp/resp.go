package resp

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// VALUE TYPE IDENTIFIERS
const (
	IDENTIFIER_SIMPLE_ERROR  = "-"
	IDENTIFIER_INTEGER       = ":"
	IDENTIFIER_SIMPLE_STRING = "+"
	IDENTIFIER_BULK_STRING   = "$"
	IDENTIFIER_ARRAY         = "*"
	IDENTIFIER_NULL          = "_"
	TERMINATOR               = "\r\n"
)

// ERRORS
var (
	INVALID_CHARACTERS_ERROR     = errors.New("illegal character present: '\\r' or '\\n'")
	INVALID_INTEGER_ERROR        = errors.New("input is not an integer")
	INVALID_TERMINATION_SEQUENCE = errors.New("input not terminated by '\\r\\n'")
)

type CashewValue interface {
	GetValue() any
	Marshal() string
}

type Null struct {
	value any
}

func NewNull() (Null, error) {
	return Null{}, nil
}

func (n Null) GetValue() any {
	return n.value
}

func (n Null) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFIER_NULL)
	b.WriteString(TERMINATOR)
	return b.String()
}

type SimpleString struct {
	value string
}

func NewSimpleString(s string) (SimpleString, error) {
	for _, rune := range s {
		if rune == '\r' || rune == '\n' {
			return SimpleString{}, INVALID_CHARACTERS_ERROR
		}
	}
	return SimpleString{s}, nil
}

func (s SimpleString) GetValue() any {
	return s.value
}

func (s SimpleString) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFIER_SIMPLE_STRING)
	b.WriteString(s.value)
	b.WriteString(TERMINATOR)
	return b.String()
}

type SimpleError struct {
	value string
}

func NewSimpleError(s string) (SimpleError, error) {
	if hasInvalidCharacters(s) {
		return SimpleError{}, INVALID_CHARACTERS_ERROR
	}

	return SimpleError{s}, nil
}

func (s SimpleError) GetValue() any {
	return s.value
}

func (s SimpleError) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFIER_SIMPLE_ERROR)
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

type Integer struct {
	value int64
}

func (s Integer) GetValue() any {
	return s.value
}

func NewInteger(s string) (Integer, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return Integer{}, INVALID_INTEGER_ERROR
	}
	return Integer{n}, nil
}

func (i Integer) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFIER_INTEGER)
	b.WriteString(strconv.FormatInt(i.value, 10))
	b.WriteString(TERMINATOR)
	return b.String()
}

type BulkString struct {
	value string
}

func NewBulkString(s string) (BulkString, error) {
	return BulkString{s}, nil
}

func (s BulkString) GetValue() any {
	return s.value
}

func (s BulkString) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFIER_BULK_STRING)
	b.WriteString(strconv.FormatInt(int64(len(s.value)), 10))
	b.WriteString(TERMINATOR)
	b.WriteString(s.value)
	b.WriteString(TERMINATOR)
	return b.String()
}

type Array struct {
	values []CashewValue
}

func NewArray(values []CashewValue) (Array, error) {
	return Array{values}, nil
}

func (a Array) GetValue() any {
	return a.values
}

func (a Array) Marshal() string {
	var b strings.Builder
	b.WriteString(IDENTIFIER_ARRAY)
	b.WriteString(strconv.FormatInt(int64(len(a.values)), 10))
	b.WriteString(TERMINATOR)
	// Each element in the sequence adds its own termination
	for _, value := range a.values {
		b.WriteString(value.Marshal())
	}
	return b.String()
}

// UNMARSHAL

func Unmarshal(b *bufio.Reader) (CashewValue, error) {
	identifier, err := b.Peek(1)
	if err != nil {
		return nil, err
	}
	switch string(identifier) {
	case IDENTIFIER_SIMPLE_STRING:
		return unmarshalSimpleString(b)
	case IDENTIFIER_SIMPLE_ERROR:
		return unmarshalSimpleError(b)
	case IDENTIFIER_INTEGER:
		return unmarshalInteger(b)
	case IDENTIFIER_NULL:
		return unmarshalNull(b)
	case IDENTIFIER_ARRAY:
		isPossibleNull, err := checkIsPossibleNullSequence(b)
		if err != nil {
			return nil, err
		}
		if isPossibleNull {
			return unmarshalNull(b)
		}
		return unmarshalArray(b)
	case IDENTIFIER_BULK_STRING:
		isPossibleNull, err := checkIsPossibleNullSequence(b)
		if err != nil {
			return nil, err
		}
		if isPossibleNull {
			return unmarshalNull(b)
		}
		return unmarshalBulkString(b)
	default:
		return nil, fmt.Errorf("invalid data type identifier %s", identifier)
	}
}

func unmarshalSimpleString(b *bufio.Reader) (SimpleString, error) {
	identifier, err := readIdentifier(b)
	if err != nil {
		return SimpleString{}, err
	}
	if identifier != IDENTIFIER_SIMPLE_STRING {
		return SimpleString{}, fmt.Errorf("invalid data type identifier for simple string: %s", identifier)
	}

	s, err := readUntilTerminator(b)
	if err != nil {
		return SimpleString{}, err
	}

	return NewSimpleString(s)
}

func unmarshalSimpleError(b *bufio.Reader) (SimpleError, error) {
	identifier, err := readIdentifier(b)
	if err != nil {
		return SimpleError{}, err
	}
	if identifier != IDENTIFIER_SIMPLE_ERROR {
		return SimpleError{}, fmt.Errorf("invalid data type identifier for simple error: %s", identifier)
	}

	s, err := readUntilTerminator(b)
	if err != nil {
		return SimpleError{}, err
	}

	return NewSimpleError(s)
}

func unmarshalInteger(b *bufio.Reader) (Integer, error) {
	identifier, err := readIdentifier(b)
	if err != nil {
		return Integer{}, err
	}
	if identifier != IDENTIFIER_INTEGER {
		return Integer{}, fmt.Errorf("invalid data type identifier for integer: %s", identifier)
	}

	s, err := readUntilTerminator(b)
	if err != nil {
		return Integer{}, err
	}

	return NewInteger(s)
}

func unmarshalNull(b *bufio.Reader) (Null, error) {
	identifier, err := readIdentifier(b)
	if err != nil {
		return Null{}, err
	}
	if identifier != IDENTIFIER_NULL && identifier != IDENTIFIER_ARRAY && identifier != IDENTIFIER_BULK_STRING {
		return Null{}, fmt.Errorf("invalid data type identifier for null: %s", string(identifier))
	}

	s, err := readUntilTerminator(b)
	if err != nil {
		return Null{}, err
	}

	if identifier == IDENTIFIER_NULL {
		return Null{}, nil
	}

	s, found := strings.CutPrefix(s, "-1")
	if !found {
		return Null{}, fmt.Errorf("invalid null sequence for array or bulk string: %s", s)
	}

	return Null{}, nil
}

func unmarshalBulkString(b *bufio.Reader) (BulkString, error) {
	identifier, err := readIdentifier(b)
	if err != nil {
		return BulkString{}, err
	}
	if identifier != IDENTIFIER_BULK_STRING {
		return BulkString{}, fmt.Errorf("invalid data type identifier for bulk string: %s", identifier)
	}

	length, err := readUntilTerminator(b)
	if err != nil {
		return BulkString{}, err
	}
	lengthToRead, err := strconv.Atoi(length)
	if err != nil || lengthToRead < 0 {
		return BulkString{}, fmt.Errorf("invalid bulk string length: %s", length)
	}

	var sb strings.Builder
	for range lengthToRead {
		data, err := b.ReadByte()
		if err != nil {
			return BulkString{}, err
		}
		sb.WriteByte(data)
	}
	// Read two bytes to confirm termination sequence at the end
	terminationSequence := [2]byte{}
	for i := range len(terminationSequence) {
		read, err := b.ReadByte()
		if err != nil {
			return BulkString{}, err
		}
		terminationSequence[i] = read
	}
	if string(terminationSequence[:]) != TERMINATOR {
		return BulkString{}, fmt.Errorf("termination sequence not found after data, found: %q", string(terminationSequence[:]))
	}

	return BulkString{sb.String()}, err
}

func unmarshalArray(b *bufio.Reader) (Array, error) {
	identifier, err := readIdentifier(b)
	if err != nil {
		return Array{}, err
	}
	if identifier != IDENTIFIER_ARRAY {
		return Array{}, fmt.Errorf("invalid data type identifier for array: %s", identifier)
	}

	length, err := readUntilTerminator(b)
	if err != nil {
		return Array{}, err
	}
	elementCount, err := strconv.Atoi(length)
	if err != nil || elementCount < 0 {
		return Array{}, fmt.Errorf("invalid array length: %s", length)
	}

	var values []CashewValue
	for range elementCount {
		element, err := Unmarshal(b)
		if err != nil {
			return Array{}, err
		}
		values = append(values, element)
	}

	return Array{values}, nil
}

func checkIsPossibleNullSequence(b *bufio.Reader) (bool, error) {
	bytes, err := b.Peek(2)
	if err != nil {
		return false, err
	}
	sequence := string(bytes)
	if sequence == "$-" || sequence == "*-" {
		return true, nil
	}
	return false, nil
}

func readIdentifier(b *bufio.Reader) (string, error) {
	identifier, err := b.ReadByte()
	if err != nil {
		return "", err
	}
	return string(identifier), nil
}

func readUntilTerminator(b *bufio.Reader) (string, error) {
	data, err := b.ReadString('\n')
	if err != nil {
		return "", err
	}

	s, found := strings.CutSuffix(data, TERMINATOR)
	if !found {
		return "", INVALID_TERMINATION_SEQUENCE
	}

	return s, nil
}
