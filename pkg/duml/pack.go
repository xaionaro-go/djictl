package duml

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func PackURL(in string) []byte {
	if len(in) > math.MaxUint16 {
		panic(fmt.Errorf("too long string: %d > %d", len(in), math.MaxUint16))
	}
	var buf bytes.Buffer
	binary.Write(&buf, BinaryOrder(), uint16(len(in)))
	must(buf.WriteString(in))
	return buf.Bytes()
}

func PackString(in string) []byte {
	if len(in) > math.MaxUint8 {
		panic(fmt.Errorf("too long string: %d > %d", len(in), math.MaxUint8))
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(len(in)))
	must(buf.WriteString(in))
	return buf.Bytes()
}

func UnpackStringU16BE(b []byte) (string, error) {
	if len(b) < 2 {
		return "", fmt.Errorf("too short payload: %d", len(b))
	}
	length := binary.BigEndian.Uint16(b[:2])
	if len(b) < int(2+length) {
		return "", fmt.Errorf("payload too short for length %d: %d", length, len(b))
	}
	return string(b[2 : 2+length]), nil
}
