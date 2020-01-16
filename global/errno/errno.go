package errno

import (
	"bytes"
	"fmt"
)

// Errno 结构体 Code错误码、Message错误信息
type Errno struct {
	Code    int
	Message string
}

func (e *Errno) Error() string {
	return e.Message
}

type Err struct {
	Code int
	Err  error

	buf bytes.Buffer
}

func New(errno *Errno, err error) *Err {
	e := &Err{
		Code: errno.Code,
		Err:  err,
	}

	e.buf.WriteString(errno.Message)

	return e
}

// Error 实现error接口，输出错误信息。
func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.buf.String(), err.Err)
}

// Add 新增错误信息
func (err *Err) Add(message string) error {
	err.buf.WriteString(" ")
	err.buf.WriteString(message)

	return err
}

// Addf 追加错误信息
func (err *Err) Addf(format string, args ...interface{}) error {
	err.buf.WriteString(" ")
	err.buf.WriteString(fmt.Sprintf(format, args...))

	return err
}

// IsErrUserNotFound 判断是否是用户不存在错误
func IsErrUserNotFound(err error) bool {
	code, _ := DecodeErr(err)
	return code == ErrUserNotFound.Code
}

// DecodeErr 解析err是属于那种错误类型，并返回其错误代码和错误信息。
func DecodeErr(err error) (int, string) {
	if err == nil {
		return OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Err:
		return typed.Code, typed.buf.String()
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return InternalServerError.Code, InternalServerError.Message
}
