package middleware

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/lughong/gin-api-demo/model"
	"github.com/lughong/gin-api-demo/pkg/errno"
	"github.com/lughong/gin-api-demo/pkg/token"

	"github.com/gin-gonic/gin"
)

// Auth 验证token中间件
func (m *GoMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		fmt.Println(c.Get("Authorization"))

		reg := regexp.MustCompile(`(/static|/login|/favicon.ico)`)
		if reg.MatchString(path) {
			return
		}

		response := model.Response{
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
