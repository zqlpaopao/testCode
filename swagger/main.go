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
	r.GET("/user/setPassword", src.SetPassword)
	docs.SwaggerInfo.Title = "go-siteExample API"
	docs.SwaggerInfo.Description = "This is a sample server go-site server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", "localhost", 8080)
	docs.SwaggerInfo.BasePath = ""
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run() // listen and serve on 0.0.0.0:8080
}
