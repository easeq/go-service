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

func createError(t *testing.T, id int, level int, errors ...*status.Status) *status.Status {
	msg := fmt.Sprintf(
		"Message ID: %d-%d, timestamp: %v",
		level,
		id,
		time.Now().UnixNano(),
	)

	code := getRandom(1, 16)
	invalidArg := WithCode(codes.Code(code))

	var errList []error
	for _, err := range errors {
		errList = append(errList, err.Err())
	}
	err := invalidArg(msg, errList...)

	require := require.New(t)
	require.NotEmpty(err)
	require.Equal(int(err.Code()), code)
	require.Equal(err.Message(), msg)
	require.Equal(len(err.Details()), len(errList))

	return err
}

func compareErrorDetails(t *testing.T, parent *status.Status, children ...*status.Status) {
	parentDetails := parent.Details()
	for i, child := range children {
		parentDetail := parentDetails[i].(*ErrorDetail)
		childAsDetail := FromStatusError(child.Err())

		require := require.New(t)
		require.Equal(parentDetail.Code, childAsDetail.Code)
		require.Equal(parentDetail.Status, childAsDetail.Status)
		require.Equal(parentDetail.Message, childAsDetail.Message)
		require.Equal(parentDetail.StackEntries, childAsDetail.StackEntries)
	}
}

func generateErrors(t *testing.T, n int, subN int, level int) ([]*status.Status, [][]*status.Status) {
	if n == 0 {
		n = getRandom(1, 10)
	}

	errList := make([]*status.Status, n)
	subErrList := make([][]*status.Status, n)
	for i := 0; i < n; i++ {
		if subN == -1 {
			subN = getRandom(0, 10)
		}

		var errors []*status.Status
		subErrList[i] = make([]*status.Status, subN)

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
