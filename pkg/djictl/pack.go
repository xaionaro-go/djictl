package djictl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func packURL(in string) []byte {
	if len(in) > math.MaxUint16 {
		panic(fmt.Errorf("too long string: %d > %d", len(in), math.MaxUint16))
	}
	var buf bytes.Buffer
	binary.Write(&buf, BinaryOrder(), uint16(len(in)))
	must(buf.WriteString(in))
	return buf.Bytes()
}

func packString(in string) []byte {
	if len(in) > math.MaxUint8 {
		panic(fmt.Errorf("too long string: %d > %d", len(in), math.MaxUint8))
	}
	var buf bytes.Buffer
	buf.WriteByte(byte(len(in)))
	must(buf.WriteString(in))
	return buf.Bytes()
}
