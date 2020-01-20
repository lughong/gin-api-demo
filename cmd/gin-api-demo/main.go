package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/lughong/gin-api-demo/config"
	"github.com/lughong/gin-api-demo/global/constvar"
	version2 "github.com/lughong/gin-api-demo/global/version"
	"github.com/lughong/gin-api-demo/registry"
	"github.com/lughong/gin-api-demo/router"
	"github.com/lughong/gin-api-demo/router/middleware"
	"github.com/lughong/gin-api-demo/util"
)

var (
	cfg     = pflag.StringP("config", "c", "config", "config file name.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func init() {
	pflag.Parse()

	// 设置根目录
	constvar.RootDir = interRootDir()

	// 设置配置文件路径
	c := config.NewConfig(func(c *config.Config) {
		c.Name = *cfg
	})
	// 加载配置文件信息
	if err := c.Load(); err != nil {
		log.Fatalf("Config load. %s", err)
	}
}

// @title Restful API
// @version 1.0
// @description 使用GO语言gin框架开发的Restful API example

// @contact.name lughong
// @contact.url https://github.com/lughong/gin-api-demo
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
	m := middleware.NewGoMiddleware()
	mw := []gin.HandlerFunc{
		m.Recover(),
		m.LoggerToFile(),
		m.CORS(),
		m.NoCache(),
		m.Secure(),
		m.RequestId(),
		m.Auth(),
	}
	r := router.NewRouter(mw)
	handler := r.InitHandler(ctn)

	if err := http.ListenAndServe(viper.GetString("server.addr"), handler); err != nil {
		log.Fatalf("Gin run. %s", err)
	}
}

// interRootDir 初始化项目根目录
func interRootDir() string {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Getwd. %s", err)
	}

	var inter func(d string) string
	inter = func(d string) string {
		if util.Exists(d + "/config") {
			return d
		}

		return inter(filepath.Dir(d))
	}

	return inter(cwd)
}
