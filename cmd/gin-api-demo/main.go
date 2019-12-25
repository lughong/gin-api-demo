package main

import (
	"encoding/json"
	"log"

	"github.com/lughong/gin-api-demo/app/config"
	"github.com/lughong/gin-api-demo/app/model"
	"github.com/lughong/gin-api-demo/app/pkg/redis"
	version2 "github.com/lughong/gin-api-demo/app/pkg/version"
	"github.com/lughong/gin-api-demo/app/router"
	"github.com/lughong/gin-api-demo/app/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg     = pflag.StringP("config", "c", "app/conf/config.yaml", "config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

// @title gin-api-demo Example API
// @version 1.0
// @description gin api demo

// @contact.name lughong
// @contact.url http://www.swagger.io/support
// @contact.email 1586668924@qq.com

// @host localhost:8090
// @BasePath /v1
func main() {
	pflag.Parse()

	// 获取版本信息并输出其内容
	if *version {
		v := version2.Get()
		marshalled, err := json.MarshalIndent(&v, "", " ")
		if err != nil {
			log.Fatalf("%v\r\n", err)
		}

		log.Println(string(marshalled))
		return
	}

	// 设置配置文件路径
	c := config.New(func(c *config.Config) {
		c.Name = *cfg
	})
	// 加载配置文件信息
	if err := c.Load(); err != nil {
		log.Fatalf("Config load. %s", err)
	}

	// 初始化数据库
	db, err := model.Init()
	if err != nil {
		log.Printf("Model init. %s", err)
	}
	if db != nil {
		defer db.Close()
	}

	// 初始化redis
	redis.Init()

	// 设置路由中间件
	mw := []gin.HandlerFunc{
		middleware.RequestId(),
		middleware.LoggerToFile(),
	}
	// 初始化路由器
	g := router.Init(mw)

	// 开启服务
	log.Println("Api server start...")
	if err := g.Run(viper.GetString("server.port")); err != nil {
		log.Fatalf("Gin run. %s", err)
	}
}
