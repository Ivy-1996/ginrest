# ginrest

简单、实用、可扩展gin的第三方扩展



## 路由扩展

在gin中无侵入式使用

```go
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

`AsDefaultHandleFunc`：背后帮我们做了很多事情，但前提是我们得将我们结构体"继承"`view.RestView`, 你可以通过重写某些方法来完成请求的可控化，摆脱去写多层嵌套中间件的烦恼。

`AsHandlerFunc`：`AsDefaultHandleFunc`背后调用了`AsHandlerFunc`, 你可以通过调用这个方法来完成我们自己请求的一个流程控制。



### Router

```go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ivy-1996/ginrest/router"
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

	route := router.NewRouter(server) 

	route.Register("/user", new(User)) 

	server.Run()
}
```

`Router`内部调用了`AsHandleFunc`，使得我们的结构体可以被当做一个`gin.HandlerFunc`来使用

`SetRestViewHandlerFunc`：就像`AsHandleFunc`可以定制化流程一样，我们可以通过调用生成了的`Router`"对象"的`SetRestViewHandlerFunc`方法来设置我们控制流程的函数

`SetRelativePath`：设置当前`router`路由的前缀

`SetAutoPrefix`：设置为`true`会自动格式化当前的路由，默认为`false`。如传入`user`会自动格式化为`/user`

`NewRouter`：`Router`结构体的构造函数，第一个参数为`gin.Engine`对象，后面的参数为`go`的不定长参数，传入的参数会被当做当前`router`的共有的中间件。

`Register`：注册路由，第一个参数为`url`路径，第二个参数为继承了`view.RestView`的结构体，后面的参数为`go`的不定长参数，传入的参数会被当做当前注册的结构体的私有的中间件。



## 数据校验

`gin`自带的数据校验很强大，但返回的错误并不是很好处理。



### Validator

```go
package main

import (
	"fmt"
	"github.com/ivy-1996/ginrest/validator"
)

type A struct {
	Name string `validate:"gt:10:必须大于3"`
	Age  int    `validate:"gt:10:必须大于10"`
}

func (a A) ValidateAge() *validator.ValidateErrorNodes {
	if a.Age < 10 {
		validateErrorNode := validator.NewValidateErrorNode("no_field", "Age必须大于10")
		return validator.NewValidateErrorNodes(validateErrorNode)
	} else if a.Age > 20 {
		validateErrorNode := validator.NewValidateErrorNode("no_field", "Age必须小于20")
		return validator.NewValidateErrorNodes(validateErrorNode)
	}
	return nil
}

func main() {
	validate := validator.NewValidator()

	err := validate.RunValidators(A{Name: "6", Age: 9})
	if err != nil {
		fmt.Println(err.Error())
	}
}

/*
OUT:
	[{"field_name":"Age","validate_error_nodes":[{"code":"no_field","error_message":"Age必须大于10"}]},{"field_name":"Name","validate_error_nodes":[{"code":"gt","error_mes
sage":"必须大于3"}]}]

*/

```

上面的这个简单的例子展示了我们得验证器的两种验证方式，验证器处理和自己处理

* 自己处理：

  * 想要自己处理只需要满足一个条件即可：在当前校验的结构体上写上符合条件的方法即可。详见下文。

    

* 验证器处理：

  * 想要验证器处理，首先必须满足三个条件

    

    * `tag`: 必须在结构体的`Field`里面写上规定的`tag`，默认是`validate`,可以通过创建的验证器对象(如上为：`validator.NewValidator()`)的`SetValidateTag`方法来修改默认的`tag`

    

    * 语法：`tag`的语法必须满足 `ruleName:expect:errmsg`的格式

      * `rule`：自动去`ValidateLibrary`里面匹配的key
      * `expect`: 期待的比较的值
      * `errmsg`：校验失败返回的错误信息

      

    * 自己没有实现校验方法

      * 如何自己实现的校验方法：在当前校验的`struct`写上默认验证的前缀加自己校验字段的名字，默认的前缀为`Validate`, 可以通过创建的校验器对象的`SetValidateFuncPrefix`方法来改变这一前缀。如上`ValidateAge`。返回值必须为`validator.ValidateErrorNodes`的指针。可以通过`validator.NewValidateErrorNodes`来返回多个错误对象。



`validator.RegisterLibrary`方法可以为默认的`ValidateLibrary`注册你自己的处理逻辑,但它不是`goroutine`安全的，最好在服务启动前就完成注册。

`SetValidateLibrary`：可以通过创建的校验器的`SetValidateLibrary`方法来为你的校验器设置你自己的`ValidateLibrary`

`ValidateLibrary`: 接口类型，只需要实现三个方法即可创建你自己的`ValidateLibrary`

```go
// ValidateLibrary interface
type ValidateLibrary interface {
	Register(key string, validateFunc ValidateFunc)
	Validate(field reflect.Value, node *validatorNode) ValidateErrorNodes
	LookForValidateFunc(key string) ValidateFunc
}
```



返回的`err`是一个实现了`error`类型的`type`，它的`Error`方法默认返回当前`type``json`格式的字符串。它可以很方便反序列化为`go`的数据类型供我们自己解析

参考实例

```go
if err, ok := err.(validator.ValidateErrors); ok {
		var result validator.ValidateErrors
		json.Unmarshal([]byte(err.Error()), &result)
	}
```



内部更多的小细节可以通过查看源码发掘！



## TODO  加入更多的支持 !























	









#### 