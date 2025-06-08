package djictl

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/facebookincubator/go-belt/tool/logger"
)

var (
	MessageStartMagic = []byte{0x55}
	protocolVersion   = uint8(0x04)
)

type Message struct {
	Subsystem SubsystemID
	ID        MessageID
	Type      MessageType
	Payload   []byte
}

func ParseMessage(b []byte) (_ret *Message, _err error) {
	logger.Tracef(context.TODO(), "ParseMessage(%X)", b)
	defer func() { logger.Tracef(context.TODO(), "/ParseMessage(%X): %p %v", b, _ret, _err) }()
	if len(b) < 13 {
		return nil, fmt.Errorf("%w: expected at least 13 bytes, but received only %d", io.ErrUnexpectedEOF, len(b))
	}
	return ParseMessageFromReader(bytes.NewReader(b))
}

func ParseMessageFromReader(r io.Reader) (*Message, error) {
	var buf bytes.Buffer
	r = io.TeeReader(r, &buf)

	if err := expectToRead(r, MessageStartMagic); err != nil {
		return nil, fmt.Errorf("invalid beginning magic in Message: %w", err)
	}

	length, err := readOneByte(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read the length: %w", err)
	}
	if expectedLength, ok := getLengthFromReader(r); ok {
		if uint64(length) != expectedLength {
			return nil, fmt.Errorf("invalid length: header contains: %d, while the actual Message length is %d", length, expectedLength)
		}
	}

	version, err := readOneByte(r)
	if err != nil {
		return nil, fmt.Errorf("unexpected read the version: %w", err)
	}
	if version != protocolVersion {
		return nil, fmt.Errorf("unexpected version: received:0x%02X expected:0x%02X", version, protocolVersion)
	}

	headerBytes := buf.Bytes()
	expectedHeaderCRC := crc8(headerBytes)
	headerCRC, err := readOneByte(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read the header CRC: %w", err)
	}
	if headerCRC != expectedHeaderCRC {
		return nil, fmt.Errorf("header CRC does not match: received:%02X expected:%02X (header bytes: %X)", headerCRC, expectedHeaderCRC, headerBytes)
	}

	var msg Message

	if err := binary.Read(r, binary.BigEndian, &msg.Subsystem); err != nil {
		return nil, fmt.Errorf("unable to read the subsystem ID: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &msg.ID); err != nil {
		return nil, fmt.Errorf("unable to read the ID: %w", err)
	}

	if err := msg.Type.ParseFrom(r); err != nil {
		return nil, fmt.Errorf("unable to read the SubsystemMethod: %w", err)
	}

	payloadWithCRC, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read the payload: %w", err)
	}
	if len(payloadWithCRC) < 2 {
		return nil, fmt.Errorf("not enough bytes left for CRC16: left:%d, expected:2", len(payloadWithCRC))
	}
	msg.Payload = payloadWithCRC[:len(payloadWithCRC)-2]
	msgCRC16 := BinaryOrder().Uint16(payloadWithCRC[len(payloadWithCRC)-2:])
	expectedMsgCRC16 := crc16(buf.Bytes()[:buf.Len()-2])
	if msgCRC16 != expectedMsgCRC16 {
		return nil, fmt.Errorf("the full Message CRC16 does not match: received:%04X, expected:%04X", msgCRC16, expectedHeaderCRC)
	}

	return &msg, nil
}

func (msg *Message) Bytes() []byte {
	if len(msg.Payload) > math.MaxUint8-13 {
		panic(fmt.Errorf("the payload is too long: %d > %d", len(msg.Payload), math.MaxUint8-13))
	}
	var buf bytes.Buffer
	must(buf.Write(MessageStartMagic))
	must(buf.Write([]byte{13 + uint8(len(msg.Payload))}))
	must(buf.Write([]byte{protocolVersion}))
	must(buf.Write([]byte{crc8(buf.Bytes())}))
	cannotFail(binary.Write(&buf, binary.BigEndian, msg.Subsystem))
	cannotFail(binary.Write(&buf, binary.BigEndian, msg.ID))
	must(buf.Write(msg.Type.Bytes()))
	must(buf.Write(msg.Payload))
	cannotFail(binary.Write(&buf, BinaryOrder(), crc16(buf.Bytes())))

	//           M  ln pr C8 subs id   type   payload
	// DELETEME: 55 12 04 C7 0402 AEB5 000427 0000080000B684
	//           55 0e 04 66 0207 0400 c00746 00 3528
	//           55 13 04 03 0208 b5bb 40028e 01011a000102385f
	//           55 13 04 03 0208 b4bb 40028e 01011a0001013238
	//           55 42 04 b0 0208 b3bb 400878 0032000a7017020003000...

	return buf.Bytes()
}
