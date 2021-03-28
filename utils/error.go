package utils

import (
	"database/sql"
	"errors"

	"google.golang.org/grpc/codes"
)

// GetErrorCode returns the error code to use
func GetErrorCode(err error) codes.Code {
	if errors.Is(err, sql.ErrNoRows) {
		return codes.NotFound
	}

	return codes.Internal
}
