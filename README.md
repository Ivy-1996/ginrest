# ginrest



## Example

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest"
	"github.com/ivy-1996/ginrest/view"
)

type View struct {
	view.RestView
}

func (*View) Get(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"method": ctx.Request.Method})
}

func (*View) Post(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"method": ctx.Request.Method})
}

type AuthView struct {
	view.RestView
}

const CurrentUserId = "currentUserId"

func (*AuthView) PerformAuthentication(ctx *gin.Context) {
	if userId := ctx.Query("user"); userId != "" {
		ctx.Set(CurrentUserId, userId)
	}
}

func (*AuthView) Get(ctx *gin.Context) {
	if currentUserId, exist := ctx.Get(CurrentUserId); exist {
		ctx.JSON(200, gin.H{"currentUserId": currentUserId})
	} else {
		view.NotAuthenticatedHandle(ctx)
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()
	server.Any("/", ginrest.AsDefaultHandleFunc(new(View)))
	server.Any("/auth", ginrest.AsDefaultHandleFunc(new(AuthView)))
	server.Run(":8002")
}

```