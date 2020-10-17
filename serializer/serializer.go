package serializer

import (
	"github.com/gin-gonic/gin"
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
