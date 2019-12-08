package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"time"

	"github.com/lughong/gin-api-demo/entity"
	"github.com/lughong/gin-api-demo/pkg/errno"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logging is a middleware function that logs the each request.
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqURI := c.Request.RequestURI
		reg := regexp.MustCompile(`(/v1/login|/favicon.ico)`)
		if reg.MatchString(reqURI) {
			return
		}

		// read the body context
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		}

		// Restore the io.ReadCloser to its original state
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

		var response entity.Response
		if err := json.Unmarshal(blw.body.Bytes(), &response); err != nil {
			statusCode = errno.InternalServerError.Code
			message = errno.InternalServerError.Message
		} else {
			statusCode = response.Code
			message = response.Msg
		}

		// the basic information
		clientIP := c.ClientIP()
		reqMethod := c.Request.Method

		logrus.Infof(
			"%s | %s | %s | %s | %s | {code: %d, message: %s} ",
			latencyTime,
			clientIP,
			reqMethod,
			reqURI,
			bodyBytes,
			statusCode,
			message,
		)
	}
}

func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
