package validator

import (
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"strings"
)

const (
	validateTag            = "validate"
	validateTagSplitString = ";"
	validateLabelTag       = "label"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type ValidateError map[string]map[string]string

func (v ValidateError) Error() string {
	if b, err := json.Marshal(v); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}

type validator struct {
	validateLibrary ValidateLibrary
}

func NewValidator() *validator {
	return &validator{}
}

func (va *validator) SetValidateLibrary(validateLibrary ValidateLibrary) {
	va.validateLibrary = validateLibrary
}

func (va *validator) getTagMap(value reflect.Type) map[string][]string {

	result := map[string][]string{}

	for i := 0; i < value.NumField(); i++ {

		field := value.Field(i)

		tagString := field.Tag.Get(validateTag)

		if tagString == "" {
			continue
		}

		tagSlice := strings.Split(tagString, validateTagSplitString)

		result[field.Name] = tagSlice
	}
	return result
}

func (va *validator) RunValidators(i interface{}) ValidateError {

	if va.validateLibrary == nil {
		va.validateLibrary = simpleValidateLibrary
	}

	_value := reflect.ValueOf(i)

	_type := reflect.TypeOf(i)

	tagMap := va.getTagMap(_type)

	var result ValidateError

	for i := 0; i < _type.NumField(); i++ {

		field := _type.Field(i)

		methodName := fmt.Sprintf("Validate%s", field.Name)

		if method := _value.MethodByName(methodName); method.IsValid() {
			delete(tagMap, field.Name)

			err := method.Call([]reflect.Value{})[0]

			if !err.IsNil() {

				if item, ok := err.Interface().(map[string]string); !ok {

					errString := fmt.Sprintf("func `%s` must return `map[string]string` type", methodName)

					panic(errors.New(errString))

				} else if item != nil {

					if result == nil {
						result = make(ValidateError, 0)
					}

					label := field.Tag.Get(validateLabelTag)

					result[label] = item
				}
			}
		}
	}

	for k, v := range tagMap {

		field, _ := _type.FieldByName(k)

		label := field.Tag.Get(validateLabelTag)

		if label == "" {
			label = k
		}

		item := va.validateLibrary.Validate(_value.FieldByName(k), map[string][]string{k: v})

		if item != nil {

			if result == nil {
				result = make(ValidateError, 0)
			}

			result[label] = item
		}
	}
	return result
}
