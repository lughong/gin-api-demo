package router

import (
	"net/http"
	"time"

	"github.com/lughong/gin-api-demo/api/user"
	_handler "github.com/lughong/gin-api-demo/api/user/handler/http"
	_ "github.com/lughong/gin-api-demo/docs"
	"github.com/lughong/gin-api-demo/registry"
	"github.com/lughong/gin-api-demo/router/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Router struct {
	g *gin.Engine
}

// NewRouter 创建路由
func NewRouter(mw []gin.HandlerFunc) *Router {
	// 设置运行模式
	gin.SetMode(viper.GetString("server.runMode"))

	g := gin.New()
	// 设置panic恢复中间件
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

	return &Router{
		g: g,
	}
}

// Run 运行路由
func (r *Router) Run(ctn *registry.Container) error {
	// 404
	r.g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	r.g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	pprof.Register(r.g)

	_handler.NewUserHandler(r.g, ctn.Resolve("user-logic").(user.Logic))

	return r.g.Run(viper.GetString("server.port"))
}
