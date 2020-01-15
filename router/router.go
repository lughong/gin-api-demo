package router

import (
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/lughong/gin-api-demo/docs"
	_http "github.com/lughong/gin-api-demo/handler/http"
	"github.com/lughong/gin-api-demo/model"
	"github.com/lughong/gin-api-demo/registry"
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
	g.Use(mw...)

	// 404
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	return &Router{
		g: g,
	}
}

// Run 运行路由
func (r *Router) Run(ctn *registry.Container) error {
	r.g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	pprof.Register(r.g)

	_http.NewUserHandler(r.g, ctn.Resolve("user-logic").(model.UserLogic))

	return r.g.Run(viper.GetString("server.port"))
}
