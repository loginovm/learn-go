package validators

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrIncorrectLen          = errors.New("string length does not equal specified length")
	ErrRegexpNotMatched      = errors.New("string does not match pattern")
	ErrValueNotInSet         = errors.New("value is not in specified range")
	ErrLessThanMin           = errors.New("value should be equal or greater than min value")
	ErrGreaterThanMax        = errors.New("value should be less or equal than max value")
	ErrValidatorNotSupported = errors.New("validator not supported")
)

type Validator struct {
	Name  string
	Value string
}

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return fmt.Sprintf(`FieldName:%s, ValidationMessage:%s`, v.Field, v.Err)
}

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}
	sb.Grow(len(v))
	for _, i := range v {
		sb.WriteString(i.Error())
		sb.WriteByte(13)
	}
	return sb.String()
}
