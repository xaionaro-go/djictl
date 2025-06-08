package djictl

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

func getLengthFromReader(r io.Reader) (uint64, bool) {
	switch r := r.(type) {
	case io.Seeker:
		p, err := r.Seek(0, io.SeekCurrent)
		if err != nil {
			return 0, false
		}
		return uint64(p), true
	default:
		return 0, false
	}
}
