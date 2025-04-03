package setuptest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

var serializedLogger = func() *StreamLogger {
	l := NewStreamLogger(os.Stdout)
	l.outputProgress = false
	return l
}()

// A StreamLogger is a logger that writes to a stream, such as stdout, a file or a memory buffer. It also keeps track of the number of
// log lines written, and can output a progress message every 50 lines.
type StreamLogger struct {
	stream         io.ReadWriter // The stream to write to
	mu             *sync.Mutex   // A mutex to ensure only one thread writes to the stream at a time
	logCount       int           // The number of log lines written
	outputProgress bool          // Whether or not to output a progress message every 50 lines
}

// NewMemoryLogger creates a new StreamLogger that writes to an in-memory buffer. This is useful for capturing logs in tests.
func NewMemoryLogger() *StreamLogger {
	buff := new(bytes.Buffer)
	return NewStreamLogger(buff)
}

// NewStreamLogger creates a new StreamLogger that writes to the given stream.
func NewStreamLogger(stream io.ReadWriter) *StreamLogger {
	return &StreamLogger{
		stream:         stream,
		mu:             new(sync.Mutex),
		outputProgress: false,
	}
}

// Logf logs the given arguments to the given writer, along with a prefix of the test name.
func (s *StreamLogger) Logf(t testing.TestingT, format string, args ...interface{}) {
	// Sprintf removed as we don't want the prefixes to the log lines
	// log := fmt.Sprintf(format, args...)
	doLog(t, s.stream, args...)
	s.logCount++
	if s.outputProgress && s.logCount%50 == 0 {
		logger.Log(t, fmt.Sprintf("logging sample: %s", args))
	}
}

// The PipeFrom function is a method of the StreamLogger struct. It takes a pointer to another StreamLogger object as its input parameter and returns an error.
// Inside the function, a mutex lock is acquired to ensure that the function is thread-safe.
// The PipeFrom function is useful when you want to redirect the output of one logger to another logger.
// For example, if you have a logger that writes to the console and another logger that writes to a file, you can use PipeFrom to redirect the console logger's output to the file logger:
func (s *StreamLogger) PipeFrom(srcLogger *StreamLogger) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := io.Copy(s.stream, srcLogger.stream)
	return err
}

func (s *StreamLogger) Close() error {
	defer func() {
		c, ok := s.stream.(io.Closer)
		if ok {
			_ = c.Close()
		}
	}()
	return serializedLogger.PipeFrom(s)
}

// doLog logs the given arguments to the given writer, along with a prefix of the test name.
func doLog(t testing.TestingT, writer io.Writer, args ...interface{}) {
	// date := time.Now()
	prefix := fmt.Sprintf("%s:", t.Name())
	allArgs := append([]interface{}{prefix}, args...)
	fmt.Fprintln(writer, allArgs...)
}
