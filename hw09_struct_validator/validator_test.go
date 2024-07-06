package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	vld "github.com/loginovm/learn-go/hw09_struct_validator/validators"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidateFail(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr vld.ValidationErrors
	}{
		{
			in: User{ID: "1", Name: "user", Age: 10, Email: "user", Role: "manager", Phones: []string{"111", "222"}},
			expectedErr: vld.ValidationErrors{
				vld.ValidationError{Field: "ID", Err: vld.ErrIncorrectLen},
				vld.ValidationError{Field: "Age", Err: vld.ErrLessThanMin},
				vld.ValidationError{Field: "Email", Err: vld.ErrRegexpNotMatched},
				vld.ValidationError{Field: "Role", Err: vld.ErrValueNotInSet},
				vld.ValidationError{Field: "Phones[0]", Err: vld.ErrIncorrectLen},
				vld.ValidationError{Field: "Phones[1]", Err: vld.ErrIncorrectLen},
			},
		},
		{
			in: App{Version: "1.0"},
			expectedErr: vld.ValidationErrors{
				vld.ValidationError{Field: "Version", Err: vld.ErrIncorrectLen},
			},
		},
		{
			in: Response{Code: 201, Body: "aaa"},
			expectedErr: vld.ValidationErrors{
				vld.ValidationError{Field: "Code", Err: vld.ErrValueNotInSet},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			require.IsType(t, vld.ValidationErrors{}, err)
			actual := err.(vld.ValidationErrors) //nolint:errorlint
			require.Len(t, actual, len(tt.expectedErr))
			require.True(t, areErrorsSame(actual, tt.expectedErr), "Errors differ from expected")
		})
	}
}

func TestValidateSuccess(t *testing.T) {
	tests := []any{
		User{
			ID:     "e0b000e9-3055-4675-b7a6-78cc8b7879cb",
			Name:   "user",
			Age:    18,
			Email:  "user@abc.com",
			Role:   "stuff",
			Phones: []string{"111-111-111", "222-222-222"},
		},
		User{
			ID:     "e0b000e9-3055-4675-b7a6-78cc8b7879cb",
			Name:   "user",
			Age:    50,
			Email:  "user@abc.com",
			Role:   "admin",
			Phones: []string{"111-111-111"},
		},
		App{Version: "0.0.1"},
		Response{Code: 200, Body: "aaa"},
		Token{
			Header:    []byte{0x00, 0x11},
			Payload:   []byte{0x02},
			Signature: []byte{0x03},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt)
			require.Nil(t, err)
		})
	}
}

func areErrorsSame(actual, expected vld.ValidationErrors) bool {
	for _, e := range expected {
		errFound := false
		for _, a := range actual {
			if e.Field == a.Field && errors.Is(a.Err, e.Err) {
				errFound = true
				break
			}
		}
		if !errFound {
			return false
		}
	}
	return true
}
