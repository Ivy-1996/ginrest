package view

import (
	"github.com/gin-gonic/gin"
)

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

	// Perform authentication on the incoming request.
	PerformAuthentication(ctx *gin.Context)

	// Check if the request should be permitted.
	CheckPermissions(ctx *gin.Context)

	// Check if request should be throttled.
	CheckThrottles(ctx *gin.Context)
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

// Implement RestViewer
// Overwrite to your own idea
func (b *RestView) PerformAuthentication(ctx *gin.Context) {}

// Implement RestViewer
// Overwrite to your own idea
func (b *RestView) CheckPermissions(ctx *gin.Context) {}

// Implement RestViewer
// Overwrite to your own idea
func (b *RestView) CheckThrottles(ctx *gin.Context) {}
