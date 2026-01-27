package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// https://github.com/go-dev-frame/sponge/blob/main/assets/readme-cn.md

func StaticFS(relativePath, root string) gin.HandlerFunc {
	fileSystem := http.Dir(root)
	fileServer := http.FileServer(fileSystem)

	if relativePath != "/" {
		fileServer = http.StripPrefix(relativePath, fileServer)
	}

	return func(c *gin.Context) {
		// 检查文件是否存在
		path := filepath.Join(root, c.Request.URL.Path)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// 文件不存在，返回 index.html（用于 SPA）
			c.File(filepath.Join(root, "index.html"))
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	r := gin.Default()

	// 使用自定义的静态文件服务
	r.GET("/user/*filepath", StaticFS("/user/", "./static"))

	log.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
