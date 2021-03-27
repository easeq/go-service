package utils

import (
	"database/sql"

	"google.golang.org/grpc/codes"
)

// GetErrorCode returns the error code to use
func GetErrorCode(err error) codes.Code {
	if err == sql.ErrNoRows {
		return codes.NotFound
	}

	return codes.Internal
}
