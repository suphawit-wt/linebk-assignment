package api

// API status constants represent standardized application-level response codes.
const (
	StatusSuccess              = "SUCCESS"
	StatusInvalidInput         = "INVALID_INPUT"
	StatusBadRequest           = "BAD_REQUEST"
	StatusUnauthorized         = "UNAUTHORIZED"
	StatusTokenNotMatched      = "TOKEN_NOT_MATCHED"
	StatusExpiredOTP           = "EXPIRED_OTP"
	StatusForbidden            = "FORBIDDEN"
	StatusNotFound             = "NOT_FOUND"
	StatusRequestTimeout       = "REQUEST_TIMEOUT"
	StatusConflict             = "CONFLICT"
	StatusFileTooLarge         = "FILE_TOO_LARGE"
	StatusUnsupportedMediaType = "UNSUPPORTED_MEDIA_TYPE"
	StatusUnprocessableContent = "UNPROCESSABLE_CONTENT"
	StatusTooManyRequests      = "TOO_MANY_REQUESTS"
	StatusUnexpected           = "UNEXPECTED"
)
