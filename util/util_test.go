package util_test

import (
	"net/http/httptest"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/lughong/gin-api-demo/util"
)

func TestGetReqID(t *testing.T) {
	var idTests = []struct {
		in       string
		expected string
	}{
		{"1", "1"},
		{"2", "2"},
		{"3", "3"},
		{"4", "4"},
		{"5", "5"},
		{"6", "6"},
	}

	resp := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(resp)
	for _, tt := range idTests {
		c.Set("X-Request-Id", tt.in)
		assert.Equal(t, tt.expected, util.GetReqID(c))
	}
}

func TestEncryptMD5(t *testing.T) {
	var expected string
	err := faker.FakeData(&expected)
	assert.NoError(t, err)

	actual, err := util.EncryptMD5(expected)

	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func BenchmarkEncryptMD5(b *testing.B) {
	var expected string
	err := faker.FakeData(&expected)
	assert.NoError(b, err)

	_, _ = util.EncryptMD5(expected)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = util.EncryptMD5(expected)
	}
}
