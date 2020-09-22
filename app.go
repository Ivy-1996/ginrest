package ginrest

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/view"
	"reflect"
	"strings"
)

func AsHandlerFunc(r view.RestViewer) gin.HandlerFunc {

	value := reflect.ValueOf(r)

	return func(ctx *gin.Context) {

		defer func() {
			err := recover()
			r.HandleError(ctx, err)
		}()

		if r.IsAllowRequest(ctx) {
			method := strings.Title(strings.ToLower(ctx.Request.Method))
			handle := value.MethodByName(method)
			if handle.IsValid() {
				r.PrePare(ctx)
				context := reflect.ValueOf(ctx)
				handle.Call([]reflect.Value{context})
				r.Finish(ctx)
				return
			}
		} else {
			view.ForbiddenHandleFunc(ctx)
		}
		view.MethodNotAllowedHandleFunc(ctx)
	}
}
