package middleware

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"

	"github.com/lughong/gin-api-demo/global/errno"
	"github.com/lughong/gin-api-demo/global/token"
)

// Auth 验证token中间件
func (m *GoMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		reg := regexp.MustCompile(`(/static|/login|/favicon.ico)`)
		if reg.MatchString(path) {
			return
		}

		// parse the JSON web token.
		if _, err := token.ParseRequest(c); err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": errno.ErrTokenInvalid.Code,
				"msg":  errno.ErrTokenInvalid.Message,
				"data": "",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
