package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var methodNotAllowedBody = map[string]string{"detail": "Method Not Allowed"}

var permissionDeniedBody = map[string]string{"detail": "Permission Denied"}

var notAuthenticatedBody = map[string]string{"detail": "Incorrect authentication credentials"}

var throttledBody = map[string]string{"detail": "Request was throttled"}

var MethodNotAllowedHandle = methodNotAllowedHandle

var PermissionDeniedHandle = permissionDeniedHandle

var NotAuthenticatedHandle = notAuthenticatedHandle

var ThrottledHandle = throttledHandle

func methodNotAllowedHandle(ctx *gin.Context) {
	ctx.JSON(http.StatusMethodNotAllowed, methodNotAllowedBody)
}

func permissionDeniedHandle(ctx *gin.Context) {
	ctx.JSON(http.StatusForbidden, permissionDeniedBody)
}

func notAuthenticatedHandle(ctx *gin.Context) {
	ctx.JSON(http.StatusUnauthorized, notAuthenticatedBody)
}

func throttledHandle(ctx *gin.Context) {
	ctx.JSON(http.StatusTooManyRequests, throttledBody)
}

func SetMethodNotAllowedHandle(handle gin.HandlerFunc) {
	MethodNotAllowedHandle = handle
}

func SetPermissionDeniedHandle(handle gin.HandlerFunc) {
	PermissionDeniedHandle = handle
}

func SetNotAuthenticatedHandle(handle gin.HandlerFunc) {
	NotAuthenticatedHandle = handle
}

func SetThrottledHandle(handle gin.HandlerFunc) {
	ThrottledHandle = handle
}
