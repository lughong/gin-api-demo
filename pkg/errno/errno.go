package errno

import "fmt"

// Errno 结构体 Code错误码、Message错误信息
type Errno struct {
	Code    int
	Message string
}

func (e *Errno) Error() string {
	return e.Message
}

type Err struct {
	Code    int
	Message string
	Err     error
}

func New(errno *Errno, err error) *Err {
	return &Err{Code: errno.Code, Message: errno.Message, Err: err}
}

// Error 实现error接口，输出错误信息。
func (err *Err) Error() string {
	return fmt.Sprintf("Err - code: %d, message: %s, error: %s", err.Code, err.Message, err.Err)
}

// Add 新增错误信息
func (err *Err) Add(message string) error {
	err.Message += " " + message
	return err
}

// Addf 追加错误信息
func (err *Err) Addf(format string, args ...interface{}) error {
	err.Message += " " + fmt.Sprintf(format, args...)
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
		return typed.Code, typed.Message
	case *Errno:
		return typed.Code, typed.Message
	default:
	}

	return InternalServerError.Code, InternalServerError.Message
}
