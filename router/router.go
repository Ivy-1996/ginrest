package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/http"
	"github.com/ivy-1996/ginrest/view"
	"strings"
)

const pathPrefix = "/"

type Router struct {
	engine              *gin.Engine
	restViewHandlerFunc view.RestViewHandlerFunc
	middleWare          []gin.HandlerFunc
	group               *gin.RouterGroup
	relativePath        string
	autoPrefix          bool
}

// Set AutoPrefix with this method
// Call it before Register otherwise it won't be worked
func (r *Router) SetAutoPrefix(autoPrefix bool) {
	r.autoPrefix = autoPrefix
}

// Constructor of Router
func NewRouter(engine *gin.Engine, middleWare ...gin.HandlerFunc) *Router {
	return &Router{engine: engine, middleWare: middleWare}
}

// initialize for Router
func (r *Router) prepare() {

	if r.group == nil {
		// create a new group
		// relativePath maybe a empty string
		r.group = r.engine.Group(r.relativePath)

		// public middleware for this router
		r.group.Use(r.middleWare...)

	}

	// get defaultRestViewFunc if it is nil
	if r.restViewHandlerFunc == nil {
		r.restViewHandlerFunc = view.DefaultRestViewFunc
	}
}

// Get final path from this function
// Your can call SetAutoPrefix to change the result
func (r *Router) formatPath(path string) string {
	if r.autoPrefix {
		if strings.HasPrefix(path, pathPrefix) {
			return path
		} else {
			return pathPrefix + path
		}
	}
	return path
}

// Register RestViewer to gin.Engine
func (r *Router) Register(path string, v view.RestViewer, middleware ...gin.HandlerFunc) {

	// initialize Router
	r.prepare()

	// private middleware use
	handlers := append(middleware, http.AsHandlerFunc(v, r.restViewHandlerFunc))

	// Get formatPath
	// Call SetAutoPrefix to change value
	path = r.formatPath(path)

	// register done
	r.group.Any(path, handlers...)
}

// Set your RestViewHandlerFunc with this method
// Call this before Register otherwise it won't be worked
func (r *Router) SetRestViewHandlerFunc(f view.RestViewHandlerFunc) {
	r.restViewHandlerFunc = f
}

// Set your relativePath with this method
// Call this before Register otherwise it won't be worked
func (r *Router) SetRelativePath(relativePath string) {
	r.relativePath = relativePath
}
