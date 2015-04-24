package main

import (
	"encoding/binary"
	"io"
	"math"
)

// Return a pointer to a MultiplexedReader for demuxing io.Reader s.
func Demuxer(s io.Reader) *MultiplexedReader {
	return &MultiplexedReader{s: s}
}

// MultiplexedReader wraps the stream from docker logs and demuxes the stream.
//
// Docker log streams are multiplexed, with 8-byte headers describing the
// source stream (stdin, stdout, stderr) and the length of the forthcoming
// payload. This reader demuxes such a stream, but discards the source stream
// information, munging everything into a single stream.
type MultiplexedReader struct {
	s io.Reader
	b []byte
}

// Read demuxed bytes from the stream info the byte slice.
//
// Returns the actual number of bytes placed into the buffer, not the total
// number of bytes read in processing the underlying stream.
func (r *MultiplexedReader) Read(b []byte) (int, error) {
	if len(r.b) < len(b) {
		if _, err := r.fill(); err != nil {
			return 0, err
		}
	}

	n := int(math.Min(
		float64(len(r.b)),
		float64(len(b)),
	))

	for i, v := range r.b[:n] {
		b[i] = v
	}

	r.b = r.b[n:]

	return n, nil
}

// Fills the internal buffer.
func (r *MultiplexedReader) fill() (int, error) {
	var (
		got int
		err error
	)

	header := make([]byte, 8)

	if _, err = r.s.Read(header); err == nil {
		read := 0
		size := binary.BigEndian.Uint32(header[4:])
		buf := make([]byte, 0)

		for size > 0 {
			payload := make([]byte, size)

			if read, err = r.s.Read(payload); read > 0 {
				buf = append(buf, payload...)
				size -= uint32(read)
			}
		}

		r.b = append(r.b, buf...)
		got = len(buf)
	}

	return got, err
}
