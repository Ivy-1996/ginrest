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

// Constructor for validator
func NewValidator() *validator {
	return &validator{}
}

// Implement Validator
type validator struct {
	validateLibrary  ValidateLibrary
	validateTag      string
	validateSplitTag string
	validateLabelTag string
}

// Prepare for validator
func (va *validator) prepare() {
	if va.validateTag == empty {
		va.validateTag = validateTag
	}
	if va.validateLibrary == nil {
		va.validateLibrary = simpleValidateLibrary
	}
	if va.validateSplitTag == empty {
		va.validateSplitTag = validateTagSplitString
	}
	if va.validateLabelTag == empty {
		va.validateLabelTag = validateLabelTag
	}
}

// Set validateTag with this method
// Call it before RunValidators otherwise it won't be worked
func (va *validator) SetValidateTag(validateTag string) {
	va.validateTag = validateTag
}

// Set validateLibrary with this method
// Call it before RunValidators otherwise it won't be worked
func (va *validator) SetValidateLibrary(validateLibrary ValidateLibrary) {
	va.validateLibrary = validateLibrary
}

// Set validateSplitTag with this method
// Call it before RunValidators otherwise it won't be worked
func (va *validator) SetValidateSplitTag(validateSplitTag string) {
	va.validateSplitTag = validateSplitTag
}

// Set validateLabelTag with this method
// Call it before RunValidators otherwise it won't be worked
func (va *validator) SetValidateLabelTag(validateLabelTag string) {
	va.validateLabelTag = validateLabelTag
}

// Get validateTag value from struct
// You can call SetValidateTag to change tag name
func (va *validator) getTagMap(value reflect.Type) map[string][]string {

	result := map[string][]string{}

	for i := 0; i < value.NumField(); i++ {

		// Get current field
		field := value.Field(i)

		// Get validateTag value from field
		tagString := field.Tag.Get(va.validateTag)

		// Get, but no value
		// throw it
		if tagString == empty {
			continue
		}

		// Get different condition by validateSplitTag
		tagSlice := strings.Split(tagString, va.validateSplitTag)

		// Add it to result
		result[field.Name] = tagSlice
	}

	return result
}

// Implement Validator
func (va *validator) RunValidators(i interface{}) ValidateError {

	// Set some value default
	va.prepare()

	// Get reflect
	_value := reflect.ValueOf(i)

	_type := reflect.TypeOf(i)

	// TagMap or Condition ?
	tagMap := va.getTagMap(_type)

	var result ValidateError

	// Try to get validate function if we want validate field by ourselves
	for i := 0; i < _type.NumField(); i++ {

		// Get current Field
		field := _type.Field(i)

		// Get validate function name
		methodName := fmt.Sprintf("Validate%s", field.Name)

		// If we found and it is valid
		if method := _value.MethodByName(methodName); method.IsValid() {

			// Delete it from tagMap
			delete(tagMap, field.Name)

			// Get validate function return value
			// Note: validate function must return a single value and it must be map[string]string type
			err := method.Call([]reflect.Value{})[0]

			// We can call `IsNil` because it is a map
			// Pass this field if it return nil
			if !err.IsNil() {

				// Try to translate it to `map[string]string` type
				if item, ok := err.Interface().(map[string]string); !ok {

					// But we failed
					errString := fmt.Sprintf("func `%s` must return `map[string]string` type", methodName)

					panic(errors.New(errString))

				} else if item != nil {

					if result == nil {
						result = make(ValidateError, 0)
					}

					// Try to get label from field
					label := field.Tag.Get(va.validateLabelTag)

					// Set default label if it is empty
					if label == empty {
						label = field.Name
					}

					// Add it to result
					result[label] = item
				}
			}
		}
	}

	// Now it is time to validate by validateLibrary
	for k, v := range tagMap {

		// Get current field
		// It must be exist, so ignore if it exist
		field, _ := _type.FieldByName(k)

		label := field.Tag.Get(va.validateLabelTag)

		if label == empty {
			label = k
		}

		// Validate it by validateLibrary
		// Pass this field if it return nil, otherwise add it to result
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
