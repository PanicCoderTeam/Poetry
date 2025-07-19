package utils

import (
	"poetry/src/pkg/trpc/codec/capi_error"

	"github.com/lonng/nano/session"
)

func NewErrorAndResponse(s *session.Session, code capi_error.ErrorCode, msg string, err error, a ...interface{}) *capi_error.CapiError {
	err1 := capi_error.NewError(code, msg, err, a...)
	s.Response(err1)
	return err1
}

// NewError returns a wrapError instance
func NewErrAndResponse(s *session.Session, code capi_error.ErrorCode, msg string, a ...interface{}) *capi_error.CapiError {
	err1 := capi_error.NewErr(code, msg, a...)
	s.Response(err1)
	return err1
}
