package contentbuffer

import (
	"errors"
	"io"
)

// implements io.ReadCloser
type ContentReader struct {
	buf      *ContentBuffer
	isOpened bool
	off      int
}

// Read reads the next len(p) bytes from the buffer or until the buffer
// is drained.  The return value n is the number of bytes read.  If the
// buffer has no data to return, err is io.EOF (unless len(p) is zero);
// otherwise it is nil.
// Safe for async operations.
func (cr *ContentReader) Read(p []byte) (n int, err error) {
	if !cr.isOpened {
		panic(errors.New("try read from closed content reader"))
	}
	if cr.off >= len(cr.buf.buf) {
		cr.off = 0
		if len(p) == 0 {
			return
		}
		return 0, io.EOF
	}
	n = copy(p, cr.buf.buf[cr.off:])
	cr.off += n
	return
}

func (cr *ContentReader) Close() error {
	cr.isOpened = false
	cr.buf.mutex.RUnlock()
	return nil
}
