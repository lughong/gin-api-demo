package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lughong/gin-api-demo/router/middleware"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	m := middleware.NewGoMiddleware()

	t.Run("CORS", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(resp)
		c.Request = httptest.NewRequest("GET", "/", nil)

		h := m.CORS()
		h(c)

		result := resp.Result()
		t.Log(result)
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
		t.Log(result)
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
		t.Log(result)
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
		t.Log(result)
		assert.Equal(t, http.StatusOK, result.StatusCode)
		assert.Equal(t, "*", result.Header.Get("Access-Control-Allow-Origin"))
	})
}
