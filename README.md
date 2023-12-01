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

## 12、mysql binglog监听同步工具

[github](https://github.com/go-mysql-org/go-mysql@v1.7.0/canal/config.go)
[文章](https://mp.weixin.qq.com/s/DHMmusc554o1TVvBr3A94g)

## 13、mysql 监听binglog 同步es

[github](https://github.com/go-mysql-org/go-mysql-elasticsearch)
[文章](https://mp.weixin.qq.com/s/14Jpb3uBJz1yCd8gKuTVQw)

## go-mysql 数据同步 binglog
[github](https://github.com/go-mysql-org/go-mysql)
[文章](https://mp.weixin.qq.com/s/VPYYQc8c4rWTVS6Paw0vnw)
- 主从复制（Replication）
- 增量同步（Incremental dumping）
- 客户端（Client）
- 虚拟服务端（Fake server）
- 高可用（Failover）
- mysql的驱动（database/sql like driver）

## 14、mysql

About
Bifrost ---- 面向生产环境的 MySQL 同步到Redis,MongoDB,ClickHouse,MySQL等服务的异构中间件
[Bifrost](https://github.com/brokercap/Bifrost)

## 15、规则引擎go

[govaluate](https://mp.weixin.qq.com/s/OH2KT9XiYrzjWnk4q81BGw)

## 16、go-docker-file

[文章](https://mp.weixin.qq.com/s/773INmwebAIy6zDtGHOEoQ)

[目录](go-dockerFile)

## 17、go逃过gc操作内存

[文章](https://github.com/heiyeluren/XMM)



## 18、golang启动cli参数

[github](https://github.com/urfave/cli/v2)

## 19、获取json中指定的值，判断是否存在

[github](https://github.com/tidwall/gjson)

## 20、fastjson

[github](https://github.com/valyala/fastjson)

## 21、go包性能分析工具

[pyroscope](https://github.com/pyroscope-io/pyroscope)

## 22、json只取出某个值

[github](https://github.com/tidwall/gjson)

## 23、钉钉介入

[github](https://github.com/blinkbean/dingtalk)

## 24、各种glang库

[-](https://github.com/jobbole/awesome-go-cn/)

## 25、监控tcp连接

[-](https://github.com/kevwan/tproxy)

## 26、gin-dump gin的中间件

可查看 header body等信息
[github](https://github.com/tpkeeper/gin-dump)

[文章](https://studygolang.com/topics/9104?fr=sidebar)

## 27、go操作markdown

[文章](https://mp.weixin.qq.com/s/vDwhyZyF5jrNkTJWZ2MMUg)

[github](https://github.com/yuin/goldmark)

## 28、分布式链路追踪

[github](https://github.com/jaegertracing/jaeger)

[文章](https://mp.weixin.qq.com/s/3q8KBWDqWCzvW0a5SPA1ug)

[文章1]（cnblogs.com/ExMan/p/12084524.html）

## 29、grpc连接池，提高项目的并发能力

[github](https://github.com/rfyiamcool/grpc-client-pool)

[文章](https://xiaorui.cc/archives/6001)


[高性能net库]

- 传统的net是来一个就启动一个goroutine去处理，如果有一千万就有一千万个goroutine
- gnet采用 初始化启动固定数量的线程来处理

[github](https://github.com/panjf2000/gnet/)

[文章](https://mp.weixin.qq.com/s/aBdvYvoIO2FTMTPDY_IFYQ)

[文章](https://blog.csdn.net/qq_31967569/article/details/103107707)

## 30、聚合消息推送系统 钉钉 企业微信 email等

[github](https://github.com/rxrddd/austin-go)

## 31、docker监控工具ctop

[github](https://github.com/bcicen/ctop)

[文章](https://mp.weixin.qq.com/s/hMNreJdR1yeUx6GmH-dSsQ)

## 32、打桩单元测试覆盖率

[文章](https://mp.weixin.qq.com/s/3HRHjDUSnExdb6njMtaXXw)

## 33、gocover

[文章](https://www.jianshu.com/p/e3b2b1194830)

[github](https://github.com/smartystreets/goconvey)


## 34、goreman类似supervisor的进程管理工具

[github](https://github.com/mattn/goreman)

[搭建etcd]（https://cloud.tencent.com/developer/article/1824475）

[文章](https://blog.csdn.net/zb199738/article/details/124769389)


## 35、k8s-docker 日志收集查看

Grafana日志聚合工具Loki搭建使用

[文章](https://dylanyang.top/post/2020/03/21/grafana%E6%97%A5%E5%BF%97%E8%81%9A%E5%90%88%E5%B7%A5%E5%85%B7loki%E6%90%AD%E5%BB%BA%E4%BD%BF%E7%94%A8/)

## 36、媲美es的搜索引擎-gofund

[github](https://github.com/newpanjing/gofound)
[文章](https://mp.weixin.qq.com/s/ZRJeiD57KrP9zWEOChmeMQ)

## 37、指定不同场景不同的json字段

[json-filter](https://gocn.vip/topics/3Qz4jeTKoe)


## 38、统一身份认证系统-casdoor

[github](https://github.com/casdoor/casdoor)

## 39、绘制图表-gochart

[github](github.com/go-echarts/go-echarts/v2/charts)

## 40、使用 tableflip 实现应用的优雅热升级

[文章]([使用 tableflip 实现应用的优雅热升级](https://gocn.vip/topics/qoYMdnTxo4))

[]()

[文章](https://gocn.vip/topics/yQ0meNiewd)

## 41、跨云vpc项目

[github](https://github.com/ICKelin/cframe)

## 42、高效率json库

[github](github.com/json-iterator/go)

## 43、压测工具

[文章](https://mp.weixin.qq.com/s/nLdG-LHjbhZ8iyF4-PHeRg)

## 44、定位goroutine泄漏

[文章]（https://mp.weixin.qq.com/s/LyOD2EoG_M5JFQGUyELVkQ）

## 45、文件操作

[文章](https://mp.weixin.qq.com/s/Rz8ddGV6pq8-z3swiEDpVQ)

## 46、qq群消息

[文章](https://mp.weixin.qq.com/s/klGtFmDc65WpaXSMTNE_Yg)

## 47、接口压测

[文章](https://github.com/adjust/go-wrk)

## 48、代码依赖关系

[depth](https://github.com/KyleBanks/depth)

## 49、go-swagger

[go-swagger](https://github.com/go-swagger/go-swagger)

## 50、查看代码的时间复杂度

[gocyclo](https://github.com/fzipp/gocyclo)

## 51、调用链路可视化工具

[文章](https://mp.weixin.qq.com/s/1UhSS9mJBf_C6j9Mpx14uA)

[github](https://github.com/ofabry/go-callvis)

## 52、限流中间件

[文章]（https://www.iwesen.com/post/120.html）

[github](https://github.com/didip/tollbooth)

## 53、分布式任务调度

[github](https://github.com/iwannay/jiacrontab)

[github](https://github.com/busgo/pink)

## 54、限流、address、cookie、method

[github]（https://github.com/throttled/throttled/）

[文章](https://mp.weixin.qq.com/s/PM4ocsopqlTezIheI4qZzA)

## 55、操作excel

[文](https://mp.weixin.qq.com/s/K3sC_0msgmXhNtXkMnXCIA)

## 56、单机和集群限流工具

[github](https://github.com/throttled/throttled)

## 57、自动检测及调整内存对齐 fieldalignment

[文章](https://mp.weixin.qq.com/s/1LeetxzFAqj3RCvnz_adZQ)

[github](golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment)

## 58、BFE七层负载均衡-流量转发控制

[文章](https://www.topgoer.cn/docs/bfe/bfe-1cva4i69k88se)

## 59、可视化监控工runtime 

[github](https://github.com/arl/statsviz)

## 60、数据库可视化，流程管理工单系统bytebase

[文章](https://mp.weixin.qq.com/s/lz68jC63yMrjAQT4fVBabw)
[github](https://github.com/bytebase/bytebase)

## 61、测试覆盖率godoc

[github](https://github.com/qiniu/goc)

## 62、网络质量监控

[github](https://github.com/smartping/smartping)


## 63、websocket 

[github](https://github.com/eranyanay/1m-go-websockets)

## 64、go爬虫

[github](https://github.com/gocolly/colly)

## 65、多主机管理

[gitree](https://gitee.com/openspug/spug)
![image](https://user-images.githubusercontent.com/43371021/228761667-eeb2f0cf-bd88-4f32-a7a0-fb7ebb5817b6.png)

## 66、runnergo 性能测试平台

[git](https://github.com/Runner-Go-Team/runnerGo)
https://camo.githubusercontent.com/93563d56ed2cfc21602d8efae2161ab98497100c96017bf7b4f6587a1c24f016/68747470733a2f2f617069706f73742e6f73732d636e2d6265696a696e672e616c6979756e63732e636f6d2f6b756e70656e672f696d616765732f686f6d652e6a7067


## 67、nginx 管理ui

[nginx-proxy-manager](https://mp.weixin.qq.com/s/n981QudKz_pO-7RuoVJ3wA)

## 68、go 网关

[1](https://github.com/eolinker/apinto)
[2](https://github.com/megaease/easegress)
[3](https://github.com/fagongzi/manba)


## 69、proto-doc

[protobuf文档生成](https://github.com/pseudomuto/protoc-gen-doc)

## 70、go版本日志手机系统

[七牛云](https://github.com/qiniu/logkit)

## 71、syslog解析器

[解析器](https://github.com/influxdata/go-syslog)


## 72、超牛逼的程序详情调度展示

[git](https://github.com/bcicen/grmon)



# 73、go 压测ui话
[github](https://github.com/link1st/go-stress-testing)


# 74、网关
[github](https://github.com/go-mysql-org/go-mysql-elasticsearch)

# 75、在线文档
[github](https://github.com/mindoc-org/mindoc)

# 76、部署web
[github](https://github.com/dreamans/syncd)

# 77、负载均衡算法
[github](https://github.com/zehuamama/balancer)


# 问答管理后台
[问答管理后台](https://github.com/meloalright/guora)

# parquest
[parquest](https://github.com/parsyl/parquet/tree/master)


# 设备配置文件
[设备配置文件](https://github.com/udhos/jazigo#about-jazigo)

# 设备的telmetry-gateway
[设备的telmetry-gateway](https://github.com/yahoo/panoptes-stream)

# 日志监控系统
[文章](https://mp.weixin.qq.com/s/9032e4gk2Y8kb2EVJ7eWHA)
[github](https://github.com/AutohomeCorp/frostmourne/tree/master)



可视化大屏
[github](https://github.com/dataease/dataease)
[文章](https://mp.weixin.qq.com/s/PArqTRgh-UaVjWKvuOCFGw)


id生成器
[文章](https://mp.weixin.qq.com/s/tfZ5zHo8FG_Rc1JteLBS7g)

[github](https://github.com/speps/go-hashids)

结构体转为csv
[文章](https://mp.weixin.qq.com/s/9Vg_kGeX744E8cMZ68a2_w)

[github](github.com/gocarina/gocsv)

byte转为mb、kb
[文章](https://mp.weixin.qq.com/s/UGvuxILtk7lfMR_bFSkcYg)

[github](https://github.com/dustin/go-humanize)

# 无锁 map
[文章](https://mp.weixin.qq.com/s/gDt92DXSAc8EOR0IhJS-Ig)
[github](https://github.com/orcaman/concurrent-map)

# csv-struct操作
[git](https://github.com/gocarina/gocsv)
[文章](https://mp.weixin.qq.com/s/9Vg_kGeX744E8cMZ68a2_w)

# 进程管理
[github](https://github.com/mattn/goreman)

# 服务器接口监控 
http、grpc redis等
[github](https://github.com/louislam/uptime-kuma)

# etcd健康诊断
[github](https://github.com/ahrtr/etcd-diagnosis)
[文章](https://mp.weixin.qq.com/s/yAy9tTphn7eQmDIGIrHpYA)

#[表格式输出]
[github](https://github.com/olekukonko/tablewriter)

# 爬虫平台
[github](https://github.com/crawlab-team/crawlab/)

# 高性能的golang web框架
[文章](https://mp.weixin.qq.com/s/CPlq5SiuhHRRzQ1EzElExQ)

# kafka管理工具
[Kafka监控系统-Kafka Eagle](https://mp.weixin.qq.com/s/hbmWGxL2lEAA0fhvJSLSWA)

# cmdb
[-](https://github.com/veops/cmdb)
