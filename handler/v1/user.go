package v1

import (
	"github.com/lughong/gin-api-demo/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/lughong/gin-api-demo/entity"
)

func Get(c *gin.Context) {
	entity.SendResponse(c, errno.HelloWorld, nil)
}
