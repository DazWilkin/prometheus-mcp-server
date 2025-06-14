package errors

import (
	"errors"
	"fmt"
)

const (
	msg string = "method not implemented"
)

// ErrNotImplemented is an error used to represent a method not implemented
var ErrNotImplemented = errors.New(msg)

// ErrConfig is a type that represents errors return by Config
type ErrConfig struct {
	Msg string
	Err error
}

// NewErrConfig is a function that creates a new ErrConfig
func NewErrConfig(msg string, err error) *ErrConfig {
	return &ErrConfig{
		Msg: msg,
		Err: err,
	}
}

// Error is a method that implements the error interface for ErrConfig
func (e *ErrConfig) Error() string {
	// If there's a wrapped error, include its message
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
	}

	return e.Msg
}

// GoString is a method that converts an ErrClient to its equivalent Go syntax
func (e *ErrConfig) GoString() string {
	// Handle pointer field
	if e.Err != nil {
		return fmt.Sprintf("&ErrConfig{Msg: %q, Err: %v}", e.Msg, e.Err)
	}

	return fmt.Sprintf("&ErrConfig{Msg: %q}", e.Msg)
}

// Unwrap is a method that unwraps any wrapped errors
func (e *ErrConfig) Unwrap() error {
	return e.Err
}

// ErrToolHandler is a type that represents errors returned by Client
type ErrToolHandler struct {
	Msg string
	Err error
}

// NewErrToolHandler is a function that creates a new ErrToolHandler
func NewErrToolHandler(msg string, err error) *ErrToolHandler {
	return &ErrToolHandler{
		Msg: msg,
		Err: err,
	}
}

// Error is a method that implements the error interface for ErrClient
func (e *ErrToolHandler) Error() string {
	// If there's a wrapped error, include its message
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
	}

	return e.Msg
}

// GoString is a method that converts an ErrClient to its equivalent Go syntax
func (e *ErrToolHandler) GoString() string {
	// Handle pointer field
	if e.Err != nil {
		return fmt.Sprintf("&ErrClient{Msg: %q, Err: %v}", e.Msg, e.Err)
	}

	return fmt.Sprintf("&ErrClient{Msg: %q}", e.Msg)
}

// Unwrap is a method that unwraps any wrapped errors
func (e *ErrToolHandler) Unwrap() error {
	return e.Err
}
