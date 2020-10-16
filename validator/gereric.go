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

		// Get validateFunc from library
		validateFunc := s.LookForValidateFunc(rule.Name)

		// If validated, but no passed
		if err := validateFunc(field, rule); err != nil {
			if validateNodes == nil {
				validateNodes = make(ValidateErrorNodes, 0)
			}
			// Create new ErrorNode to validateNodes
			validateNode := NewValidateErrorNode(rule.Name, err.Error())
			validateNodes = append(validateNodes, validateNode)
		}
	}

	// Return result
	// if this field passed, validateNodes must be nil type
	return validateNodes

}

var simpleValidateLibrary SimpleValidateLibrary

// required valid
// Yeah! just like that
func required(field reflect.Value, rule *validatorRule) error {

	isValid := func(f reflect.Value) bool {
		// todo add more support type here
		// help!
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

	if !isValid(field) {
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

// Support int and string only
// Int: greater than expect
// String: length of string greater than expect
func gt(field reflect.Value, rule *validatorRule) error {
	kink := field.Kind()
	expect := parseInt(rule.Expect)
	if reflect.Int <= kink && kink <= reflect.Float64 {
		if field.Int() > int64(expect) {
			return nil
		} else {
			return errors.New(rule.ErrorMessage)
		}
	} else if kink == reflect.String {
		if len(field.String()) < expect {
			return errors.New(rule.ErrorMessage)
		} else {
			return nil
		}
	} else {
		return errors.New(rule.ErrorMessage)
	}
}

// Support int and string only
// Int: lesser than expect
// String: length of string lesser than expect
func lt(field reflect.Value, rule *validatorRule) error {
	kink := field.Kind()
	expect := parseInt(rule.Expect)
	if reflect.Int <= kink && kink <= reflect.Float64 {
		if field.Int() < int64(expect) {
			return nil
		} else {
			return errors.New(rule.ErrorMessage)
		}
	} else if kink == reflect.String {
		if len(field.String()) > expect {
			return errors.New(rule.ErrorMessage)
		} else {
			return nil
		}
	} else {
		return errors.New(rule.ErrorMessage)
	}
}

// Support string only
// Check if it is a legal email
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

// Support string only
// Check if it is a legal uuid
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
	simpleValidateLibrary = make(SimpleValidateLibrary, 0)
	simpleValidateLibrary.Register("required", required)
	simpleValidateLibrary.Register("gt", gt)
	simpleValidateLibrary.Register("lt", lt)
	simpleValidateLibrary.Register("email", email)
	simpleValidateLibrary.Register("uuid", uuid)
}
