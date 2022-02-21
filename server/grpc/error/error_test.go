package error

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
)

func createError(t *testing.T, errors ...error) (error, string) {
	msg := "invalid arg message"
	invalidArg := WithCode(codes.InvalidArgument)

	err := invalidArg(msg, errors...)
	require.Error(t, err)
	require.Equal(t, err.Error(), msg)

	return err, msg
}

func TestWithCode(t *testing.T) {
	require := require.New(t)

	tests := []struct {
		name  string
		check func()
	}{
		{
			name: "EmptyErrList",
			check: func() {
				err, msg := createError(t)
				require.Error(err)
				require.Equal(err.Error(), msg)
			},
		},
		{
			name: "ErrListWithOneEntry",
			check: func() {
				oldErr, _ := createError(t)
				err, msg := createError(t, oldErr)
				require.Error(err)
				require.Equal(err.Error(), msg)
				// check details
			},
		},
		{
			name: "ErrListWithMultipleEntries",
			check: func() {
				oldErr, _ := createError(t)
				err, msg := createError(t, oldErr)
				require.Error(err)
				require.Equal(err.Error(), msg)
				// check details
			},
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			tc.check()
		})
	}
}
