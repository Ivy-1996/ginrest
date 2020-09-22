package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	methodNotAllowedResponseBody = []byte("Method Not Allowed")
	forbiddenResponseBody        = []byte("403 Forbidden")
)

func SetMethodNotAllowedBody(body []byte) {
	methodNotAllowedResponseBody = body
}

func SetForbiddenResponseBody(body []byte) {
	forbiddenResponseBody = body
}

func MethodNotAllowedHandleFunc(ctx *gin.Context) {
	ctx.String(http.StatusMethodNotAllowed, string(methodNotAllowedResponseBody))
}

func ForbiddenHandleFunc(ctx *gin.Context) {
	ctx.String(http.StatusForbidden, string(forbiddenResponseBody))
}

type RestViewer interface {
	// before call your handle
	PrePare(ctx *gin.Context)

	// after call your handle
	Finish(ctx *gin.Context)

	Get(ctx *gin.Context)
	Post(ctx *gin.Context)
	Put(ctx *gin.Context)
	Patch(ctx *gin.Context)
	Delete(ctx *gin.Context)

	IsAllowRequest(ctx *gin.Context) bool

	HandleError(ctx *gin.Context, err interface{})
}

type RestView struct{}

func (b *RestView) PrePare(ctx *gin.Context) {}

func (b *RestView) Finish(ctx *gin.Context) {}

func (b *RestView) IsAllowRequest(ctx *gin.Context) bool {
	return true
}

func (b *RestView) Get(ctx *gin.Context) {
	MethodNotAllowedHandleFunc(ctx)
}
func (b *RestView) Post(ctx *gin.Context) {
	MethodNotAllowedHandleFunc(ctx)
}
func (b *RestView) Put(ctx *gin.Context) {
	MethodNotAllowedHandleFunc(ctx)
}
func (b *RestView) Patch(ctx *gin.Context) {
	MethodNotAllowedHandleFunc(ctx)
}
func (b *RestView) Delete(ctx *gin.Context) {
	MethodNotAllowedHandleFunc(ctx)
}

func (b *RestView) HandleError(ctx *gin.Context, err interface{}) {
	panic(err)
}
