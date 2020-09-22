package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	methodNotAllowedResponseBody = []byte("Method Not Allowed")
	forbiddenResponseBody        = []byte("403 Forbidden")
)

// Change default MethodNotAllowed tips
func SetMethodNotAllowedBody(body []byte) {
	methodNotAllowedResponseBody = body
}

// Change default Forbidden tips
func SetForbiddenResponseBody(body []byte) {
	forbiddenResponseBody = body
}

func MethodNotAllowedHandleFunc(ctx *gin.Context) {
	ctx.String(http.StatusMethodNotAllowed, string(methodNotAllowedResponseBody))
}

func ForbiddenHandleFunc(ctx *gin.Context) {
	ctx.String(http.StatusForbidden, string(forbiddenResponseBody))
}

// Abstract RestViewer
type RestViewer interface {
	// Before call your handle
	PrePare(ctx *gin.Context)

	// After call your handle
	Finish(ctx *gin.Context)

	// Continue precess if return true
	// ForbiddenHandleFunc will be called if return false
	IsAllowRequest(ctx *gin.Context) bool

	// Implement this method to handle panic error
	HandleError(ctx *gin.Context, err interface{})
}

// Base RestViewer implement
// All RestView must extend this struct
type RestView struct{}

// Implement RestViewer
// Overwrite to your own idea
func (b *RestView) PrePare(ctx *gin.Context) {}

// Implement RestViewer
// Overwrite to your own idea
func (b *RestView) Finish(ctx *gin.Context) {}

// Implement RestViewer
// Default always return true
// Overwrite to your own idea
func (b *RestView) IsAllowRequest(ctx *gin.Context) bool {
	return true
}

// Implement RestViewer
// Default panic
// Overwrite to your own idea
func (b *RestView) HandleError(ctx *gin.Context, err interface{}) {
	panic(err)
}



type TemplateViewer interface {
	RestViewer
	GetTemplateName(ctx *gin.Context) string
	GetTemplateContext(ctx *gin.Context) interface{}
}

// Template view
type TemplateView struct {
	RestView
}

// Overwrite this method
// Return your template name
func (t *TemplateView) GetTemplateName(ctx *gin.Context) string {
	panic("implement me")
}

// Set template Context
func (t *TemplateView) GetTemplateContext(ctx *gin.Context) interface{} {
	return nil
}

func (t *TemplateView) Get(ctx *gin.Context) {
	panic("implement me")
}
