package duml

import (
	"bytes"
	"fmt"
	"io"
)

func read(r io.Reader, b []byte) error {
	n, err := r.Read(b)
	if err != nil {
		return err
	}
	if n != len(b) {
		return fmt.Errorf("expected to read %d bytes, but read %d", len(b), n)
	}
	return nil
}

func expectToRead(r io.Reader, b []byte) error {
	buf := make([]byte, len(b))
	err := read(r, buf)
	if err != nil {
		return err
	}
	if !bytes.Equal(buf, b) {
		return fmt.Errorf("expected to read %X, but read %X", b, buf)
	}
	return nil
}

func readOneByte(r io.Reader) (byte, error) {
	var buf [1]byte
	err := read(r, buf[:])
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func readBytes(r io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	err := read(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
