package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const empty = ""

type ValidateFunc func(reflect.Value, []string) error

type ValidateLibrary interface {
	Register(key string, validateFunc ValidateFunc)
	Validate(field reflect.Value, arg map[string][]string) map[string]string
	LookForValidateFunc(key string) ValidateFunc
}

type SimpleValidateLibrary map[string]ValidateFunc

func (s SimpleValidateLibrary) LookForValidateFunc(key string) ValidateFunc {
	validateFunc := s[key]

	if validateFunc == nil {
		panic(errors.New(fmt.Sprintf("no key named `%s`, did your forgetten to register it in your library?", key)))
	}
	return validateFunc
}

func (s SimpleValidateLibrary) Register(key string, validateFunc ValidateFunc) {
	if s[key] != nil {
		panic(fmt.Sprintf("key `%s` has already exist", key))
	}
	s[key] = validateFunc
}

func (s SimpleValidateLibrary) Validate(field reflect.Value, arg map[string][]string) map[string]string {

	var errData map[string]string

	for _, v := range arg {

		for _, tag := range v {

			item := strings.Split(tag, ":")

			if len(item) != 3 {
				panic(errors.New("wrong tag style, example: `gt:10:must great than 10`"))
			}

			validateFunc := s.LookForValidateFunc(item[0])

			if err := validateFunc(field, item); err != nil {

				if errData == nil {
					errData = make(map[string]string, 0)
				}

				errData[item[0]] = err.Error()
			}
		}
	}

	return errData
}

var simpleValidateLibrary SimpleValidateLibrary

func required(field reflect.Value, s []string) error {

	lambda := func(f reflect.Value) bool {
		// todo add more type support here
		switch field.Kind() {
		case reflect.Int:
			return !f.IsZero()
		case reflect.String:
			return f.String() != empty
		case reflect.Bool:
			return f.Bool()
		default:
			return f.IsNil()
		}
	}

	if !lambda(field) {
		return errors.New(s[2])
	}
	return nil
}

func parseInt(s string) int {
	expect, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return expect
}

func gt(field reflect.Value, s []string) error {
	kink := field.Kind()
	if 2 <= kink && kink <= 14 {
		expect := parseInt(s[1])
		if field.Int() > int64(expect) {
			return nil
		} else {
			return errors.New(s[2])
		}
	} else {
		return errors.New(s[2])
	}
}

func lt(field reflect.Value, s []string) error {
	kink := field.Kind()
	if 2 <= kink && kink <= 14 {
		expect := parseInt(s[1])
		if field.Int() < int64(expect) {
			return nil
		} else {
			return errors.New(s[2])
		}
	} else {
		return errors.New(s[2])
	}
}

func email(field reflect.Value, s []string) error {
	switch field.Kind() {
	case reflect.String:
		if EmailRequired(field.String()) {
			return nil
		}
		return errors.New(s[2])
	default:
		return errors.New(s[2])
	}
}

func uuid(field reflect.Value, s []string) error {
	switch field.Kind() {
	case reflect.String:
		if UuidRequired(field.String()) {
			return nil
		}
		return errors.New(s[2])
	default:
		return errors.New(s[2])
	}
}

func init() {
	simpleValidateLibrary = make(map[string]ValidateFunc, 0)
	simpleValidateLibrary.Register("required", required)
	simpleValidateLibrary.Register("gt", gt)
	simpleValidateLibrary.Register("lt", gt)
	simpleValidateLibrary.Register("email", email)
	simpleValidateLibrary.Register("uuid", uuid)
}
