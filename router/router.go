package router

import (
	"net/http"

	"github.com/lughong/gin-api-demo/api/user"
	_handler "github.com/lughong/gin-api-demo/api/user/handler/http"
	_ "github.com/lughong/gin-api-demo/docs"
	"github.com/lughong/gin-api-demo/registry"

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

	_handler.NewUserHandler(r.g, ctn.Resolve("user-logic").(user.Logic))

	return r.g.Run(viper.GetString("server.port"))
}
