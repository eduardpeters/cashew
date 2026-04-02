package commands_test

import (
	"strconv"
	"testing"
	"testing/synctest"
	"time"

	"github.com/eduardpeters/cashew/internal/commands"
	"github.com/eduardpeters/cashew/internal/resp"
	"github.com/eduardpeters/cashew/internal/store"
)

func TestHandleSet(t *testing.T) {
	tests := []struct {
		name     string
		input    []resp.CashewValue
		expected commands.Result
	}{
		{"replies with simple string OK",
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan")},
			commands.Result{mustNewSimpleString(t, "OK").Marshal(), false},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()

			r, err := commands.HandleSet(s, tt.input)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if r.Content != tt.expected.Content {
				t.Errorf("incorrect result content want %q, got %q", tt.expected.Content, r.Content)
			}
			if r.CloseConn != tt.expected.CloseConn {
				t.Errorf("incorrect close connection value want %v, got %v", tt.expected.CloseConn, r.CloseConn)
			}

			key := tt.input[0].(resp.BulkString)
			stored, err := s.Get(key)
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}

			if stored.GetValue() != tt.input[1].GetValue() {
				t.Errorf("incorrect stored value want %q, got %q", tt.input[1].GetValue(), stored.GetValue())
			}
		})
	}
}

func TestHandleSetInvalidArguments(t *testing.T) {
	tests := []struct {
		name  string
		input []resp.CashewValue
	}{
		{"Fails on empty argument",
			[]resp.CashewValue{},
		},
		{"Fails on single argument",
			[]resp.CashewValue{mustNewBulkString(t, "key")},
		},
		{"Fails on non bulk string arguments - first simple string",
			[]resp.CashewValue{mustNewSimpleString(t, "key"), mustNewBulkString(t, "value")},
		},
		{"Fails on non bulk string arguments - second simple string",
			[]resp.CashewValue{mustNewBulkString(t, "key"), mustNewSimpleString(t, "value")},
		},
		{"Fails on non bulk string arguments - all simple string",
			[]resp.CashewValue{mustNewSimpleString(t, "key"), mustNewSimpleString(t, "value")},
		},
		{"Fails on missing expiry seconds",
			[]resp.CashewValue{mustNewBulkString(t, "key"), mustNewBulkString(t, "value"), mustNewBulkString(t, "EX")},
		},
		{"Fails on negative expiry seconds",
			[]resp.CashewValue{mustNewBulkString(t, "key"), mustNewBulkString(t, "value"), mustNewBulkString(t, "EX"), mustNewBulkString(t, "-1")},
		},
		{"Fails on unknown option",
			[]resp.CashewValue{mustNewBulkString(t, "key"), mustNewBulkString(t, "value"), mustNewBulkString(t, "expire"), mustNewBulkString(t, "1")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := store.NewStore()

			_, err := commands.HandleSet(s, tt.input)
			if err == nil {
				t.Errorf("Expected error for %q - got: %v", tt.input, err)
			}
		})
	}
}

func TestHandleSetEX(t *testing.T) {
	tests := []struct {
		name          string
		input         []resp.CashewValue
		expireSeconds int
	}{
		{"sets an expiration for 1 second",
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan"), mustNewBulkString(t, "EX"), mustNewBulkString(t, "1")},
			1,
		},
		{"sets an expiration for 5 seconds",
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan"), mustNewBulkString(t, "EX"), mustNewBulkString(t, "5")},
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := store.NewStore()

				_, err := commands.HandleSet(s, tt.input)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				time.Sleep(time.Second*time.Duration(tt.expireSeconds-1) + time.Millisecond)

				key := tt.input[0].(resp.BulkString)
				stored, err := s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != tt.input[1].GetValue() {
					t.Fatalf("incorrect stored value want %q, got %q", tt.input[1].GetValue(), stored.GetValue())
				}

				time.Sleep(time.Second*time.Duration(tt.expireSeconds) + time.Millisecond)

				stored, err = s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != nil {
					t.Errorf("incorrect value after expiration want %v, got %q", nil, stored.GetValue())
				}
			})
		})
	}
}

func TestHandleSetPX(t *testing.T) {
	tests := []struct {
		name         string
		input        []resp.CashewValue
		expireMillis int
	}{
		{"sets an expiry for 512 milliseconds",
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan"), mustNewBulkString(t, "PX"), mustNewBulkString(t, "512")},
			512,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := store.NewStore()

				_, err := commands.HandleSet(s, tt.input)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				time.Sleep(time.Millisecond*time.Duration(tt.expireMillis-100) + time.Millisecond)

				key := tt.input[0].(resp.BulkString)
				stored, err := s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != tt.input[1].GetValue() {
					t.Fatalf("incorrect stored value want %q, got %q", tt.input[1].GetValue(), stored.GetValue())
				}

				time.Sleep(time.Millisecond * time.Duration(tt.expireMillis+1))

				stored, err = s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != nil {
					t.Errorf("incorrect value after expiration want %v, got %q", nil, stored.GetValue())
				}
			})
		})
	}
}

func TestHandleSetEXAT(t *testing.T) {
	tests := []struct {
		name                string
		input               []resp.CashewValue
		expireSecondsOffset int
	}{
		{"sets an expiry for a timestamp 5 seconds forward",
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan"), mustNewBulkString(t, "EXAT")},
			5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := store.NewStore()

				expiration := time.Now().Add(time.Second * time.Duration(tt.expireSecondsOffset)).Unix()
				cmds := append(tt.input, mustNewBulkString(t, strconv.Itoa(int(expiration))))
				_, err := commands.HandleSet(s, cmds)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				time.Sleep(time.Second*time.Duration(tt.expireSecondsOffset-1) + time.Millisecond)

				key := tt.input[0].(resp.BulkString)
				stored, err := s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != tt.input[1].GetValue() {
					t.Fatalf("incorrect stored value want %q, got %q", tt.input[1].GetValue(), stored.GetValue())
				}

				time.Sleep(time.Second * time.Duration(tt.expireSecondsOffset+1))

				stored, err = s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != nil {
					t.Errorf("incorrect value after expiration want %v, got %q", nil, stored.GetValue())
				}
			})
		})
	}
}

func TestHandleSetPXAT(t *testing.T) {
	tests := []struct {
		name               string
		input              []resp.CashewValue
		expireMillisOffset int
	}{
		{"sets an expiry for a timestamp 500 milliseconds forward",
			[]resp.CashewValue{mustNewBulkString(t, "name"), mustNewBulkString(t, "juan"), mustNewBulkString(t, "PXAT")},
			500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				s := store.NewStore()

				expiration := time.Now().Add(time.Millisecond * time.Duration(tt.expireMillisOffset)).UnixMilli()
				cmds := append(tt.input, mustNewBulkString(t, strconv.Itoa(int(expiration))))
				_, err := commands.HandleSet(s, cmds)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				time.Sleep(time.Millisecond*time.Duration(tt.expireMillisOffset-100) + time.Millisecond)

				key := tt.input[0].(resp.BulkString)
				stored, err := s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != tt.input[1].GetValue() {
					t.Fatalf("incorrect stored value want %q, got %q", tt.input[1].GetValue(), stored.GetValue())
				}

				time.Sleep(time.Millisecond * time.Duration(tt.expireMillisOffset+1))

				stored, err = s.Get(key)
				if err != nil {
					t.Fatalf("Unexpected error %v", err)
				}

				if stored.GetValue() != nil {
					t.Errorf("incorrect value after expiration want %v, got %q", nil, stored.GetValue())
				}
			})
		})
	}
}
