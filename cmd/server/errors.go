package main

import (
	"errors"
	"fmt"
)

const (
	msg string = "method not implemented"
)

// ErrNotImplemented is an error used to represent a method not implemented
var ErrNotImplemented = errors.New(msg)

// ErrPrometheusClient
type ErrPrometheusClient struct {
	Msg string
	Err error
}

// NewErrPrometheusClient is a function that creates a new ErrPrometheusClient
func NewErrPrometheusClient(msg string, err error) *ErrPrometheusClient {
	return &ErrPrometheusClient{
		Msg: msg,
		Err: err,
	}
}

// Error is a method that implements the error interface for ErrPrometheusClient
func (e *ErrPrometheusClient) Error() string {
	// If there's a wrapped error, include its message
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", e.Msg, e.Err.Error())
	}

	return e.Msg
}

// GoString is a method that converts an ErrPrometheusClient to its equivalent Go syntax
func (e *ErrPrometheusClient) GoString() string {
	// Handle pointer field
	if e.Err != nil {
		return fmt.Sprintf("&ErrPrometheusClient{Msg: %q, Err: %v}", e.Msg, e.Err)
	}

	return fmt.Sprintf("&ErrPrometheusClient{Msg: %q}", e.Msg)
}

// Unwrap is a method that unwraps any wrapped errors
func (e *ErrPrometheusClient) Unwrap() error {
	return e.Err
}
