# ginrest

> ginrest is simple, helpful , Customizable, realize your idea with it!

[Document](./docs/docs.md)

## Limitation
```shell script
golang version >= 1.11
```

## Installation

####  `go get`

```shell
go get github.com/Ivy-1996/ginrest
```

#### `go mod`

```shell 
require github.com/Ivy-1996/ginrest latest
```



## Example

#### Use with `gin`

```go
package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getUserList(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"result": "userList"})
}

func createUser(context *gin.Context) {
	context.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func main() {

	server := gin.Default()

	server.GET("/user/", getUserList)

	server.POST("/user/", createUser)

	server.Run()
}

```



#### Use with `ginrest`

```GO
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest"
	"github.com/ivy-1996/ginrest/view"
	"net/http"
)

type User struct {
	view.RestView
}

func (*User) Get(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"result": "userList"})
}

func (*User) Post(context *gin.Context) {
	context.JSON(http.StatusCreated, gin.H{"status": "ok"})
}

func main() {
	server := gin.Default()
	server.Any("/user/", ginrest.AsDefaultHandleFunc(new(User)))
	server.Run()
}
```




