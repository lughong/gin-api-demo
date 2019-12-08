package router

import (
	v1 "github.com/lughong/gin-api-demo/handler/v1"
	"net/http"
	"time"

	"github.com/lughong/gin-api-demo/router/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init(mw []gin.HandlerFunc) *gin.Engine {
	gin.SetMode(viper.GetString("server.runMode"))

	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(middleware.Secure)
	g.Use(middleware.NoCache)
	g.Use(mw...)

	// 跨域资源共享 CORS 配置
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:  viper.GetBool("cors.allowAllOrigins"),
		AllowMethods:     viper.GetStringSlice("cors.allowMethods"),
		AllowHeaders:     viper.GetStringSlice("cors.allowHeaders"),
		ExposeHeaders:    viper.GetStringSlice("cors.exposeHeaders"),
		AllowCredentials: viper.GetBool("cors.allowCredentials"),
		MaxAge:           time.Duration(viper.GetInt64("cors.maxAge")) * time.Hour,
	}))

	V1 := g.Group("v1")
	{
		V1.GET("user", v1.Get)
	}

	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	return g
}
