package errs

type errNotFound struct{}

func (e errNotFound) New(opts ...Option) *Error {
	return NewError(CodeNotFound, "resource not found", opts...)
}

var ErrNotFound = errNotFound{}
var _ ErrorBuilder = ErrNotFound
