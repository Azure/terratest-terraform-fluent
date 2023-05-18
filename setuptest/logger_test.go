package setuptest

import (
	"bytes"
	"io"
	"math"
	"testing"

	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStreamLoggerShouldLogSth(t *testing.T) {
	buff := new(bytes.Buffer)
	l := NewStreamLogger(buff)
	log := "hello"
	l.Logf(t, "", log)
	assert.Contains(t, buff.String(), log)
}

func TestStreamLoggerCanPipeLogsToAnotherStreamLogger(t *testing.T) {
	log := "hello"
	srcBuff := bytes.NewBufferString(log)
	srcLogger := NewStreamLogger(srcBuff)
	destBuff := new(bytes.Buffer)
	destLogger := NewStreamLogger(destBuff)
	err := destLogger.PipeFrom(srcLogger)
	require.Nil(t, err)
	assert.Contains(t, destBuff.String(), log)
}

func TestStreamLoggerClose(t *testing.T) {
	log := "hello"
	srcBuff := bytes.NewBufferString(log)
	srcLogger := NewStreamLogger(srcBuff)

	destBuff := new(bytes.Buffer)
	dummyLogger := NewStreamLogger(destBuff)
	stub := gostub.Stub(&serializedLogger, dummyLogger)
	defer stub.Reset()

	err := srcLogger.Close()
	require.Nil(t, err)
	assert.Contains(t, destBuff.String(), log)
}

type testStream struct {
	content byte
	length  int
	finish  chan int
}

func newTestStream(content byte, length int) *testStream {
	return &testStream{
		content: content,
		length:  length,
		finish:  make(chan int),
	}
}

func (t *testStream) Read(p []byte) (n int, err error) {
	l := int(math.Min(float64(len(p)), float64(t.length)))
	for i := 0; i < l; i++ {
		p[i] = t.content
	}
	t.length -= l
	err = io.EOF
	if t.length != 0 {
		err = nil
	}
	return l, err
}

func (t *testStream) Write(p []byte) (n int, err error) {
	panic("Write is not expected to be called here")
}

var _ io.ReadWriter = new(testStream)

func TestStreamLogParallelLogShouldBePipeToStdoutSerialized(t *testing.T) {
	destBuff := new(bytes.Buffer)
	dummyLogger := NewStreamLogger(destBuff)
	stub := gostub.Stub(&serializedLogger, dummyLogger)
	defer stub.Reset()

	large1 := newTestStream(byte(1), 1024*1024*50)
	l1 := NewStreamLogger(large1)
	large2 := newTestStream(byte(2), 1024*1024*50)
	l2 := NewStreamLogger(large2)

	go func() {
		_ = l1.Close()
		large1.finish <- 1
	}()
	go func() {
		_ = l2.Close()
		large2.finish <- 1
	}()
	// Wait until two large streams are closed.
	for i := 0; i < 2; i++ {
		select {
		case <-large1.finish:
		case <-large2.finish:
		}
	}
	buff := destBuff.Bytes()
	gap := false
	for i := 1; i < len(buff); i++ {
		if buff[i] != buff[i-1] {
			if gap {
				t.Fatalf("not serialized output")
			}
			gap = true
		}
	}
}
