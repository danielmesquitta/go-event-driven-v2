package errs

type errInvalidFormat struct{}

func (e errInvalidFormat) New(opts ...Option) *Error {
	return NewError(CodeBadRequest, "invalid format", opts...)
}

var ErrInvalidFormat = errInvalidFormat{}
var _ ErrorBuilder = ErrInvalidFormat
