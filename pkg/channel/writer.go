package channel

import (
	"io"
	"sync"

	"github.com/pkg/errors"
)

// WriteCloser is an io.WriteCloser that that proxies Write() calls to a channel
// The []byte buffer of the Write() is queued on the channel as one message.
type WriteCloser interface {
	io.WriteCloser
	Chan() <-chan []byte
}

type writeCloser struct {
	ch  chan []byte
	mux sync.Mutex
}

// NewWriter initializes a new channel writer
func NewWriter(c chan []byte) WriteCloser {
	return &writeCloser{
		ch: c,
	}
}

// Chan returns the R/O channel behind WriteCloser
func (w *writeCloser) Chan() <-chan []byte {
	return w.ch
}

// Write method for WriteCloser
func (w *writeCloser) Write(b []byte) (int, error) {
	if w == nil || w.ch == nil {
		return 0, errors.New("use channel.NewWriter() to initialize a WriteCloser")
	}

	w.mux.Lock()
	buf := make([]byte, len(b))
	copy(buf, b)
	w.ch <- buf
	w.mux.Unlock()

	return len(b), nil
}

// Close method for WriteCloser
func (w *writeCloser) Close() error {
	close(w.ch)
	return nil
}