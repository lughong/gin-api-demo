package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/lughong/gin-api-demo/router/middleware"
)

func TestMain(m *testing.M) {
	src, _ := os.OpenFile("/var/log/go/system.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	// 设置日志输出
	logrus.SetOutput(src)

	m.Run()
}

func TestMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	m := middleware.NewGoMiddleware()

	t.Run("CORS", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(resp)
		c.Request = httptest.NewRequest("GET", "/", nil)

		h := m.CORS()
		h(c)

		result := resp.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "*", result.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("Options", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(resp)
		c.Request = httptest.NewRequest("OPTIONS", "/", nil)

		h := m.Options()
		h(c)

		result := resp.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "*", result.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("NoCache", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(resp)
		c.Request = httptest.NewRequest("GET", "/", nil)

		h := m.NoCache()
		h(c)

		result := resp.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "no-cache, no-store, max-age=0, must-revalidate, value", result.Header.Get("Cache-Control"))
	})

	t.Run("Secure", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(resp)
		c.Request = httptest.NewRequest("GET", "/", nil)

		h := m.Secure()
		h(c)

		result := resp.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "*", result.Header.Get("Access-Control-Allow-Origin"))
	})

	t.Run("LoggerToFile", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(resp)
		c.Request = httptest.NewRequest("GET", "/", nil)

		h := m.LoggerToFile()
		h(c)

		result := resp.Result()
		assert.Equal(t, http.StatusOK, result.StatusCode)
	})
}
