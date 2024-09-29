package e

import (
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	msg := "An error occurred"
	code := NotFound

	err := New(msg, code)

	if err.GetMessage() != msg {
		t.Errorf("Expected message %q, got %q", msg, err.GetMessage())
	}

	if err.GetCode() != code {
		t.Errorf("Expected code %v, got %v", code, err.GetCode())
	}
}

func TestWithMessage(t *testing.T) {
	initialMsg := "Initial error"
	code := Internal
	err := New(initialMsg, code)

	newMsg := "Updated error"
	err.WithMessage(newMsg)

	if err.GetMessage() != newMsg {
		t.Errorf("Expected message %q, got %q", newMsg, err.GetMessage())
	}
}

func TestWithCode(t *testing.T) {
	msg := "Some error"
	initialCode := Forbidden
	err := New(msg, initialCode)

	newCode := Internal
	err.WithCode(newCode)

	if err.GetCode() != newCode {
		t.Errorf("Expected code %v, got %v", newCode, err.GetCode())
	}
}

func TestToHttpCode(t *testing.T) {
	tests := []struct {
		code     statusType
		expected int
	}{
		{Internal, http.StatusInternalServerError},
		{NotFound, http.StatusNotFound},
		{BadInput, http.StatusBadRequest},
		{Unauthorize, http.StatusUnauthorized},
		{Forbidden, http.StatusForbidden},
		{Conflict, http.StatusConflict},
		{99, http.StatusInternalServerError}, // Testing default case
	}

	for _, tt := range tests {
		err := New("", tt.code)
		assert.Equal(t, tt.expected, err.ToHttpCode())
	}
}

func TestError(t *testing.T) {
	msg := "This is an error message"
	err := New(msg, Internal)

	got := err.Error()
	if got != msg {
		t.Errorf("Error() = %q; want %q", got, msg)
	}
}

func TestToGRPCErr(t *testing.T) {
	msg := "This is a gRPC error"
	err := New(msg, Internal)

	grpcErr := err.ToGRPCErr()
	if grpcErr == nil {
		t.Fatal("Expected non-nil gRPC error")
	}

	// Check if the gRPC error message is as expected
	assert.ErrorIs(t, grpcErr, status.Errorf(codes.Internal, msg))

	// Check if the gRPC error code is as expected
	grpcStatus, ok := status.FromError(grpcErr)
	if !ok {
		t.Fatal("Expected gRPC status from error")
	}

	if grpcStatus.Code() != codes.Internal {
		t.Errorf("Expected gRPC code %v; got %v", codes.Internal, grpcStatus.Code())
	}
}

func TestToGRPCCode(t *testing.T) {
	tests := []struct {
		code     statusType
		expected codes.Code
	}{
		{code: Internal, expected: codes.Internal},
		{code: NotFound, expected: codes.NotFound},
		{code: BadInput, expected: codes.InvalidArgument},
		{code: Conflict, expected: codes.AlreadyExists},
		{code: 99, expected: codes.Internal},
	}

	for _, tc := range tests {
		err := New("", tc.code)
		assert.Equal(t, tc.expected, err.ToGRPCCode())
	}
}

func TestSlErr(t *testing.T) {
	msg := "test error for slog"
	code := Internal
	err := New(msg, code)

	assert.Equal(t, slog.String("error", msg), err.SlErr())
}
