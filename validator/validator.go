package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	validateTag            = "validate"
	validateTagSplitString = ";"
	validateLabelTag       = "label"
)

type validator struct {
	validateLibrary ValidateLibrary
}

func NewValidator() *validator {
	return &validator{}
}

func (v *validator) SetValidateLibrary(validateLibrary ValidateLibrary) {
	v.validateLibrary = validateLibrary
}

func (v *validator) getTagMap(value reflect.Type) map[string][]string {

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

func (va *validator) RunValidators(i interface{}) map[string]map[string]string {

	if va.validateLibrary == nil {
		va.validateLibrary = simpleValidateLibrary
	}

	_value := reflect.ValueOf(i)

	_type := reflect.TypeOf(i)

	tagMap := va.getTagMap(_type)

	var result map[string]map[string]string

	for k, v := range tagMap {

		field, _ := _type.FieldByName(k)

		label := field.Tag.Get(validateLabelTag)

		if label == "" {
			label = k
		}

		var item map[string]string

		methodName := fmt.Sprintf("Validate%s", k)

		method := _value.MethodByName(methodName)

		if method.IsValid() {

			args := reflect.ValueOf(map[string][]string{k: v})

			err := method.Call([]reflect.Value{args})[0]

			if !err.IsNil() {

				var ok bool

				item, ok = err.Interface().(map[string]string)

				if !ok {
					errString := fmt.Sprintf("func `%s` must return `map[string]string` type", methodName)
					panic(errors.New(errString))
				}
			}

		} else {
			item = va.validateLibrary.Validate(_value.FieldByName(k), map[string][]string{k: v})
		}

		if item != nil {

			if result == nil {
				result = make(map[string]map[string]string, 0)
			}

			result[label] = item
		}
	}
	return result
}
