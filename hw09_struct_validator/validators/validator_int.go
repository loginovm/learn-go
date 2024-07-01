package validators

import (
	"fmt"
	"strconv"
	"strings"
)

func ValidateInt(field string, value int, v Validator) error {
	switch v.Name {
	default:
		return fmt.Errorf("%w: %s", ErrValidatorNotSupported, v.Name)
	case "min":
		return validateMin(field, value, v)
	case "max":
		return validateMax(field, value, v)
	case "in":
		return validateIntIn(field, value, v)
	}
}

func validateMin(field string, value int, v Validator) error {
	min, err := strconv.Atoi(v.Value)
	if err != nil {
		return fmt.Errorf(`failed to parse min validator: "%s"`, v.Value)
	}
	if value < min {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w: %d", ErrLessThanMin, min),
		}
	}

	return nil
}

func validateMax(field string, value int, v Validator) error {
	max, err := strconv.Atoi(v.Value)
	if err != nil {
		return fmt.Errorf(`failed to parse max validator: "%s"`, v.Value)
	}
	if value > max {
		return ValidationError{
			Field: field,
			Err:   fmt.Errorf("%w: %d", ErrGreaterThanMax, max),
		}
	}

	return nil
}

func validateIntIn(field string, value int, v Validator) error {
	var isValInRange bool
	for _, s := range strings.Split(v.Value, ",") {
		d, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			return fmt.Errorf(`failed to parse in validator: "%s"`, v.Value)
		}
		if value == d {
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
