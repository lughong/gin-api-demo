package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// NoCache 在header头追加设置不缓存的信息，以防止客户端缓存HTTP响应的中间件。
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))

	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Next()
}

// Secure 附加安全性和资源访问头的中间件。
func Secure(c *gin.Context) {
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
