package router

import (
	"net/http"
	"time"

	v1 "github.com/lughong/gin-api-demo/app/handler/v1"
	"github.com/lughong/gin-api-demo/app/router/middleware"
	_ "github.com/lughong/gin-api-demo/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Init 初始化路由器
func Init(mw []gin.HandlerFunc) *gin.Engine {
	// 设置运行模式
	gin.SetMode(viper.GetString("server.runMode"))

	g := gin.New()

	// 设置panic恢复中间件
	g.Use(gin.Recovery())
	g.Use(middleware.Secure)
	g.Use(middleware.NoCache)
	g.Use(mw...)
	// 404
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	pprof.Register(g)

	// 跨域资源共享 CORS 配置
	g.Use(cors.New(cors.Config{
		AllowAllOrigins:  viper.GetBool("cors.allowAllOrigins"),
		AllowMethods:     viper.GetStringSlice("cors.allowMethods"),
		AllowHeaders:     viper.GetStringSlice("cors.allowHeaders"),
		ExposeHeaders:    viper.GetStringSlice("cors.exposeHeaders"),
		AllowCredentials: viper.GetBool("cors.allowCredentials"),
		MaxAge:           time.Duration(viper.GetInt64("cors.maxAge")) * time.Hour,
	}))

	V1 := g.Group("/v1")
	{
		V1.GET("/user", v1.GetUser)
	}

	return g
}
