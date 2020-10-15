package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const empty = ""

// Implement ValidateLibrary
type SimpleValidateLibrary map[string]ValidateFunc

// Get ValidateFunc by key
// It would course panic if key does not exist
func (s SimpleValidateLibrary) LookForValidateFunc(key string) ValidateFunc {
	validateFunc := s[key]

	if validateFunc == nil {
		panic(errors.New(fmt.Sprintf("no key named `%s`, did your forgetten to register it in your library?", key)))
	}
	return validateFunc
}

// Register your ValidateFunc to Library
// It would course panic if key has already exist
func (s SimpleValidateLibrary) Register(key string, validateFunc ValidateFunc) {
	if s[key] != nil {
		panic(fmt.Sprintf("key `%s` has already exist", key))
	}
	s[key] = validateFunc
}

// Validate field if is legal
func (s SimpleValidateLibrary) Validate(field reflect.Value, node *validatorNode) ValidateErrorNodes {

	var validateNodes ValidateErrorNodes

	for _, rule := range node.Rules {

		validateFunc := s.LookForValidateFunc(rule.Name)

		if err := validateFunc(field, rule); err != nil {
			if validateNodes == nil {
				validateNodes = make(ValidateErrorNodes, 0)
			}
			validateNode := NewValidateErrorNode(rule.Name, err.Error())
			validateNodes = append(validateNodes, validateNode)
		}
	}
	return validateNodes

}

var simpleValidateLibrary SimpleValidateLibrary

func required(field reflect.Value, rule *validatorRule) error {

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
		return errors.New(rule.ErrorMessage)
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

func gt(field reflect.Value, rule *validatorRule) error {
	kink := field.Kind()
	if 2 <= kink && kink <= 14 {
		expect := parseInt(rule.Expect)
		if field.Int() > int64(expect) {
			return nil
		} else {
			return errors.New(rule.ErrorMessage)
		}
	} else {
		return errors.New(rule.ErrorMessage)
	}
}

func lt(field reflect.Value, rule *validatorRule) error {
	kink := field.Kind()
	if 2 <= kink && kink <= 14 {
		expect := parseInt(rule.Expect)
		if field.Int() < int64(expect) {
			return nil
		} else {
			return errors.New(rule.ErrorMessage)
		}
	} else {
		return errors.New(rule.ErrorMessage)
	}
}

func email(field reflect.Value, rule *validatorRule) error {
	switch field.Kind() {
	case reflect.String:
		if EmailRequired(field.String()) {
			return nil
		}
		return errors.New(rule.ErrorMessage)
	default:
		return errors.New(rule.ErrorMessage)
	}
}

func uuid(field reflect.Value, rule *validatorRule) error {
	switch field.Kind() {
	case reflect.String:
		if UuidRequired(field.String()) {
			return nil
		}
		return errors.New(rule.ErrorMessage)
	default:
		return errors.New(rule.ErrorMessage)
	}
}

func init() {
	simpleValidateLibrary = make(map[string]ValidateFunc, 0)
	simpleValidateLibrary.Register("required", required)
	simpleValidateLibrary.Register("gt", gt)
	simpleValidateLibrary.Register("lt", lt)
	simpleValidateLibrary.Register("email", email)
	simpleValidateLibrary.Register("uuid", uuid)
}
