package errno

var (
	// Common errors
	OK                  = &Errno{Code: 0, Message: "OK"}
	InternalServerError = &Errno{Code: 10001, Message: "Internal server error"}
	ErrBind             = &Errno{Code: 10002, Message: "Error occurred while binding the request body to the struct."}

	// user errors
	HelloWorld		= &Errno{Code: 20100, Message: "hello world"}
	ErrTokenInvalid = &Errno{Code: 20101, Message: "The token was invalid."}
	ErrUserNotFound = &Errno{Code: 20102, Message: "User was not found."}
)
