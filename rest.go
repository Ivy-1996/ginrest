package ginrest

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/http"
	"github.com/ivy-1996/ginrest/router"
	"github.com/ivy-1996/ginrest/view"
)

// Shortcut for http.AsHandlerFunc
func AsHandlerFunc(r view.RestViewer, restHandlerFunc view.RestViewHandlerFunc) gin.HandlerFunc {
	return http.AsHandlerFunc(r, restHandlerFunc)
}

// Shortcut for view.SimpleRestViewFunc
func AsSimpleHandleFunc(r view.RestViewer) gin.HandlerFunc {
	return AsHandlerFunc(r, view.SimpleRestViewFunc)
}

// Shortcut for view.DefaultRestViewFunc
func AsDefaultHandleFunc(r view.RestViewer) gin.HandlerFunc {
	return AsHandlerFunc(r, view.DefaultRestViewFunc)
}

// Shortcut for router.NewDefaultRouter
func NewRouter(engine *gin.Engine, middleware ...gin.HandlerFunc) *router.Router {
	return router.NewRouter(engine, middleware...)
}

const author = "github.com/ivy-1996"
