package middleware

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/lughong/gin-api-demo/pkg/errno"
	"github.com/lughong/gin-api-demo/pkg/token"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware 验证token中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		fmt.Println(c.Get("Authorization"))

		reg := regexp.MustCompile(`(/static|/login|/favicon.ico)`)
		if reg.MatchString(path) {
			return
		}

		response := Response{
			Code: errno.ErrTokenInvalid.Code,
			Msg:  errno.ErrTokenInvalid.Message,
			Data: nil,
		}
		// parse the JSON web token.
		if _, err := token.ParseRequest(c); err != nil {
			c.JSON(http.StatusOK, response)
			c.Abort()
			return
		}

		c.Next()
	}
}
