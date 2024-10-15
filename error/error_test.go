package e

import (
	"errors"
	"log/slog"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	msg := "An error occurred"
	testErr := errors.New("some error")
	code := NotFound

	err := New(msg, code, testErr)

	if err.GetMessage() != msg {
		t.Errorf("Expected message %q, got %q", msg, err.GetMessage())
	}

	if errors.Is(err.GetError(), err) {
		t.Errorf("Expected message %s, got %s", err.Error(), err.GetError().Error())
	}

	if err.GetCode() != code {
		t.Errorf("Expected code %v, got %v", code, err.GetCode())
	}
}

func TestWithMessage(t *testing.T) {
	initialMsg := "Initial error"
	testErr := errors.New("some error")
	code := Internal
	err := New(initialMsg, code, testErr)

	newMsg := "Updated error"
	err.WithMessage(newMsg)

	if err.GetMessage() != newMsg {
		t.Errorf("Expected message %q, got %q", newMsg, err.GetMessage())
	}
}

func TestWithCode(t *testing.T) {
	msg := "Some msg"
	testErr := errors.New("Some error")
	initialCode := Forbidden
	err := New(msg, initialCode, testErr)

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
		err := New("", tt.code, nil)
		assert.Equal(t, tt.expected, err.ToHttpCode())
	}
}

func TestError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    []error
		msg    string
		expect string
	}{{
		name:   "AllFields",
		err:    []error{errors.New("some error")},
		msg:    "some message",
		expect: "some message: some error",
	},
		{
			name:   "Only Message",
			err:    []error{},
			msg:    "some message",
			expect: "some message",
		},
		{
			name:   "Only error",
			err:    []error{errors.New("some error")},
			msg:    "",
			expect: "some error",
		},
		{
			name:   "Nil fileds",
			err:    []error{},
			msg:    "",
			expect: "nil",
		}}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := New(tc.msg, Internal, tc.err...)

			assert.Equal(t, tc.expect, err.Error())
		})
	}

}

func TestToGRPCErr(t *testing.T) {
	msg := "This is a gRPC error message"
	testErr := errors.New("some grpc error")
	err := New(msg, Internal, testErr)

	grpcErr := err.ToGRPCErr()
	if grpcErr == nil {
		t.Fatal("Expected non-nil gRPC error")
	}

	// Check if the gRPC error message is as expected
	assert.ErrorIs(t, grpcErr, status.Error(codes.Internal, msg))

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
		err := New("", tc.code, nil)
		assert.Equal(t, tc.expected, err.ToGRPCCode())
	}
}

func TestSlErr(t *testing.T) {
	msg := "test error for slog"
	testErr := errors.New("some error")
	code := Internal
	err := New(msg, code, testErr)

	assert.Equal(t, slog.String("error", err.Error()), err.SlErr())
}

func TestE(t *testing.T) {
	t.Parallel()

	testErr := errors.New("some error")
	tests := []struct {
		name string
		err  error
		want Error
	}{
		{
			name: "Nil error",
			err:  nil,
			want: nil,
		},
		{
			name: "Not nil error",
			err:  testErr,
			want: New("", Internal, testErr),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, E(tt.err))
		})
	}
}
