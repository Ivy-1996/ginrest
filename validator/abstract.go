package validator

import (
	jsoniter "github.com/json-iterator/go"
	"reflect"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// ValidateLibrary interface
type ValidateLibrary interface {
	Register(key string, validateFunc ValidateFunc)
	Validate(field reflect.Value, node *validatorNode) ValidateErrorNodes
	LookForValidateFunc(key string) ValidateFunc
}

// Validator interface
type Validator interface {
	RunValidators(interface{}) ValidateError
}

// Validate data if it is valid
// ValidateFunc must register to ValidateLibrary or it won't be worked
type ValidateFunc func(reflect.Value, *validatorRule) error

type ValidateErrorNode struct {
	Code         string `json:"code"`
	ErrorMessage string `json:"error_message"`
}

type ValidateErrorNodes []*ValidateErrorNode

type ValidateError struct {
	FieldName          string              `json:"field_name"`
	ValidateErrorNodes *ValidateErrorNodes `json:"validate_error_nodes"`
}

type ValidateErrors []*ValidateError

// Constructor for ValidateErrorNode
func NewValidateErrorNode(code string, errorMessage string) *ValidateErrorNode {
	return &ValidateErrorNode{Code: code, ErrorMessage: errorMessage}
}

// Constructor for ValidateErrorNodes
func NewValidateErrorNodes(validateErrorNode ...*ValidateErrorNode) *ValidateErrorNodes {
	return (*ValidateErrorNodes)(&validateErrorNode)
}

// Constructor for ValidateError
func NewValidateError(fieldName string, validateErrorNodes *ValidateErrorNodes) *ValidateError {
	return &ValidateError{FieldName: fieldName, ValidateErrorNodes: validateErrorNodes}
}

// Implement error
func (v ValidateErrors) Error() string {
	if b, err := json.Marshal(v); err != nil {
		panic(err)
	} else {
		return string(b)
	}
}
