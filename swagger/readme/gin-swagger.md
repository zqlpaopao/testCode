# 1、安装swagger

```go
go get -u github.com/swaggo/swag/cmd/swag
```

验证

```go
swag -v                                                            
swag version v1.7.0
```



# 2、 安装gin-swagger

```go
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/gin-swagger/swaggerFiles
```



# 3、编写swagger注释

```go
// @Tags Base
// @Summary 用户登录
// @Produce  application/json
// @Param data body request.LoginRequest true "用户登录接口"
// @Success 200 {string} string "{"success":true,"data": { "user": { "username": "asong", "nickname": "", "avatar": "" }, "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VybmFtZSI6ImFzb25nIiwiZXhwIjoxNTk2OTAyMzEyLCJpc3MiOiJhc29uZyIsIm5iZiI6MTU5Njg5NDExMn0.uUS1TreZusX-hL3nKOSNYZIeZ_0BGrxWjKI6xdpdO40", "expiresAt": 1596902312000 },,"msg":"操作成功"}"
// @Router /base/login [post]
func login(c *gin.Context)  {}

// @Tags User
// @Summary 用户修改密码
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body request.ChangePassword true "用户修改密码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/setPassword [PUT]
func setPassword(c *gin.Context) {}
```

# 4、生成文档

```go
 swag init -o /Users/zhangsan/Documents/GitHub/testCode/swagger/docs
```

-o 指定文档生成位置

```go
swagger/docs
├── docs.go
├── swagger.json
└── swagger.yaml
```

在路由注册的地方引入这个docs

```go
/**
 * @Author: zhangsan
 * @Description:
 * @File:  main
 * @Version: 1.0.0
 * @Date: 2021/3/17 上午10:11
 */

package  main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"test/swagger/src"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"test/swagger/docs"//需要引入
)

func main() {
	r := gin.Default()
	r.GET("/login", src.Login)
	r.GET("/set-password", src.SetPassword)
	docs.SwaggerInfo.Title = "go-siteExample API"
	docs.SwaggerInfo.Description = "This is a sample server go-site server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", "localhost", 8080)
	docs.SwaggerInfo.BasePath = ""
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run() // listen and serve on 0.0.0.0:8080
}
```

设置版本信息等

# 5、 访问

```go
http://127.0.0.1:8080/docs/index.html
```

![image-20210317113647157](gin-swagger.assets/image-20210317113647157.png)