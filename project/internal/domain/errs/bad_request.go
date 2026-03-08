package errs

type errBadRequest struct{}

func (e errBadRequest) New(opts ...Option) *Error {
	return NewError(CodeBadRequest, "bad request", opts...)
}

var ErrBadRequest = errBadRequest{}
var _ ErrorBuilder = ErrBadRequest
