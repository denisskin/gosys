package gosys

import "io"

var DevNull *emptyStream

//----------------------------------------------
type emptyStream struct {
}

func (*emptyStream) Read(buf []byte) (int, error) {
	return len(buf), nil
}

func (*emptyStream) Write(buf []byte) (int, error) {
	return len(buf), nil
}

func (*emptyStream) Close() error {
	return nil
}

func NewProgressWriter(w io.Writer, progress func(written int64)) io.Writer {
	return &progressStream{w, progress}
}

type progressStream struct {
	w  io.Writer
	fn func(written int64)
}

func (s *progressStream) Write(buf []byte) (n int, err error) {
	n, err = s.w.Write(buf)
	s.fn(int64(n))
	return
}
