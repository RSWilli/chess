package uci

import (
	"fmt"
	"io"
)

type loggingWriter struct{}

// Write implements io.Writer.
func (l loggingWriter) Write(p []byte) (n int, err error) {
	fmt.Println(string(p))

	return len(p), nil
}

var _ io.Writer = loggingWriter{}
