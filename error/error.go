package e

import (
	"log/slog"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Error interface defines a custom error type that provides additional context
// and functionality beyond the standard error interface in Go. It is designed
// to encapsulate error messages, status codes, and conversion methods for gRPC
// and HTTP responses.
type Error interface {
	// GetMessage returns the error message as a string.
	// This method allows users to retrieve a human-readable description of the error.
	GetMessage() string

	// GetCode returns the status code associated with the error.
	// The statusType can be a custom type that represents various error codes.
	GetCode() statusType

	// WithMessage sets a new error message for the error instance.
	// This method allows users to update the error message dynamically.
	WithMessage(string)

	// WithCode sets a new status code for the error instance.
	// This method allows users to update the error code dynamically.
	WithCode(statusType)

	// ToGRPCCode converts the error's status code to a gRPC error code.
	// This method facilitates interoperability with gRPC services by providing
	// an appropriate error code representation.
	ToGRPCCode() codes.Code

	// ToHttpCode converts the error's status code to an HTTP status code.
	// This method helps in mapping application-specific errors to standard HTTP responses.
	ToHttpCode() int

	// Error returns the string representation of the error.
	// This method implements the standard error interface, allowing the error
	// to be used in contexts where a simple error message is required.
	Error() string

	// SlErr returns structured logging attributes for the error.
	// This method provides a way to log the error with additional context,
	// making it easier to analyze issues in logs.
	SlErr() slog.Attr

	// ToGRPCErr converts the custom error into a standard Go error type suitable for gRPC.
	// This method allows seamless integration with gRPC error handling mechanisms.
	ToGRPCErr() error
}

type errorStruct struct {
	Message string     `json:"error"`
	Code    statusType `json:"-"`
}

const (
	Internal statusType = iota
	NotFound
	BadInput
	Conflict
	Forbidden
	Unauthorize
)

type statusType int

// New returns type Error with message.
func New(msg string, status statusType) Error {
	return &errorStruct{
		Message: msg,
		Code:    status,
	}
}

func (e *errorStruct) GetMessage() string {
	return e.Message
}

func (e *errorStruct) GetCode() statusType {
	return e.Code
}

func (e *errorStruct) WithMessage(msg string) {
	e.Message = msg
}

func (e *errorStruct) WithCode(status statusType) {
	e.Code = status
}

// ToHttpCode convert Error to http status code.
func (e *errorStruct) ToHttpCode() int {
	switch e.Code {

	case Internal:
		return http.StatusInternalServerError

	case NotFound:
		return http.StatusNotFound

	case BadInput:
		return http.StatusBadRequest

	case Unauthorize:
		return http.StatusUnauthorized

	case Forbidden:
		return http.StatusForbidden

	case Conflict:
		return http.StatusConflict

	default:
		return http.StatusInternalServerError

	}
}

func (e *errorStruct) Error() string {
	return e.Message
}

func (e *errorStruct) ToGRPCErr() error {
	return status.Errorf(e.ToGRPCCode(), e.Message)
}

// ToGRPCCode convert Error to grpc status code.
func (e *errorStruct) ToGRPCCode() codes.Code {
	switch e.Code {

	case Internal:
		return codes.Internal

	case NotFound:
		return codes.NotFound

	case BadInput:
		return codes.InvalidArgument

	case Unauthorize:
		return codes.Unauthenticated

	case Conflict:
		return codes.AlreadyExists

	default:
		return codes.Internal

	}
}

// SlErr returns slog.Attr with key "error" and err value.
func (e *errorStruct) SlErr() slog.Attr {
	return slog.String("error", e.Message)
}

func E(err error) Error {
	return New(err.Error(), Internal)
}

func EC(err error, code statusType) Error {
	return New(err.Error(), code)
}
