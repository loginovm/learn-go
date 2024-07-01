package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	vl "github.com/loginovm/learn-go/hw09_struct_validator/validators"
)

var ErrValueIsNotStruct = errors.New("value is not a struct")

func Validate(v interface{}) error {
	var vErr vl.ValidationErrors
	rVal := reflect.ValueOf(v)
	if rVal.Kind() != reflect.Struct {
		return ErrValueIsNotStruct
	}

	rType := reflect.TypeOf(v)
	for i := 0; i < rType.NumField(); i++ {
		if tv, ok := rType.Field(i).Tag.Lookup("validate"); ok {
			v := parseValidator(tv)
			err := validateField(rType.Field(i).Name, rVal.Field(i), v)

			ve, ok := checkErrs(err)
			if ok {
				vErr = append(vErr, ve...)
			} else if err != nil {
				return err
			}
		}
	}

	return vErr
}

func validateField(field string, value reflect.Value, validators []vl.Validator) error {
	var vErr vl.ValidationErrors
	for _, v := range validators {
		var err error
		switch value.Kind() { //nolint:exhaustive
		default:
			return nil
		case reflect.String:
			err = vl.ValidateString(field, value.String(), v)
		case reflect.Int:
			err = vl.ValidateInt(field, int(value.Int()), v)
		case reflect.Array:
		case reflect.Slice:
			for i := 0; i < value.Len(); i++ {
				item := value.Index(i)
				fieldName := fmt.Sprintf("%s[%d]", field, i)
				err2 := validateField(fieldName, item, []vl.Validator{v})
				ve, ok := checkErrs(err2)
				if ok {
					vErr = append(vErr, ve...)
				} else if err2 != nil {
					return err2
				}
			}
		}

		ve, ok := checkErr(err)
		if ok {
			vErr = append(vErr, ve)
		} else if err != nil {
			return err
		}
	}

	if len(vErr) == 0 {
		return nil
	}

	return vErr
}

func checkErr(err error) (vl.ValidationError, bool) {
	var t vl.ValidationError
	if errors.As(err, &t) {
		return t, true
	}
	return vl.ValidationError{}, false
}

func checkErrs(err error) (vl.ValidationErrors, bool) {
	var t vl.ValidationErrors
	if errors.As(err, &t) {
		return t, true
	}
	return nil, false
}

func parseValidator(vs string) []vl.Validator {
	validators := make([]vl.Validator, 0)
	parts := strings.Split(vs, "|")
	for _, p := range parts {
		valParts := strings.Split(p, ":")
		if len(valParts) > 1 {
			v := vl.Validator{
				Name:  strings.TrimSpace(valParts[0]),
				Value: strings.TrimSpace(valParts[1]),
			}
			validators = append(validators, v)
		}
	}
	return validators
}
