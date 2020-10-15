package validator

import (
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ValidateError map[string]map[string]string

type ValidateFunc func(reflect.Value, []string) error

// Implement error
func (v ValidateError) Error() string {
	if b, err := json.Marshal(v); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}

type Validator interface {
	RunValidators(interface{}) ValidateError
}
