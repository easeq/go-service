package error

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func getRandom(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func createError(t *testing.T, id int, level int, errors ...error) error {
	msg := fmt.Sprintf(
		"Message ID: %d-%d, timestamp: %v",
		level,
		id,
		time.Now().UnixNano(),
	)

	code := getRandom(1, 16)
	errType := WithCode(codes.Code(code))

	var errList []error
	for _, err := range errors {
		errList = append(errList, err)
	}
	err := errType(msg, errList...)

	require := require.New(t)
	require.NotEmpty(err)

	st, _ := status.FromError(err)

	require.Equal(int(st.Code()), code)
	require.Equal(st.Message(), msg)
	require.Equal(len(st.Details()), len(errList))

	return err
}

func compareErrorDetails(t *testing.T, parent error, children ...error) {
	parentSt, _ := status.FromError(parent)
	parentDetails := parentSt.Details()
	for i, child := range children {
		parentDetail := parentDetails[i].(*ErrorDetail)
		childAsDetail := FromStatusError(child)

		require := require.New(t)
		require.Equal(parentDetail.Code, childAsDetail.Code)
		require.Equal(parentDetail.Status, childAsDetail.Status)
		require.Equal(parentDetail.Message, childAsDetail.Message)
		require.Equal(parentDetail.StackEntries, childAsDetail.StackEntries)
	}
}

func generateErrors(t *testing.T, n int, subN int, level int) ([]error, [][]error) {
	if n == 0 {
		n = getRandom(1, 10)
	}

	errList := make([]error, n)
	subErrList := make([][]error, n)
	for i := 0; i < n; i++ {
		if subN == -1 {
			subN = getRandom(0, 10)
		}

		var errors []error
		subErrList[i] = make([]error, subN)

		if subN > 0 {
			errors, _ = generateErrors(t, subN, 0, level+1)
			subErrList[i] = errors
		}

		errList[i] = createError(t, i, level, errors...)
	}

	return errList, subErrList
}

func TestWithCode(t *testing.T) {
	tests := []struct {
		name  string
		check func()
	}{
		{
			name: "EmptyErrList",
			check: func() {
				createError(t, 1, 1)
			},
		},
		{
			name: "ErrWith_1_ErrDetail",
			check: func() {
				errors, subErrList := generateErrors(t, 1, 1, 1)
				compareErrorDetails(t, errors[0], subErrList[0]...)
			},
		},
		{
			name: "ErrWith_2_ErrDetails",
			check: func() {
				errors, subErrList := generateErrors(t, 1, 2, 1)
				compareErrorDetails(t, errors[0], subErrList[0]...)
			},
		},
		{
			name: "ErrWith_{0..10}_ErrDetails",
			check: func() {
				errors, subErrList := generateErrors(t, 1, -1, 1)
				compareErrorDetails(t, errors[0], subErrList[0]...)
			},
		},
		{
			name: "ErrWith_1_ErrDetails_With_2_StackEntries",
			check: func() {
				errors, _ := generateErrors(t, 1, 2, 2)
				err := createError(t, 100, 1, errors...)
				compareErrorDetails(t, err, errors...)
			},
		},
		{
			name: "ErrWith_1_ErrDetails_With_{0..10}_StackEntries",
			check: func() {
				errors, _ := generateErrors(t, 1, -1, 2)
				err := createError(t, 100, 1, errors...)
				compareErrorDetails(t, err, errors...)
			},
		},
		{
			name: "ErrWith_2_ErrDetails_With_3_StackEntriesEach",
			check: func() {
				errors, _ := generateErrors(t, 2, 3, 2)
				err := createError(t, 100, 1, errors...)
				compareErrorDetails(t, err, errors...)
			},
		},
		{
			name: "ErrWith_2_ErrDetails_With_{0..10}_StackEntries",
			check: func() {
				errors, _ := generateErrors(t, 2, -1, 2)
				err := createError(t, 100, 1, errors...)
				compareErrorDetails(t, err, errors...)
			},
		},
		{
			name: "ErrWith_{0..10}_ErrDetails_With_{0..10}_StackEntries",
			check: func() {
				errors, _ := generateErrors(t, 0, -1, 2)
				err := createError(t, 100, 1, errors...)
				compareErrorDetails(t, err, errors...)
			},
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.check()
		})
	}
}

func errorWithCode(t *testing.T, fn GrpcError, code codes.Code) {
	require := require.New(t)

	msg := fmt.Sprintf("%s message", code.String())
	detail := fmt.Errorf("%s detail", code.String())
	err := fn(msg, detail)

	require.Error(err)

	st := status.Convert(err)
	require.Equal(st.Code(), code)
	require.Equal(st.Message(), msg)
}

func TestAborted(t *testing.T) {
	errorWithCode(t, Aborted, codes.Aborted)
}

func TestAlreadyExists(t *testing.T) {
	errorWithCode(t, AlreadyExists, codes.AlreadyExists)
}

func TestCanceled(t *testing.T) {
	errorWithCode(t, Canceled, codes.Canceled)
}

func TestDataLoss(t *testing.T) {
	errorWithCode(t, DataLoss, codes.DataLoss)
}

func TestDeadlineExceeded(t *testing.T) {
	errorWithCode(t, DeadlineExceeded, codes.DeadlineExceeded)
}

func TestFailedPrecondition(t *testing.T) {
	errorWithCode(t, FailedPrecondition, codes.FailedPrecondition)
}

func TestInternal(t *testing.T) {
	errorWithCode(t, Internal, codes.Internal)
}

func TestInvalidArgument(t *testing.T) {
	errorWithCode(t, InvalidArgument, codes.InvalidArgument)
}

func TestNotFound(t *testing.T) {
	errorWithCode(t, NotFound, codes.NotFound)
}

func TestOutOfRange(t *testing.T) {
	errorWithCode(t, OutOfRange, codes.OutOfRange)
}

func TestPermissionDenied(t *testing.T) {
	errorWithCode(t, PermissionDenied, codes.PermissionDenied)
}

func TestResourceExhausted(t *testing.T) {
	errorWithCode(t, ResourceExhausted, codes.ResourceExhausted)
}

func TestUnavailable(t *testing.T) {
	errorWithCode(t, Unavailable, codes.Unavailable)
}

func TestUnauthenticated(t *testing.T) {
	errorWithCode(t, Unauthenticated, codes.Unauthenticated)
}

func TestUnimplemented(t *testing.T) {
	errorWithCode(t, Unimplemented, codes.Unimplemented)
}

func TestUnknown(t *testing.T) {
	errorWithCode(t, Unknown, codes.Unknown)
}
