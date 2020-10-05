package ginrest

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/view"
	"reflect"
	"strings"
)

func handleError(r view.RestViewer, ctx *gin.Context) {
	err := recover()
	r.HandleError(ctx, err)
}

// Make your RestView become a gin.HandlerFunc
func AsHandlerFunc(r view.RestViewer, restHandlerFunc view.RestViewHandlerFunc) gin.HandlerFunc {

	// Get reflect view value
	// It will finish before this application start
	// Don't care about consuming too much performance
	value := reflect.ValueOf(r)

	return func(ctx *gin.Context) {

		defer handleError(r, ctx)

		if r.IsAllowRequest(ctx) {

			// Make request method from upper to title
			// Example: GET -> get -> Get
			method := strings.Title(strings.ToLower(ctx.Request.Method))

			// Try to find your handle in your view
			handler := value.MethodByName(method)

			// If we find
			// Check if valid
			if handler.IsValid() {
				// Call PrePare
				r.PrePare(ctx)
				// Call restHandlerFunc
				restHandlerFunc(handler, r, ctx)
				// Call Finish
				r.Finish(ctx)
			} else {
				// If request is allowed
				// But don't have any handle to call
				view.MethodNotAllowedHandle(ctx)
			}
		} else {
			// Handle forbidden request
			view.PermissionDeniedHandle(ctx)
		}
	}
}

// Shortcut for view.SimpleRestViewFunc
func AsSimpleHandleFunc(r view.RestViewer) gin.HandlerFunc {
	return AsHandlerFunc(r, view.SimpleRestViewFunc)
}

// Shortcut for view.DefaultRestViewFunc
func AsDefaultHandleFunc(r view.RestViewer) gin.HandlerFunc {
	return AsHandlerFunc(r, view.DefaultRestViewFunc)
}
