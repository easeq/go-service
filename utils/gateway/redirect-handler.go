package gateway

import (
	"context"
	"net/http"
	"strconv"

	"google.golang.org/protobuf/runtime/protoiface"
)

// RedirectHandler gRPC metadata and redirect if redirection headers are set
func RedirectHandler(ctx context.Context, w http.ResponseWriter, resp protoiface.MessageV1) error {
	headers := w.Header()
	if location, ok := headers["Grpc-Metadata-Location"]; ok {
		w.Header().Set("Location", location[0])

		if code, ok := headers["Grpc-Metadata-Code"]; ok {
			codeInt, err := strconv.Atoi(code[0])
			if err != nil {
				return err
			}

			w.WriteHeader(codeInt)
		} else {
			w.WriteHeader(http.StatusFound)
		}
	}

	return nil
}
