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
	conditionSplitTag      = ":"
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

// Get ValidatorNotes from struct
// You can call SetValidateTag to change tag name
func (va *validator) getValidatorNotes(value reflect.Type) validatorNodes {

	var nodes validatorNodes

	for i := 0; i < value.NumField(); i++ {

		var rules validatorRules

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

		for _, tag := range tagSlice {

			ruleAttrSlice := strings.Split(tag, conditionSplitTag)

			if len(ruleAttrSlice) != 3 {
				panic(`wrong tag style, example: gt:10:your value must greater than 10`)
			}

			name, expect, errMessage := ruleAttrSlice[0], ruleAttrSlice[1], ruleAttrSlice[2]

			rule := newValidatorRule(name, expect, errMessage)

			if rules == nil {
				rules = make(validatorRules, 0)
			}

			rules = append(rules, rule)
		}

		if nodes == nil {
			nodes = make(validatorNodes, 0)
		}

		node := newValidatorNote(field.Name, rules)

		nodes = append(nodes, node)
	}

	return nodes
}

func (va *validator) RunValidators(i interface{}) ValidateErrors {

	// Set some value default
	va.prepare()

	// Get reflect
	_value := reflect.ValueOf(i)

	_type := reflect.TypeOf(i)

	// Get rules nodes
	nodes := va.getValidatorNotes(_type)

	var validateErrors ValidateErrors

	// Try to get validate function if we want validate field by ourselves
	for i := 0; i < _type.NumField(); i++ {

		// Get current Field
		field := _type.Field(i)

		// Get validate function name
		methodName := fmt.Sprintf("Validate%s", field.Name)

		// If we found and it is valid
		if method := _value.MethodByName(methodName); method.IsValid() {

			// Remove it from nodes
			nodes = RemoveNodeByName(nodes, field.Name)

			// Get validate function return value
			// Note: validate function must return a single value and it must be ValidateNode
			err := method.Call([]reflect.Value{})[0]

			// We can call `IsNil` because it is a struct
			// Pass this field if it return nil
			if !err.IsNil() {

				// Try to translate it to `map[string]string` type
				if validateErrorNodes, ok := err.Interface().(*ValidateErrorNodes); !ok {

					// But we failed
					errString := fmt.Sprintf("func `%s` must return `*ValidateErrorNodes` type", methodName)

					panic(errors.New(errString))

				} else if validateErrorNodes != nil {

					if validateErrors == nil {
						validateErrors = make(ValidateErrors, 0)
					}

					// Try to get label from field
					label := field.Tag.Get(va.validateLabelTag)

					// Set default label if it is empty
					if label == empty {
						label = field.Name
					}
					validateError := NewValidateError(field.Name, validateErrorNodes)
					// Add it to result
					validateErrors = append(validateErrors, validateError)
				}
			}
		}
	}

	// Now it is time to validate by validateLibrary
	for _, node := range nodes {

		// Get current field
		// It must be exist, so ignore if it exist
		field, _ := _type.FieldByName(node.FieldName)

		label := field.Tag.Get(va.validateLabelTag)

		if label == empty {
			label = node.FieldName
		}

		// Validate it by validateLibrary
		// Pass this field if it return nil, otherwise add it to result
		validateErrorNodes := va.validateLibrary.Validate(_value.FieldByName(node.FieldName), node)

		if validateErrorNodes != nil {

			if validateErrors == nil {
				validateErrors = make(ValidateErrors, 0)
			}
			validateError := NewValidateError(node.FieldName, &validateErrorNodes)
			validateErrors = append(validateErrors, validateError)
		}
	}

	return validateErrors
}

type validatorRule struct {
	Name         string
	Expect       string
	ErrorMessage string
}

func newValidatorRule(name string, expect string, errorMessage string) *validatorRule {
	return &validatorRule{Name: name, Expect: expect, ErrorMessage: errorMessage}
}

type validatorRules []*validatorRule

type validatorNode struct {
	FieldName string
	Rules     validatorRules
}

func newValidatorNote(fieldName string, rules validatorRules) *validatorNode {
	return &validatorNode{FieldName: fieldName, Rules: rules}
}

type validatorNodes []*validatorNode

func RemoveNodeByName(nodes validatorNodes, name string) validatorNodes {
	for index, node := range nodes {
		if node.FieldName == name {
			nodes := append(nodes[:index], nodes[index+1:]...)
			return nodes
		}
	}
	return nodes
}
