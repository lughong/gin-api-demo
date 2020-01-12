package errno_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lughong/gin-api-demo/pkg/errno"
)

func TestAdd(t *testing.T) {
	addStr := "test add func."
	errStr := "Test errno. "
	err := errors.New(errStr)

	expected := fmt.Sprintf(
		"Err - code: %d, message: %s %s, error: %s",
		errno.ErrPasswordIncorrect.Code,
		errno.ErrPasswordIncorrect.Message,
		addStr,
		errStr,
	)

	e := errno.New(errno.ErrPasswordIncorrect, err)
	_ = e.Add(addStr)

	if e.Error() != expected {
		t.Error("actual not equal expected.")
	}
}

func BenchmarkAdd(b *testing.B) {
	expected := "test"
	e := new(errno.Err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = e.Add(expected)
	}
}