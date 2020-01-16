package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/lughong/gin-api-demo/global/errno"
)

// bodyLogWriter结构体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 这个方法覆盖了*bytes.Buffer里面的Write方法。
// 目的是为了截取响应信息。
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerToFile 把日志记录在 File。这是一个中间件函数，可以记录每一次客户端请求的信息。
func (m *GoMiddleware) LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取请求body内容
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		// 因为请求body只能读取一次，所以读取后需要重写入到request里面。
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		startTime := time.Now().UTC()

		c.Next()

		endTime := time.Now().UTC()
		latencyTime := endTime.Sub(startTime)

		var statusCode, message = -1, ""

		// 读取响应信息
		var response struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}
		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			statusCode = errno.InternalServerError.Code
			message = errno.InternalServerError.Message
		} else {
			statusCode = response.Code
			message = response.Msg
		}

		// the basic information
		reqURI := c.Request.RequestURI
		clientIP := c.ClientIP()
		reqMethod := c.Request.Method

		logrus.Infof(
			"| %s | %s | %s | %s | %d | %s | %s | %s | {code: %d, message: %s} |",
			clientIP,
			reqMethod,
			reqURI,
			c.Request.Proto,
			c.Writer.Status(),
			latencyTime,
			c.Request.UserAgent(),
			bodyBytes,
			statusCode,
			message,
		)
	}
}

// LoggerToMongo 把日志记录在 MongoDB
func (m *GoMiddleware) LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// LoggerToMQ 把日志记录在 MQ
func (m *GoMiddleware) LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
