package view

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var methodNotAllowedBody = map[string]string{"detail": "Method Not Allowed"}

var permissionDeniedBody = map[string]string{"detail": "Permission Denied"}

var notAuthenticatedBody = map[string]string{"detail": "Incorrect authentication credentials"}

var pageNotFoundBody = map[string]string{"detail": "404 Not Found"}

var throttledBody = map[string]string{"detail": "Request was throttled"}

var MethodNotAllowedHandle = methodNotAllowedHandle

var PermissionDeniedHandle = permissionDeniedHandle

var NotAuthenticatedHandle = notAuthenticatedHandle

var ThrottledHandle = throttledHandle

var PageNotFoundHandle = pageNotFoundHandle

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

func pageNotFoundHandle(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, pageNotFoundBody)
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

func SerPageNotFoundHandle(handle gin.HandlerFunc) {
	PageNotFoundHandle = handle
}
