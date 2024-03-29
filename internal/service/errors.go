package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	admin_movies_service "github.com/Falokut/admin_movies_service/pkg/admin_movies_service/v1/protos"
	"github.com/Falokut/grpc_errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInternal        = errors.New("internal error")
	ErrInvalidArgument = errors.New("invalid input data")
	ErrInvalidImage    = errors.New("invalid image")
	ErrAlreadyExists   = errors.New("already exists")
)

var errorCodes = map[error]codes.Code{
	ErrNotFound:        codes.NotFound,
	ErrInvalidArgument: codes.InvalidArgument,
	ErrInternal:        codes.Internal,
	ErrInvalidImage:    codes.InvalidArgument,
	ErrAlreadyExists:   codes.AlreadyExists,
}

type errorHandler struct {
	logger *logrus.Logger
}

func newErrorHandler(logger *logrus.Logger) errorHandler {
	return errorHandler{
		logger: logger,
	}
}

func IsContextError(msg string) bool {
	switch msg {
	case context.Canceled.Error(), context.DeadlineExceeded.Error():
		return true
	default:
		if strings.Contains(msg, "context canceled") {
			return true
		}

		if strings.Contains(msg, context.DeadlineExceeded.Error()) {
			return true
		}

		return false
	}
}

func (e *errorHandler) createErrorResponceWithSpan(span opentracing.Span, err error, developerMessage string) error {
	if err == nil {
		return nil
	}
	if IsContextError(developerMessage) {
		err = context.Canceled
		span.SetTag("grpc.status", codes.Canceled)
		ext.LogError(span, err)
	} else {
		span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
		ext.LogError(span, err)
	}

	return e.createErrorResponce(err, developerMessage)
}

func (e *errorHandler) createErrorResponce(err error, developerMessage string) error {
	if errors.Is(err, context.Canceled) || IsContextError(developerMessage) {
		err = status.Error(codes.Canceled, developerMessage)
		e.logger.Error(err)
		return err
	}

	var msg string
	if len(developerMessage) == 0 {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("%s. error: %v", developerMessage, err)
	}

	err = status.Error(grpc_errors.GetGrpcCode(err), msg)
	e.logger.Error(err)
	return err
}

func (e *errorHandler) createExtendedErrorResponceWithSpan(span opentracing.Span,
	err error, developerMessage, userMessage string) error {
	if err == nil {
		return nil
	}
	if IsContextError(developerMessage) {
		err = context.Canceled
		span.SetTag("grpc.status", codes.Canceled)
		ext.LogError(span, err)
	} else {
		span.SetTag("grpc.status", grpc_errors.GetGrpcCode(err))
		ext.LogError(span, err)
	}

	return e.createExtendedErrorResponce(err, developerMessage, userMessage)
}

func (e *errorHandler) createExtendedErrorResponce(err error, developerMessage, userMessage string) error {
	if errors.Is(err, context.Canceled) || IsContextError(developerMessage) {
		err = status.Error(codes.Canceled, developerMessage)
		e.logger.Error(err)
		return err
	}

	var msg string
	if developerMessage != "" {
		msg = fmt.Sprintf("%s. error: %v", developerMessage, err)
	} else {
		msg = err.Error()
	}

	extErr := status.New(grpc_errors.GetGrpcCode(err), msg)
	if len(userMessage) > 0 {
		extErr, _ = extErr.WithDetails(&admin_movies_service.UserErrorMessage{Message: userMessage})
		if extErr == nil {
			e.logger.Error(err)
			return err
		}
	}

	e.logger.Error(extErr)
	return extErr.Err()
}

func init() {
	grpc_errors.RegisterErrors(errorCodes)
}
