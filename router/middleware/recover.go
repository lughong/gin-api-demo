package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/lughong/gin-api-demo/global/constvar"
	"github.com/lughong/gin-api-demo/global/errno"
	"github.com/lughong/gin-api-demo/util"
)

// Recover 捕获异常中间件
func (m *GoMiddleware) Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {

				DebugStack := ""
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + "<br>"
				}

				subject := fmt.Sprintf("【重要错误】%s 项目出错了！", viper.GetString("server.appName"))

				body := strings.ReplaceAll(constvar.MailTemplate, "{ErrorMsg}", fmt.Sprintf("%s", err))
				body = strings.ReplaceAll(body, "{RequestTime}", time.Now().Format("2006/01/02 15:04:05"))
				body = strings.ReplaceAll(body, "{RequestURL}", c.Request.Method+"  "+c.Request.Host+c.Request.RequestURI)
				body = strings.ReplaceAll(body, "{RequestUA}", c.Request.UserAgent())
				body = strings.ReplaceAll(body, "{RequestIP}", c.ClientIP())
				body = strings.ReplaceAll(body, "{DebugStack}", DebugStack)

				_ = util.SendMail(viper.GetStringSlice("mail.to"), subject, body)

				logrus.Error(DebugStack)

				c.JSON(http.StatusOK, gin.H{
					"code": errno.InternalServerError.Code,
					"msg":  errno.InternalServerError.Message,
					"data": "",
				})
				c.Abort()
				return
			}
		}()

		c.Next()
	}
}
