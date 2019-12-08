package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lughong/gin-api-demo/config"
	"github.com/lughong/gin-api-demo/model"
	"github.com/lughong/gin-api-demo/pkg/redis"
	version2 "github.com/lughong/gin-api-demo/pkg/version"
	"github.com/lughong/gin-api-demo/router"
	"github.com/lughong/gin-api-demo/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	cfg     = pflag.StringP("config", "c", "conf/config.yaml", "config file path.")
	version = pflag.BoolP("version", "v", false, "show version info.")
)

func main() {
	pflag.Parse()

	if *version {
		v := version2.Get()
		marshalled, err := json.MarshalIndent(&v, "", " ")
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(marshalled))
		return
	}

	c := config.New(func(c *config.Config) {
		c.Name = *cfg
	})
	// load config
	if err := c.Load(); err != nil {
		fmt.Printf("Config load. %s", err)
		os.Exit(1)
	}

	// init database
	db, err := model.Init()
	if err != nil {
		fmt.Printf("Model init. %s", err)
		os.Exit(1)
	}
	defer db.Close()

	// init redis
	redis.Init()

	// init router
	mw := []gin.HandlerFunc{
		middleware.RequestId(),
		middleware.LoggerToFile(),
	}
	g := router.Init(mw)

	logrus.Info("Start...")
	if err := g.Run(viper.GetString("server.port")); err != nil {
		logrus.Fatalf("Gin Run. %s", err)
	}
}
