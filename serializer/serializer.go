package serializer

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/validator"
)

type Serializer interface {
	// Validate request with this method
	Validate(ctx *gin.Context, i interface{}) error
	// OnSuccess will be called if Validate return nil
	// Handle allow request with this method
	OnSuccess(ctx *gin.Context)
	// OnFail will be called if Validate return not nil
	// Handle deny request with this method
	OnFail(ctx *gin.Context, err error)
}

type SimpleJsonSerializer struct{}

func (s SimpleJsonSerializer) Validate(ctx *gin.Context, i interface{}) error {
	buff := make([]byte, 0)

	var (
		total int
		err   error
	)

	// Load json error
	if total, err = ctx.Request.Body.Read(buff); err != nil {
		panic(err)
	}

	// Parse json error
	if err := json.Unmarshal(buff[0:total], &i); err != nil {
		panic(err)
	}

	validate := validator.NewValidator()

	if err := validate.RunValidators(i); err != nil {
		return err
	}

	return nil
}

func (s SimpleJsonSerializer) OnSuccess(ctx *gin.Context) {
	panic("implement me")
}

func (s SimpleJsonSerializer) OnFail(ctx *gin.Context, err error) {

	if err, ok := err.(validator.ValidateErrors); ok {

		var result validator.ValidateErrors

		if err := json.Unmarshal([]byte(err.Error()), &result); err != nil {
			panic(err)
		}

		ctx.JSON(400, result)

	} else {

		ctx.JSON(400, nil)

	}
}
