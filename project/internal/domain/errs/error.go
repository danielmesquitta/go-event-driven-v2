package errs

import (
	"encoding/json"
	"fmt"
)

type ErrorBuilder interface {
	New(opts ...Option) *Error
}

type Error struct {
	code     Code
	message  string
	metadata map[string]any
}

type Code string

const (
	CodeInternal     Code = "internal"
	CodeNotFound     Code = "not_found"
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeBadRequest   Code = "bad_request"
)

type Option func(*Error)

func WithMetadata(key string, value any) Option {
	return func(e *Error) {
		if e.metadata == nil {
			e.metadata = make(map[string]any)
		}
		e.metadata[key] = value
	}
}

func NewError(
	code Code,
	message string,
	opts ...Option,
) *Error {
	e := &Error{
		code:    code,
		message: message,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

func (e *Error) Error() string {
	jsonBytes, err := json.Marshal(e.metadata)
	if err != nil {
		return e.message
	}

	return fmt.Sprintf("%s: %s", e.message, string(jsonBytes))
}

func (e *Error) Code() Code {
	return e.code
}

func (e *Error) Set(key string, value any) {
	if e.metadata == nil {
		e.metadata = make(map[string]any)
	}
	e.metadata[key] = value
}
