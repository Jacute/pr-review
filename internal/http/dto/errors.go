package dto

type ErrorCode string

type Response struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func Error(code ErrorCode, msg string) *Response {
	return &Response{
		Code:    code,
		Message: msg,
	}
}
