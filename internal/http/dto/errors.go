package dto

var (
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"
	ErrCodeInternal   ErrorCode = "INTERNAL_ERROR"
)

var (
	ErrContentTypeNotJson = Error(
		ErrCodeBadRequest,
		"Content-Type must be application/json",
	)
	ErrInvalidBody = Error(
		ErrCodeBadRequest,
		"invalid body",
	)
	ErrInternal = Error(
		ErrCodeInternal,
		"internal error",
	)
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
