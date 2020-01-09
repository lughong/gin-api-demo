package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS 处理CORS中间件
func (m *GoMiddleware) CORS() gin.HandlerFunc {
	/*return cors.New(cors.Config{
		AllowAllOrigins:  viper.GetBool("cors.allowAllOrigins"),
		AllowMethods:     viper.GetStringSlice("cors.allowMethods"),
		AllowHeaders:     viper.GetStringSlice("cors.allowHeaders"),
		ExposeHeaders:    viper.GetStringSlice("cors.exposeHeaders"),
		AllowCredentials: viper.GetBool("cors.allowCredentials"),
		MaxAge:           time.Duration(viper.GetInt64("cors.maxAge")) * time.Hour,
	})*/

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Access-Control-Expose-Headers", "Cache-Control, Content-Language, Content-Type, Expires, Last-Modified, Pragma")
		c.Header("Access-Control-Allow-Credentials", "false")
		c.Header("Access-Control-Allow-Max-Age", "3000")

		c.Next()
	}
}

// Options 为Options请求附加头部信息
func (m *GoMiddleware) Options() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != "OPTIONS" {
			c.Next()
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
			c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
			c.Header("Content-Type", "application/json; charset=utf-8")

			c.AbortWithStatus(200)
		}
	}
}

// NoCache 在header头追加设置不缓存的信息，以防止客户端缓存HTTP响应的中间件。
func (m *GoMiddleware) NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
		c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
		c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
		c.Header("Content-Type", "application/json; charset=utf-8")

		c.Next()
	}
}

// Secure 附加安全性和资源访问头的中间件。
func (m *GoMiddleware) Secure() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000")
		}
		// 还可以考虑添加Content-Security-Policy头信息
		// c.Header("Content-Security-Policy", "script-src 'self' https://cdnjs.cloudflare.com")

		c.Next()
	}
}
