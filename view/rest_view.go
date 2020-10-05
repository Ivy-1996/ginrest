package view

import (
	"github.com/gin-gonic/gin"
	"reflect"
)

type RestViewHandlerFunc func(reflect.Value, RestViewer, *gin.Context)

// SimpleRestViewFunc
// Just call the handler
func SimpleRestViewFunc(handler reflect.Value, r RestViewer, ctx *gin.Context) {
	context := reflect.ValueOf(ctx)
	handler.Call([]reflect.Value{context})
}

// DefaultRestViewFunc
// Add check request is permitted
func DefaultRestViewFunc(handler reflect.Value, r RestViewer, ctx *gin.Context) {

	// Ensure that the incoming request is permitted
	r.PerformAuthentication(ctx)
	r.CheckPermissions(ctx)
	r.CheckThrottles(ctx)

	// Call handler
	SimpleRestViewFunc(handler, r, ctx)
}
