package apperror

import "errors"

const (
	CodeInternal                = "INTERNAL_ERROR"
	CodeInvalidJSON             = "INVALID_JSON"
	CodeValidationFailed        = "VALIDATION_FAILED"
	CodeInvalidID               = "INVALID_ID"
	CodeInvalidQuery            = "INVALID_QUERY"
	CodeUnauthorized            = "UNAUTHORIZED"
	CodeAuthInvalidCredentials  = "AUTH_INVALID_CREDENTIALS"
	CodeAuthInvalidRefreshToken = "AUTH_INVALID_REFRESH_TOKEN"
	CodeProjectNotFound         = "PROJECT_NOT_FOUND"
	CodeUploadFileMissing       = "UPLOAD_FILE_MISSING"
	CodeUploadFileTooLarge      = "UPLOAD_FILE_TOO_LARGE"
	CodeUploadInvalidType       = "UPLOAD_INVALID_TYPE"
	CodeOrderIPBlocked          = "ORDER_IP_BLOCKED"
	CodeOrderRateLimited        = "ORDER_RATE_LIMITED"
	CodeOrderNotFound           = "ORDER_NOT_FOUND"
	CodeReviewInvalidPhone      = "REVIEW_INVALID_PHONE"
	CodeReviewNotAllowed        = "REVIEW_NOT_ALLOWED"
)

type Error struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Meta    map[string]any `json:"meta,omitempty"`
	cause   error
}

func New(code, message string, meta map[string]any) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Meta:    meta,
	}
}

func Wrap(err error, code, message string, meta map[string]any) *Error {
	appErr := New(code, message, meta)
	appErr.cause = err
	return appErr
}

func Internal(err error) *Error {
	return Wrap(err, CodeInternal, "Internal server error", nil)
}

func From(err error) (*Error, bool) {
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.cause != nil {
		return e.Message + ": " + e.cause.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}
