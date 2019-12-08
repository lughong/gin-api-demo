package middleware

import (
	"fmt"
	"regexp"

	"github.com/lughong/gin-api-demo/entity"
	"github.com/lughong/gin-api-demo/pkg/errno"
	"github.com/lughong/gin-api-demo/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		fmt.Println(c.Get("Authorization"))

		reg := regexp.MustCompile(`(/static|/login|/favicon.ico)`)
		if reg.MatchString(path) {
			return
		}

		// parse the JSON web token.
		if _, err := token.ParseRequest(c); err != nil {
			entity.SendResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
