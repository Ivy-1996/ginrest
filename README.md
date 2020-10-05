## ginrest


## Getting started
```shell
go get github.com/Ivy-1996/ginrest
```

## A Simple Example

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

func main() {
	server := gin.Default()
	server.Any("/", ginrest.AsDefaultHandleFunc(new(View)))
	server.Run()
}

```

