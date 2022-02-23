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
			errDetails := Convert(errors[i])
			st, _ = st.WithDetails(errDetails)
		}

		return st.Err()
	}
}

// Convert error to ErrorDetail
func Convert(err error) *ErrorDetail {
	st := status.Convert(err)

	var entries []string
	for _, detail := range st.Details() {
		errDetail, _ := detail.(*ErrorDetail)
		entries = append(entries, errDetail.Message)
		entries = append(entries, errDetail.StackEntries...)
	}

	return &ErrorDetail{
		Code:         int32(st.Code()),
		Status:       st.Code().String(),
		Message:      st.Message(),
		StackEntries: entries,
	}
}
