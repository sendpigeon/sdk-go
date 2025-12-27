package sendpigeon

import "fmt"

// ErrorCode represents the type of error.
type ErrorCode string

const (
	ErrorCodeNetwork ErrorCode = "network_error"
	ErrorCodeAPI     ErrorCode = "api_error"
	ErrorCodeTimeout ErrorCode = "timeout_error"
)

// Error represents an error from the SendPigeon API or SDK.
type Error struct {
	Message string    `json:"message"`
	Code    ErrorCode `json:"code"`
	APICode string    `json:"api_code,omitempty"`
	Status  int       `json:"status,omitempty"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.APICode != "" {
		return fmt.Sprintf("[%s] %s", e.APICode, e.Message)
	}
	return e.Message
}

// NewError creates a new Error.
func NewError(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// NewAPIError creates a new API error with status and API code.
func NewAPIError(status int, apiCode, message string) *Error {
	return &Error{
		Code:    ErrorCodeAPI,
		Status:  status,
		APICode: apiCode,
		Message: message,
	}
}
