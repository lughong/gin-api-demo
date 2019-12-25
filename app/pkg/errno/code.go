package errno

var (
	// 系统类型错误码
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	// 校验类型错误码
	ErrValidation = &Errno{Code: 20001, Message: "Validation failed."}
	ErrToken      = &Errno{Code: 20003, Message: "Error occurred while signing the JSON web token."}

	// 用户类型错误码
	ErrTokenInvalid      = &Errno{Code: 20101, Message: "The token was invalid."}
	ErrUserNotFound      = &Errno{Code: 20102, Message: "User was not found."}
	ErrGetUserDetail     = &Errno{Code: 20103, Message: "Get user detail was fail."}
	ErrPasswordIncorrect = &Errno{Code: 20104, Message: "Password incorrect."}
)
