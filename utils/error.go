package utils

import (
	"database/sql"
	"errors"
	"sync"

	"google.golang.org/grpc/codes"
)

// GetErrorCode returns the error code to use
func GetErrorCode(err error) codes.Code {
	if errors.Is(err, sql.ErrNoRows) {
		return codes.NotFound
	}

	return codes.Internal
}

// MergeErrors merges multiple channels of errors.
func MergeErrors(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error, len(cs))

	// Handle individual error channels
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	for _, c := range cs {
		if c == nil {
			continue
		}

		wg.Add(1)
		go output(c)
	}

	// Wait and close "out" channel once done
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// WaitForError waits for results from all error channels
// Returns on first non-nil error or returns nil after all channels are closed
func WaitForError(errs ...<-chan error) error {
	errc := MergeErrors(errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}
