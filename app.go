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

// make your RestView to become a gin.HandlerFunc
func AsHandlerFunc(r view.RestViewer) gin.HandlerFunc {

	// get reflect view value
	// it will finish before this application start
	// don't care about consuming too much performance
	value := reflect.ValueOf(r)

	return func(ctx *gin.Context) {

		defer handleError(r, ctx)

		if r.IsAllowRequest(ctx) {

			// Make request method from upper to title
			// GET -> get -> Get
			method := strings.Title(strings.ToLower(ctx.Request.Method))

			// try to find your handle in your view
			handle := value.MethodByName(method)

			// if we find
			// check if valid
			if handle.IsValid() {

				r.PrePare(ctx)
				context := reflect.ValueOf(ctx)
				handle.Call([]reflect.Value{context})
				r.Finish(ctx)
				return
			}
		} else {
			// handle forbidden request
			view.ForbiddenHandleFunc(ctx)
		}
		// if request is allowed
		// but don't have any handle to call
		view.MethodNotAllowedHandleFunc(ctx)
	}
}

// RenderTemplateView
// shortcut for template render
func RenderTemplateView(ctx *gin.Context, t view.TemplateViewer) {
	template := t.GetTemplateName(ctx)
	context := t.GetTemplateContext(ctx)
	ctx.HTML(200, template, context)
}
