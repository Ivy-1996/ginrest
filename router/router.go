package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/core"
	"github.com/ivy-1996/ginrest/view"
)

type Router interface {
	Register(view.RestViewer)
}

type DefaultRouter struct {
	engine              *gin.Engine
	restViewHandlerFunc view.RestViewHandlerFunc
}

func NewDefaultRouter(engine *gin.Engine) *DefaultRouter {
	return &DefaultRouter{engine: engine}
}

// Register RestViewer to gin.Engine
func (d *DefaultRouter) Register(path string, v view.RestViewer) {

	group := d.engine.Group(path)

	if d.restViewHandlerFunc == nil {
		d.restViewHandlerFunc = view.DefaultRestViewFunc
	}

	group.Any(path, core.AsHandlerFunc(v, d.restViewHandlerFunc))
}

// Set RestViewHandlerFunc 
func (d *DefaultRouter) SetRestViewHandlerFunc(f view.RestViewHandlerFunc) {
	d.restViewHandlerFunc = f
}
