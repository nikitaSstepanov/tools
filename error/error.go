package e

import (
	"net/http"

	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

type Error *errorStruct

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

func New(msg string, status statusType) Error {
	return &errorStruct{
		Message: msg,
		Code:    status,
	}
}

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

func (e *errorStruct) ToGRPCErr() error {
	return status.Errorf(e.ToGRPCCode(), e.Message)
}

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
