package validator

import (
	"fmt"
	"testing"
)

type A struct {
	Age1 int `validate:"gt:10:Age1必须大于10;lt:20:Age1必须小于20"`
	Age2 int `validate:"gt:10:必须大于10"`
}

func (a A) ValidateAge2() *ValidateErrorNodes {
	validateErrorNodes := make(ValidateErrorNodes, 0)
	if a.Age2 < 10 {
		validateErrorNode := NewValidateErrorNode("no_field", "Age2必须大于10")
		return NewValidateErrorNodes(validateErrorNode)
	} else if a.Age2 > 20 {
		validateErrorNode := NewValidateErrorNode("no_field", "Age2必须小于20")
		validateErrorNodes = append(validateErrorNodes, validateErrorNode)
		return &validateErrorNodes
	}
	return nil
}

func TestValidate(t *testing.T) {
	validate := NewValidator()
	err := validate.RunValidators(A{Age1: 9, Age2: 9})
	if err != nil {
		fmt.Println(err.Error())
	}
}
