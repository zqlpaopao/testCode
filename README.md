# testCode
## 1、浮点数计算
[文章](https://mp.weixin.qq.com/s/7Jd5m1pPfivi727R6TTIpA)

[浮点数计算](github.com/shopspring/decimal)

## 2、go的资源使用情况
[文章](https://mp.weixin.qq.com/s/eh7Wlc___Z4uJz2cVtig1g)

[gopsutil](https://github.com/shirou/gopsutil)
- cpu：系统CPU 相关模块；
- disk：系统磁盘相关模块；
- docker：docker 相关模块；
- mem：内存相关模块；
- net：网络相关；
- process：进程相关模块；
- winservices：Windows 服务相关模块。

## 3、打印内部细节
[文章](https://mp.weixin.qq.com/s/ry8NC3xBYp0FaV1ytCadAw)

[spew](github.com/davecgh/go-spew/spew)

## 4、全文检索库
[文章](https://mp.weixin.qq.com/s/eNaVTI90_lvzMjioN4vmJQ)

[bluge](https://github.com/blugelabs/bluge)

## 5、gosnmp
[gosnmp](https://github.com/gosnmp/gosnmp)


## 6、skywalking监控服务
[文章](https://mp.weixin.qq.com/s/2uSN1SHtQGSoVZkWPFVwRw)
[skywalking](github.com/SkyAPM/go2sky)
- 服务、服务实例、端点指标分析
- 根本原因分析。在运行时分析代码
- 服务拓扑图分析
- 服务、服务实例和端点依赖分析
- 检测到缓慢的服务和端点
- 性能优化
- 分布式跟踪和上下文传播
- 数据库访问指标。检测慢速数据库访问语句（包括 SQL 语句）
- 警报
- 浏览器性能监控
- 基础设施（VM、网络、磁盘等）监控
- 跨指标、跟踪和日志的协作
- 对比
- [zipkin-go-opentracing](https://github.com/openzipkin-contrib/zipkin-go-opentracing)

## 7、gojieba中文分词
[文章](https://mp.weixin.qq.com/s/zRmAjQ0o9n8FE1R0WcnUtQ)

[gojieba](https://github.com/yanyiwu/gojieba)

## 8、rk-boot多服务管理
[rk-boot](https://github.com/rookie-ninja/rk-boot)

## 9、Elasticsearch Api
[elastic](github.com/olivere/elastic/v7)
[文章](https://mp.weixin.qq.com/s/iHIxsEZf3w06GbO2sHSuRA)

## 10、go传输工具croc
[文章](https://mp.weixin.qq.com/s/LliaEz8JY5PS4gDI2Fq4Qw)

## 11、canal配置同步mysql的binlog
[文章](https://mp.weixin.qq.com/s/RS12LsrTvnbNQ3TpCmPh5Q)

[目录canal]()

## 12、规则引擎go
[govaluate](https://mp.weixin.qq.com/s/OH2KT9XiYrzjWnk4q81BGw)

## 13、go-docker-file
[文章](https://mp.weixin.qq.com/s/773INmwebAIy6zDtGHOEoQ)

[目录](go-dockerFile)

## go逃过gc操作内存
[文章](https://github.com/heiyeluren/XMM)

## json只取出某个值
[github](https://github.com/tidwall/gjson)

## golang启动cli参数
[github](https://github.com/urfave/cli/v2)

## 获取json中指定的值，判断是否存在
[github](https://github.com/tidwall/gjson)

## fastjson
[github](https://github.com/valyala/fastjson)

## go包性能分析工具
[pyroscope](https://github.com/pyroscope-io/pyroscope)

## mysql
About
Bifrost ---- 面向生产环境的 MySQL 同步到Redis,MongoDB,ClickHouse,MySQL等服务的异构中间件
[Bifrost](https://github.com/brokercap/Bifrost)

## 钉钉介入
[github](https://github.com/blinkbean/dingtalk)

## 各种glang库
[-](https://github.com/jobbole/awesome-go-cn/)

## 监控tcp连接
[-](https://github.com/kevwan/tproxy)

## gin-dump gin的中间件
可查看 header body等信息
[github](https://github.com/tpkeeper/gin-dump)

[文章](https://studygolang.com/topics/9104?fr=sidebar)

## go操作markdown
[文章](https://mp.weixin.qq.com/s/vDwhyZyF5jrNkTJWZ2MMUg)

[github](https://github.com/yuin/goldmark)

## 分布式链路追踪
[github](https://github.com/jaegertracing/jaeger)

[文章](https://mp.weixin.qq.com/s/3q8KBWDqWCzvW0a5SPA1ug)

[文章1]（cnblogs.com/ExMan/p/12084524.html）

## grpc连接池，提高项目的并发能力
[github](https://github.com/rfyiamcool/grpc-client-pool)

[文章](https://xiaorui.cc/archives/6001)


[高性能net库]
- 传统的net是来一个就启动一个goroutine去处理，如果有一千万就有一千万个goroutine
- gnet采用 初始化启动固定数量的线程来处理

[github](https://github.com/panjf2000/gnet/)

[文章](https://mp.weixin.qq.com/s/aBdvYvoIO2FTMTPDY_IFYQ)

[文章](https://blog.csdn.net/qq_31967569/article/details/103107707)

## 聚合消息推送系统 钉钉 企业微信 email等
[github](https://github.com/rxrddd/austin-go)

## docker监控工具ctop
[github](https://github.com/bcicen/ctop)

[文章](https://mp.weixin.qq.com/s/hMNreJdR1yeUx6GmH-dSsQ)

## 打桩单元测试覆盖率
[文章](https://mp.weixin.qq.com/s/3HRHjDUSnExdb6njMtaXXw)

## gocover
[文章](https://www.jianshu.com/p/e3b2b1194830)

[github](https://github.com/smartystreets/goconvey)


## goreman类似supervisor的进程管理工具
[github](https://github.com/mattn/goreman)

[搭建etcd]（https://cloud.tencent.com/developer/article/1824475）

[文章](https://blog.csdn.net/zb199738/article/details/124769389)


## k8s-docker 日志收集查看
Grafana日志聚合工具Loki搭建使用

[文章](https://dylanyang.top/post/2020/03/21/grafana%E6%97%A5%E5%BF%97%E8%81%9A%E5%90%88%E5%B7%A5%E5%85%B7loki%E6%90%AD%E5%BB%BA%E4%BD%BF%E7%94%A8/)

## 媲美es的搜索引擎-gofund
[github](https://github.com/newpanjing/gofound)
[文章](https://mp.weixin.qq.com/s/ZRJeiD57KrP9zWEOChmeMQ)

## 指定不同场景不同的json字段
[json-filter](https://gocn.vip/topics/3Qz4jeTKoe)


## 统一身份认证系统-casdoor
[github](https://github.com/casdoor/casdoor)

## 绘制图表-gochart
[github](github.com/go-echarts/go-echarts/v2/charts)

## 使用 tableflip 实现应用的优雅热升级
[文章]([使用 tableflip 实现应用的优雅热升级](https://gocn.vip/topics/qoYMdnTxo4))

[]()

[文章](https://gocn.vip/topics/yQ0meNiewd)

## 跨云vpc项目
[github](https://github.com/ICKelin/cframe)

## 高效率json库
[github](github.com/json-iterator/go)

## 压测工具
[文章](https://mp.weixin.qq.com/s/nLdG-LHjbhZ8iyF4-PHeRg)

## 定位goroutine泄漏
[文章]（https://mp.weixin.qq.com/s/LyOD2EoG_M5JFQGUyELVkQ）

## 文件操作
[文章](https://mp.weixin.qq.com/s/Rz8ddGV6pq8-z3swiEDpVQ)

## qq群消息
[文章](https://mp.weixin.qq.com/s/klGtFmDc65WpaXSMTNE_Yg)

## 接口压测
[文章](https://github.com/adjust/go-wrk)

## 代码依赖关系
[depth](https://github.com/KyleBanks/depth)

## go-swagger
[go-swagger](https://github.com/go-swagger/go-swagger)

## 查看代码的时间复杂度
[gocyclo](https://github.com/fzipp/gocyclo)

## 调用链路可视化工具
[文章](https://mp.weixin.qq.com/s/1UhSS9mJBf_C6j9Mpx14uA)

[github](https://github.com/ofabry/go-callvis)
