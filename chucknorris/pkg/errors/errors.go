package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// NoType error
	NoType = "NoType"
	// BadRequest error
	BadRequest = "BadRequest"
	// NotFound error
	NotFound = "NotFound"
)

// ErrorType is the type of an error
type ErrorType string

// customError will be returned by the applications's APIs if a status code other than 200 is returned
type customError struct {
	errorType     ErrorType
	code          int
	originalError error
	context       errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// Error returns the description of an customError
func (err customError) Error() string {
	return err.originalError.Error()
}

// New creates a new customError
func (errType ErrorType) New(msg string, code int) error {
	return customError{errorType: errType,
		code:          code,
		originalError: errors.New(msg),
	}
}

// Newf creates a new customError with formatted message
func (errType ErrorType) Newf(msg string, code int, args ...interface{}) error {
	return customError{errorType: errType, originalError: fmt.Errorf(msg, args...)}
}

// Wrap creates a new wrapped error
func (errType ErrorType) Wrap(err error, msg string) error {
	return errType.Wrap(err, msg)
}

// Wrapf creates a new wrapped error with formatted message
func (errType ErrorType) Wrapf(err error, code int, msg string, args ...interface{}) error {
	return customError{errorType: errType, originalError: errors.Wrapf(err, msg, args...)}
}

// NewErr creates a no type error
func NewErr(msg string) error {
	return customError{errorType: NoType, originalError: errors.New(msg)}
}

// NewfErr creates a no type error
func NewfErr(msg string, args ...interface{}) error {
	return customError{errorType: NoType, originalError: errors.New(fmt.Sprintf(msg, args...))}
}

// Wrap returns an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Wrapf returns an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	wrappedError := errors.Wrapf(err, msg, args...)
	if customErr, ok := err.(customError); ok {
		return customError{
			errorType:     customErr.errorType,
			originalError: wrappedError,
			context:       customErr.context,
		}
	}
	return customError{errorType: NoType, originalError: wrappedError}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(customError); ok {
		return customError{errorType: customErr.errorType, originalError: customErr.originalError, context: context}
	}
	return customError{errorType: NoType, originalError: err, context: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(customError); ok || customErr.context != emptyContext {
		return map[string]string{"field": customErr.context.Field, "message": customErr.context.Message}
	}
	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}
	return NoType
}
