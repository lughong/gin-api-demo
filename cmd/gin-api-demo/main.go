package main

import (
	"encoding/json"
	"github.com/lughong/gin-api-demo/router"
	"log"

	"github.com/lughong/gin-api-demo/config"
	version2 "github.com/lughong/gin-api-demo/pkg/version"
	"github.com/lughong/gin-api-demo/registry"
	"github.com/lughong/gin-api-demo/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
)

var (
	cfg     = pflag.StringP("config", "c", "conf/config.yaml", "config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func init() {
	pflag.Parse()

	// 设置配置文件路径
	c := config.NewConfig(func(c *config.Config) {
		c.Name = *cfg
	})
	// 加载配置文件信息
	if err := c.Load(); err != nil {
		log.Fatalf("Config load. %s", err)
	}
}

// @title gin-api-demo Example API
// @version 1.0
// @description gin api demo

// @contact.name lughong
// @contact.url http://www.swagger.io/support
// @contact.email 1586668924@qq.com

// @host localhost:8090
// @BasePath /v1
func main() {
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

	// 配置注册表
	ctn, err := registry.NewContainer()
	if err != nil {
		log.Fatalf("registry NewContainer. %s", err)
	}
	defer ctn.Delete()

	// 配置路由
	mw := []gin.HandlerFunc{middleware.RequestId(), middleware.LoggerToFile()}
	r := router.NewRouter(mw)
	if err := r.Run(ctn); err != nil {
		log.Fatalf("Gin run. %s", err)
	}
}
