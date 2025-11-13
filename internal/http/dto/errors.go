package dto

var (
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"
	ErrCodeInternal   ErrorCode = "INTERNAL_ERROR"
)

type ErrorCode string

type ErrorResponse struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
}

func Error(code ErrorCode, msg string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: msg,
	}
}
