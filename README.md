项目介绍
~~~
使用gin开发的一个api框架。
项目地址：https://github.com/lughong/gin-api-demo
~~~

目录结构
~~~
├── api                        # 接口目录
│   ├── user                   # 实际业务处理函数存放位置
│   │   ├── handler            # 类似MVC架构中的C，用来读取输入，并将处理流程转发给实际的处理函数（CLI、web、REST、gRPC)
│   │   ├── repository         # 仓库实现层（NoSQL、RDBMS、Micro-Services）
│   │   ├── logic              # 逻辑实现层
│   │   ├── repository.go      # 仓库层接口
│   │   └── logic.go           # 逻辑层接口
│   └── response.go            
├── cmd                        
│   └── gin-api-demo           # Go程序唯一入口      
├── conf                       # 配置文件统一存放目录(运行前需要调整database、redis等参数配置)
├── config                     # 专门用来处理配置和配置文件的Go package
├── data                       
│   └── logs                   # 运行时记录的日志目录
├── docs                       # api文档存放目录
├── model                      # 数据库相关的操作统一放在这里，包括数据库初始化和对表的增删改查
├── pkg                        # 引用的包
│   ├── auth                   # 认证包
│   ├── constvar               # 常量统一存放位置
│   ├── errno                  # 错误码存放位置
│   ├── redis                  # redis包
│   ├── token                  # token包
│   └── version                # 版本包
├── registry                   # 依赖注入容器
├── router                     # 路由相关处理
│   ├── middleware             # API服务器用的是Gin Web框架，Gin中间件存放位置
│   └── router.go
├── util                       # 工具类函数存放目录
├── gin_api_demo.sql           # 在部署新环境时，可以登录MySQL客户端，执行source gin_api_demo.sql创建数据库和表
├── Makefile                   # Makefile文件
├── README.md                  
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
mysql> source 项目根目录/data/sql/gin_api_demo.sql
mysql> CREATE user 'apiuser'@'127.0.0.1' IDENTIFIED BY '123456';
mysql> GRANT ALL PRIVILEGES ON gin_api_demo.* TO 'apiuser'@'127.0.0.1';
mysql> FLUSH PRIVILEGES;
~~~

运行项目
~~~
$ go run main.go
~~~

运行过程若出现module失败，可以尝试设置GOPROXY环境变量
~~~
$ export GO111MODULE=on
$ export GOPROXY=https://goproxy.cn,direct
~~~

测试服务是否正常
~~~
$ curl -XGET -H "Content-Type: application/json; charset=utf8;" -d'{"username":"zhangsan","password":""}' http://localhost:8090/user
输出结果：{"code":0,"msg":"OK","data":{"age":18,"username":"zhangsan"}}
~~~

构建脚本运行项目
~~~
$ make clean
$ make
$ ./gin-api-demo
~~~