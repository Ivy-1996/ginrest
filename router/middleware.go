package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/validator"
	"github.com/ivy-1996/ginrest/view"
)

func InterRequired(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exist := ctx.Params.Get(key)
		if exist && validator.InterRequired(value) {
			ctx.Next()
		} else {
			view.PageNotFoundHandle(ctx)
			ctx.Abort()
		}
	}
}

func UuidRequired(key string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		value, exist := ctx.Params.Get(key)
		if exist && validator.UuidRequired(value) {
			ctx.Next()
		} else {
			view.PageNotFoundHandle(ctx)
			ctx.Abort()
		}
	}
}

func LoginRequired(validFunc func(ctx *gin.Context) bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if validFunc(ctx) {
			ctx.Next()
		} else {
			view.NotAuthenticatedHandle(ctx)
			ctx.Abort()
		}
	}
}
