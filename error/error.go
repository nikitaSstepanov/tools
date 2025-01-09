package e

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/nikitaSstepanov/tools/ctx"
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

	// GetError returns the underlying error.
	// This method provides access to the original error that may contain more details.
	GetError() error

	GetTag(key string) interface{}

	// GetCode returns the status code associated with the error.
	// The StatusType can be a custom type that represents various error codes.
	GetCode() StatusType

	Log(msg ...string)

	// WithMessage sets a new error message for the error instance.
	// This method allows users to update the error message dynamically.
	WithMessage(string) Error

	// WithErr sets a new underlying error for the error instance.
	// This method allows users to associate a different error with this custom error type.
	WithErr(error) Error

	WithTag(key string, value interface{}) Error

	WithCtx(c ctx.Context) Error

	// WithCode sets a new status code for the error instance.
	// This method allows users to update the error code dynamically.
	WithCode(StatusType) Error

	// ToJson() returns erro struct with json tags.
	ToJson() JsonError

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
	message string
	errs    []error
	tags    map[string]interface{}
	code    StatusType
	log     *slog.Logger
}

// New returns type Error with message.
func New(msg string, status StatusType, errs ...error) Error {
	if errs == nil {
		errs = []error{}
	}

	return &errorStruct{
		message: msg,
		errs:    errs,
		tags:    make(map[string]interface{}),
		code:    status,
		log:     slog.Default(),
	}
}

func (e *errorStruct) GetMessage() string {
	return e.message
}

func (e *errorStruct) GetError() error {
	return errors.Join(e.errs...)
}

func (e *errorStruct) GetTag(key string) interface{} {
	return e.tags[key]
}

func (e *errorStruct) GetCode() StatusType {
	return e.code
}

func (e *errorStruct) Log(msg ...string) {
	l := e.log

	l = l.With(e.SlErr())

	for key, value := range e.tags {
		l = l.With(key, value)
	}

	message := ""

	if len(msg) != 0 {
		message = strings.Join(msg, " ")
	}

	l.Error(message)
}

func (e *errorStruct) WithMessage(msg string) Error {
	return New(msg, e.code, e.errs...)
}

func (e *errorStruct) WithErr(err error) Error {
	return New(e.message, e.code, append(e.errs, err)...)
}

func (e *errorStruct) WithTag(key string, value interface{}) Error {
	err := New(e.message, e.code, e.errs...).(*errorStruct)

	for key, value := range e.tags {
		err.tags[key] = value
	}

	err.tags[key] = value

	return err
}

func (e *errorStruct) WithCtx(c ctx.Context) Error {
	err := New(e.message, e.code, e.errs...).(*errorStruct)

	ctxErr := c.Err()
	if ctxErr != nil {
		err.errs = append(err.errs, c.Err())
	}

	for key, value := range e.tags {
		err.tags[key] = value
	}

	for key, value := range c.GetValues() {
		if value.Share {
			err.tags[key] = value.Val
		}
	}

	err.log = slog.New(c.SlHandler())

	c.AddErr(err)

	return err
}

func (e *errorStruct) WithCode(status StatusType) Error {
	return New(e.message, status, e.errs...)
}

func (e *errorStruct) ToJson() JsonError {
	return JsonError{
		Error: e.message,
	}
}

// ToHttpCode convert Error to http status code.
func (e *errorStruct) ToHttpCode() int {
	switch e.code {

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
	if e.message == "" && (e.errs == nil || (len(e.errs) == 0)) {
		return "nil"
	}
	if e.errs == nil || (len(e.errs) == 0) {
		return e.message
	}
	if e.message == "" {
		return errors.Join(e.errs...).Error()
	}
	return e.message + ": " + errors.Join(e.errs...).Error()
}

func (e *errorStruct) ToGRPCErr() error {
	return status.Error(e.ToGRPCCode(), e.message)
}

func FromGRPCErr(err error) Error {
	stat, _ := status.FromError(err)

	var code StatusType

	switch stat.Code() {

	case codes.Internal:
		code = Internal

	case codes.NotFound:
		code = NotFound

	case codes.InvalidArgument:
		code = BadInput

	case codes.Unauthenticated:
		code = Unauthorize

	case codes.AlreadyExists:
		code = Conflict

	default:
		code = Internal

	}

	return New(stat.Message(), code)
}

// ToGRPCCode convert Error to grpc status code.
func (e *errorStruct) ToGRPCCode() codes.Code {
	switch e.code {

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

func (e *errorStruct) SlErr() slog.Attr {
	return slog.String("error", e.Error())
}

// E creates a new custom error instance if the provided error is not nil.
// It initializes the custom error with an empty message and associates the given error
// with an internal status code by default. If the provided error is nil, it returns nil.
func E(err error) Error {
	if err != nil {
		if _, ok := err.(Error); ok {
			return err.(Error)
		}

		return New("", Internal, err)
	}

	return nil
}
