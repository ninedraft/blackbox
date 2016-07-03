package contentbuffer

import (
	"bytes"
	"io"
	"sync"
)

var (
	ErrTooLarge = bytes.ErrTooLarge
)

const (
	MinRead = 512
)

// makeSlice allocates a slice of size n. If the allocation fails, it panics
// with ErrTooLarge.
func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	return make([]byte, n)
}

type ContentBuffer struct {
	buf       []byte
	mutex     sync.RWMutex
	off       int
	bootstrap [64]byte
}

func (cb *ContentBuffer) Len() int { return len(cb.buf) }

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (cb *ContentBuffer) grow(n int) int {
	m := cb.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && cb.off != 0 {
		// mutex locks here!
		cb.truncate(0)
	}
	if len(cb.buf)+n > cap(cb.buf) {
		var buf []byte
		if cb.buf == nil && n <= len(cb.bootstrap) {
			buf = cb.bootstrap[0:]
		} else if m+n <= cap(cb.buf)/2 {
			// We can slide things down instead of allocating a new
			// slice. We only need m+n <= cap(b.buf) to slide, but
			// we instead let capacity get twice as large so we
			// don't spend all our time copying.
			copy(cb.buf[:], cb.buf[cb.off:])
			buf = cb.buf[:m]
		} else {
			// not enough space anywhere
			buf = makeSlice(2*cap(cb.buf) + n)
			copy(buf, cb.buf[cb.off:])
		}
		cb.buf = buf
		cb.off = 0
	}
	cb.buf = cb.buf[0 : cb.off+m+n]
	return cb.off + m
}

func (cb *ContentBuffer) truncate(n int) {
	switch {
	case n < 0 || n > cb.Len():
		panic("bytes.Buffer: truncation out of range")
	case n == 0:
		// Reuse buffer space.
		cb.off = 0
	}
	cb.buf = cb.buf[0 : cb.off+n]
}

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
// Locks the mutex
func (cb *ContentBuffer) Truncate(n int) {
	cb.mutex.Lock()
	cb.truncate(n)
	cb.mutex.Unlock()
}

// Write appends the contents of p to the buffer, growing the buffer as
// needed. The return value n is the length of p; err is always nil. If the
// buffer becomes too large, Write will panic with ErrTooLarge.
// Locks the mutex
func (cb *ContentBuffer) Write(p []byte) (n int, err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	m := cb.grow(len(p))
	BytesCopied := copy(cb.buf[m:], p)
	return BytesCopied, nil
}

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with ErrTooLarge.
func (cb *ContentBuffer) ReadFrom(r io.Reader) (n int64, err error) {
	// If buffer is empty, reset to recover space.
	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	if cb.off >= len(cb.buf) {
		cb.truncate(0)
	}
	for {
		if free := cap(cb.buf) - len(cb.buf); free < MinRead {
			// not enough space at end
			newBuf := cb.buf
			if cb.off+free < MinRead {
				// not enough space using beginning of buffer;
				// double buffer capacity
				newBuf = makeSlice(2*cap(cb.buf) + MinRead)
			}
			copy(newBuf, cb.buf[cb.off:])
			cb.buf = newBuf[:len(cb.buf)-cb.off]
			cb.off = 0
		}
		m, e := r.Read(cb.buf[len(cb.buf):cap(cb.buf)])
		cb.buf = cb.buf[0 : len(cb.buf)+m]
		n += int64(m)
		if e == io.EOF {
			break
		}
		if e != nil {
			return n, e
		}
	}
	return n, nil // err is EOF, so return nil explicitly
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (cb *ContentBuffer) WriteTo(w io.Writer) (n int64, err error) {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	if cb.off < len(cb.buf) {
		nBytes := cb.Len()
		m, e := w.Write(cb.buf[cb.off:])
		if m > nBytes {
			panic("bytes.Buffer.WriteTo: invalid Write count")
		}
		cb.off += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	return
}
