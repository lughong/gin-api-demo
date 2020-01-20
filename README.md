项目介绍
~~~
Go API(REST + JSON)
项目地址：https://github.com/lughong/gin-api-demo
~~~

目录结构
~~~
├── cmd                        # Go程序唯一入口
├── config                     # 专门用来处理配置和配置文件的Go package
├── docs                       # 文档存放目录
├── global                     # 全局包目录
│   ├── constvar               # 常量统一存放位置
│   ├── errno                  # 错误码存放位置
│   ├── redis                  # redis包
│   ├── token                  # token包
│   └── version                # 版本包
├── handler                    # 类似MVC架构中的C，用来读取输入，并将处理流程转发给实际的处理函数（CLI、web、REST、gRPC)
│   └── http                   # REST user处理示例
├── logic                      # 逻辑实现层
├── mock                       # 模拟库
├── model                      # 数据库相关的操作统一放在这里，包括数据库初始化和对表的增删改查
├── registry                   # 依赖注入容器
├── repository                 # 仓库实现层（NoSQL、RDBMS、Micro-Services）
├── router                     # 路由相关处理
│   └── middleware             # API服务器用的是Gin Web框架，Gin中间件存放位置
├── testdata                   # 测试数据目录
├── util                       # 工具类函数存放目录
├── Makefile                   # Makefile文件
├── README.md                  # README.md文件
├── go.mod                     # 记录依赖包及其版本号
└── go.sum                     
~~~

克隆项目
~~~
$ git clone https://github.com/lughong/gin-api-demo.git
~~~

数据库配置
~~~
mysql> CREATE DATABASE IF NOT EXISTS gin_api_demo DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;
mysql> use gin_api_demo
mysql> source {app_root}/testdata/gin_api_demo.sql
mysql> CREATE user 'apiuser'@'127.0.0.1' IDENTIFIED BY '123456';
mysql> GRANT ALL PRIVILEGES ON gin_api_demo.* TO 'apiuser'@'127.0.0.1';
mysql> FLUSH PRIVILEGES;
~~~

直接运行项目
~~~
$ go run github.com/lughong/gin-api-demo/cmd/...
~~~

运行过程若出现module失败，可以尝试设置GOPROXY环境变量
~~~
$ export GO111MODULE=on
$ export GOPROXY=https://goproxy.cn,direct
~~~

测试服务是否正常
~~~
$ curl -XPOST -H "Content-Type: application/json; charset=utf8;" http://localhost:8090/v1/login -d '{"username":"admin", "password":"admin"}'
响应结果：{"code":0,"msg":"OK","data":{"token":"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIwLTAxLTE1VDA5OjU3OjIyLjA3ODAwMjk0OCswODowMCIsImlhdCI6MTU3OTA1MzQxMiwiaWQiOjEsIm5iZiI6MTU3OTA1MzQxMiwidXNlcm5hbWUiOiJhZG1pbiJ9.IT_X3ElBuUEksGGmnD57fDF3MFwnUDf74bAikaSdLqo"}}

$ curl -XPOST -H "Content-Type: application/json; charset=utf8;" http://localhost:8090/v1/user -d '{"username":"zhangsan", "password":"123456", "age":18}' -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIwLTAxLTE1VDA5OjU3OjIyLjA3ODAwMjk0OCswODowMCIsImlhdCI6MTU3OTA1MzQxMiwiaWQiOjEsIm5iZiI6MTU3OTA1MzQxMiwidXNlcm5hbWUiOiJhZG1pbiJ9.IT_X3ElBuUEksGGmnD57fDF3MFwnUDf74bAikaSdLqo"
输出结果：{"code":0,"msg":"OK","data":{"username":"zhangsan"}}

$ curl -XGET -H "Content-Type: application/json; charset=utf8;" http://localhost:8090/v1/user/zhangsan -H "Authorization:Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOiIyMDIwLTAxLTE1VDA5OjU3OjIyLjA3ODAwMjk0OCswODowMCIsImlhdCI6MTU3OTA1MzQxMiwiaWQiOjEsIm5iZiI6MTU3OTA1MzQxMiwidXNlcm5hbWUiOiJhZG1pbiJ9.IT_X3ElBuUEksGGmnD57fDF3MFwnUDf74bAikaSdLqo"
输出结果：{"code":0,"msg":"OK","data":{"age":18,"username":"zhangsan"}}
~~~

构建脚本运行项目
~~~
$ make gotest
$ make
$ make run
~~~

自动生成API文档
~~~
$ go get -u github.com/swaggo/swag/cmd/swag
$ swag init -g ./cmd/gin-api-demo/main.go
~~~