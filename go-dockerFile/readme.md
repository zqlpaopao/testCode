有一些方法容器化 `Golang` 工程，尤其是当您使用 `Docker` 运行 `Go` 项目的可执行文件时。我们可以从我们的项目中创建我们的镜像，简单地在您的本地计算机上运行它，甚至可以通过从 `harbour` 中提取您的镜像运行它。

##  要求

- docker
- Go 编程环境
- 仓库: https://github.com/strictnerd/go-dockerfile

##  入门

首先，您需要使用`systemctl start docker`或启动您的 `docker` 守护进程`service docker start`，如果需要必要权限请使用`sudo`。

然后我们将创建我们简单的 `go HTTP` 代码。

```
$ mkdir go-dockerfile && cd go-dockerfile 
$ go mod init myapp 
$ touch server.go
```

server.go：

```
package main

import (
 "os"

 "github.com/gin-gonic/gin"
 "github.com/joho/godotenv"
)

func init() {
 godotenv.Load()
}

func main() {
 port := os.Getenv("PORT")
 if port == "" {
  port = "8080"
 }

 router := gin.Default()

 router.GET("/", func(c *gin.Context) {
  c.String(200, "Hello World")
 })

 router.GET("/env", func(c *gin.Context) {
  c.String(200, "Hello %s", os.Getenv("NAME"))
 })

 router.Run(":" + port)
}
```

我们将包含一个简单的`gin`路由器和可选的`.server.gogodotenv``/path` 将返回`“Hello World”，/envpath` 将返回`“Hello ${NAME}”`。

##  Dockerfile

有多种编写方式`Dockerfile`，但我将使用不同的基础图像制作

3 个示例：`golang docker、alpine、scratch`。

###  官方Dockerfile

```
FROM golang:1.16-alpine

WORKDIR /project/go-docker/

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .
RUN go build -o /project/go-docker/build/myapp .

EXPOSE 8080
ENTRYPOINT [ "/project/go-docker/build/myapp" ]
```

在此Dockerfile，我们将其分为几个部分：

- FROM golang:1.16-alpine，我们将golang:1.16-alpine用作此 Docker 构建的基础镜像。
- WORKDIR, 将是我们命令的工作目录/下一个命令的路径。
- COPY go.* ./，我们将从我们的项目复制go.mod&go.sum文件到工作目录。
- RUN go mod download , 从 go modules 下载项目依赖。
- COPY . . ，将我们项目中的所有内容复制到工作目录中。
- RUN go build -o /project/go-docker/build/myapp ., 在工作目录中构建我们的项目并将其project/go-docker/build/myapp作为二进制文件输出。
- EXPOSE 8080，告诉 docker 我们的代码将暴露端口8080。
- ENTRYPOINT ["/project/go-docker/build/myapp"] ，当我们运行这个镜像的容器时，它将从我们的构建可执行文件开始执行。

之后我们需要运行这个命令：

```
docker build -f Dockerfile -t test-go-docker:latest .
```

- -f flag 是我们的Dockerfile.
- -t tag 是镜像标签。
- .命令是当前文件夹下的Dockerfile.

尝试运行此命令`docker images`，例如：

![图片](https://mmbiz.qpic.cn/mmbiz_png/hvZjCFh6diaRDs53Rs3qGl0OBx7E2Oj186EHic4ibVvspWic6pnlgNHUvLPnFUKqkeWBXBoLHxyCMXicvmsSOHf45VA/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

450M

###  alpine

```
FROM golang:1.16-alpine as builder

WORKDIR /project/go-docker/

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .
RUN go build -o /project/go-docker/build/myapp .

FROM alpine:latest
COPY --from=builder /project/go-docker/build/myapp /project/go-docker/build/myapp

EXPOSE 8080
ENTRYPOINT [ "/project/go-docker/build/myapp" ]
```

与golang 官方基础镜像的区别：

- FROM golang:1.16-alpine as builder，我们将使用golang:1.16-alpine并标记它，builder因为稍后将使用它。
- FROM alpine:latest，我们用alpine作为基础镜像.
- COPY --from=builder /project/go-docker/build/myapp /project/go-docker/build/myapp ，将构建的二进制文件复制到新的 alpine 镜像中。

这个`Dockerfile`生成的镜像比之前小得多。

![图片](https://mmbiz.qpic.cn/mmbiz_png/hvZjCFh6diaRDs53Rs3qGl0OBx7E2Oj18pvr5W3v5YBMMKtlpUlgruOXtJKSicP5TR8ZpA10raGSNxsacfhSRf3A/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

15.1M

###  Scratch

```
FROM golang:1.16-alpine as builder

WORKDIR /project/go-docker/

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .
RUN go build -o /project/go-docker/build/myapp .

FROM scratch
COPY --from=builder /project/go-docker/build/myapp /project/go-docker/build/myapp

EXPOSE 8080
ENTRYPOINT [ "/project/go-docker/build/myapp" ]
```

对于最后一个 `Dockerfile`，我们只将`alpine`基础镜像更改为`scratch`. `Scratch` 是一个空镜像，所以一旦容器运行，我们就无法执行到容器中，因为它没有 `shell` 命令。如下是输出的 `docker images`。

![图片](https://mmbiz.qpic.cn/mmbiz_png/hvZjCFh6diaRDs53Rs3qGl0OBx7E2Oj18dBFAmc54SsPiaGW0icuqOw2HEGQriaiaQJSP2c8w5gGicPygQYUFgj7IALg/640?wx_fmt=png&tp=webp&wxfrom=5&wx_lazy=1&wx_co=1)

9.52M

运行镜像`docker run -d -p 8080:8080 test-go-docker:latest`，它将端口8080从容器转发到物理节点8080端口并可以访问`http://localhost:8080`.

##  结论

就我个人而言，我会选择第二个`Dockerfile`。为什么？因为体积小而且它还有几个命令和一个`shell`命令所以我们可以`docker exec`进入容器并访问它。如果我们使用`scratch`基础镜像，因为我们无法执行它，所以将很难调试正在运行的容器。