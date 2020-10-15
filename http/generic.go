package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/serializer"
)

func ValidateHandler(ctx *gin.Context, serializer serializer.Serializer) {
	var err error
	// validate data
	if err = serializer.Validate(ctx); err == nil {
		serializer.OnSuccess(ctx)
	} else {
		serializer.OnFail(ctx, err)
	}
}
