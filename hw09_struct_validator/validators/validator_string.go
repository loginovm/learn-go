package validators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ValidateString(field string, value string, v Validator) error {
	switch v.Name {
	default:
		return fmt.Errorf("%w: %s", ErrValidatorNotSupported, v.Name)
	case "len":
		return validateLen(field, value, v)
	case "regexp":
		return validateRegexp(field, value, v)
	case "in":
		return validateStringIn(field, value, v)
	}
}

func validateLen(field string, value string, v Validator) error {
	length, err := strconv.Atoi(v.Value)
	if err != nil {
		return fmt.Errorf(`failed to parse len validator: "%s"`, v.Value)
	}
	if len(value) != length {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf(`%w: %d`, ErrIncorrectLen, length),
		}
	}

	return nil
}

func validateRegexp(field string, value string, v Validator) error {
	r, err := regexp.Compile(v.Value)
	if err != nil {
		return err
	}
	if !r.MatchString(value) {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf(`%w: "%s"`, ErrRegexpNotMatched, v.Value),
		}
	}

	return nil
}

func validateStringIn(field string, value string, v Validator) error {
	var isValInRange bool
	for _, v := range strings.Split(v.Value, ",") {
		v = strings.TrimSpace(v)
		if value == v {
			isValInRange = true
			break
		}
	}
	if !isValInRange {
		return ValidationError{
			Field: field,
			Err:   ErrValueNotInSet,
		}
	}

	return nil
}
