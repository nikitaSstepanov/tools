package e

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/nikitaSstepanov/tools/ctx"
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

func TestGetMessage(t *testing.T) {
	msg := "Error"
	testErr := errors.New("some error")
	code := Internal
	err := New(msg, code, testErr)

	assert.Equal(t, msg, err.GetMessage())
}

func TestGetError(t *testing.T) {
	msg := "error"
	testErr1 := errors.New("some error 1")
	testErr2 := errors.New("some error 2")
	code := Internal

	joinedErr := errors.Join(testErr1, testErr2)
	err := New(msg, code, testErr1, testErr2)

	assert.Equal(t, joinedErr, err.GetError())
}

func TestGetTag(t *testing.T) {
	msg := "error"
	code := Internal
	err := New(msg, code)

	key := "key"
	value := "value"

	err = err.WithTag(key, value)

	assert.Equal(t, value, err.GetTag(key))
}

func TestGetCode(t *testing.T) {
	msg := "error"
	code := Internal
	err := New(msg, code)

	assert.Equal(t, code, err.GetCode())
}

func TestLog(t *testing.T) {
	var buff bytes.Buffer

	mockHandler := slog.NewJSONHandler(&buff, &slog.HandlerOptions{AddSource: false, Level: slog.LevelInfo})
	mockLog := slog.New(mockHandler)

	slog.SetDefault(mockLog)

	msg := "error"
	code := Internal
	err := New(msg, code)

	err.Log()

	expectedLevel := `"level":"ERROR"`
	expectedMsg := `"msg":""`
	expectedAttr := `"error":"error"`

	parts := strings.Split(buff.String()[:83], ",")
	fmt.Println(parts)

	assert.Equal(t, expectedLevel, parts[1])
	assert.Equal(t, expectedMsg, parts[2])
	assert.Equal(t, expectedAttr, parts[3])
}

func TestWithMessage(t *testing.T) {
	initialMsg := "Initial error"
	testErr := errors.New("some error")
	code := Internal
	err := New(initialMsg, code, testErr)

	newMsg := "Updated error"
	err = err.WithMessage(newMsg)

	if err.GetMessage() != newMsg {
		t.Errorf("Expected message %q, got %q", newMsg, err.GetMessage())
	}
}

func TestWithErr(t *testing.T) {
	msg := "error"
	code := Internal
	err := New(msg, code)

	testErr := errors.New("some error")
	joined := errors.Join(testErr)

	err = err.WithErr(testErr)

	assert.Equal(t, joined, err.GetError())
}

func TestWithTag(t *testing.T) {
	msg := "error"
	code := Internal
	err := New(msg, code)

	key := "key"
	value := "value"

	err = err.WithTag(key, value)

	assert.Equal(t, value, err.GetTag(key))
}

func TestWithCtx(t *testing.T) {
	c := ctx.New(slog.Default())

	key := "key"
	value := "value"

	c.AddValue(key, value, true)

	msg := "error"
	code := Internal

	err := New(msg, code)

	err = err.WithCtx(c)

	assert.Equal(t, value, err.GetTag(key))
}

func TestWithCode(t *testing.T) {
	msg := "Some msg"
	testErr := errors.New("Some error")
	initialCode := Forbidden
	err := New(msg, initialCode, testErr)

	newCode := Internal
	err = err.WithCode(newCode)

	if err.GetCode() != newCode {
		t.Errorf("Expected code %v, got %v", newCode, err.GetCode())
	}
}

func TestToJson(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  Error
		want JsonError
	}{
		{
			name: "Default",
			err:  New("some error", Internal),
			want: JsonError{Error: "some error"},
		},
		{
			name: "With err",
			err:  New("some error", Internal, errors.New("invalid data")),
			want: JsonError{Error: "some error"},
		},
		{
			name: "Empty message",
			err:  New("", Internal),
			want: JsonError{Error: ""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.err.ToJson())
		})
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
