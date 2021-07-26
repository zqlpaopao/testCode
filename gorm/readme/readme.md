# 1、gorm的执行信息
 <font color=red size=5x>**获取信息**</font>

- 执行 SQL 的当前时间；
- 执行 SQL 的文件地址和行号；
- 执行 SQL 的花费时长；
- 执行 SQL 的影响行数；
- 执行的 SQL 语句；



 <font color=red size=5x>**支持信息**</font>

- 支持设置 trace_id
- 支持设置 request 信息
- 支持设置 response 信息
- 支持设置 third_party_requests 三方请求信息
- 支持设置 debugs 打印调试信息
- 支持设置 sqls 执行 SQL 信息
- 可记录 cost_seconds 执行时长



编写 `CallBacks` 插件代码，GORM 的 Plugin 接口的编写非常简单，只需要实现两个方法即可。

```go
// Plugin GORM plugin interface
type Plugin interface {
	Name() string
	Initialize(*DB) error
}
```

