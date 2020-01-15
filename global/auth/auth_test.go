package auth_test

import (
	"testing"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/assert"

	"github.com/lughong/gin-api-demo/global/auth"
)

func TestEncrypt(t *testing.T) {
	var expected string
	err := faker.FakeData(&expected)
	assert.NoError(t, err)

	actual, err := auth.Encrypt(expected)
	assert.NoError(t, err)
	assert.NotNil(t, actual)
}

func BenchmarkEncrypt(b *testing.B) {
	var expected string
	err := faker.FakeData(&expected)
	assert.NoError(b, err)
	_, _ = auth.Encrypt(expected)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = auth.Encrypt(expected)
	}
}

func BenchmarkCompare(b *testing.B) {
	var expected string
	err := faker.FakeData(&expected)
	assert.NoError(b, err)
	hasStr, err := auth.Encrypt(expected)
	assert.NoError(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = auth.Compare(hasStr, expected)
	}
}
