package error

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GrpcError returns a status error from the given message and list of errors
type GrpcError func(msg string, errors ...error) error

var (
	Aborted            = WithCode(codes.Aborted)
	AlreadyExists      = WithCode(codes.AlreadyExists)
	Canceled           = WithCode(codes.Canceled)
	DataLoss           = WithCode(codes.DataLoss)
	DeadlineExceeded   = WithCode(codes.DeadlineExceeded)
	FailedPrecondition = WithCode(codes.FailedPrecondition)
	Internal           = WithCode(codes.Internal)
	InvalidArgument    = WithCode(codes.InvalidArgument)
	NotFound           = WithCode(codes.NotFound)
	OutOfRange         = WithCode(codes.OutOfRange)
	PermissionDenied   = WithCode(codes.PermissionDenied)
	ResourceExhausted  = WithCode(codes.ResourceExhausted)
	Unavailable        = WithCode(codes.Unavailable)
	Unauthenticated    = WithCode(codes.Unauthenticated)
	Unimplemented      = WithCode(codes.Unimplemented)
	Unknown            = WithCode(codes.Unknown)
)

// WithCode returns a prepared function with the respective code
func WithCode(code codes.Code) GrpcError {
	return func(msg string, errors ...error) error {
		st := status.New(code, msg)
		for i := range errors {
			st, _ = st.WithDetails(FromStatusError(errors[i]))
		}

		return st.Err()
	}
}

// FromStatusError converts status error to ErrorDetail
func FromStatusError(err error) *ErrorDetail {
	st, _ := status.FromError(err)
	return &ErrorDetail{
		Code:    int32(st.Code()),
		Status:  st.Code().String(),
		Message: st.Message(),
	}
}
