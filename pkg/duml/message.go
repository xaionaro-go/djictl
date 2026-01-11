package duml

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/facebookincubator/go-belt/tool/logger"
)

const (
	MessageStartMagicByte = 0x55
)

var (
	MessageStartMagic = []byte{MessageStartMagicByte}
	ProtocolVersion   = uint8(0x01)
)

type Message struct {
	Interface InterfaceID
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

const (
	// the contribution of the headers to the total length
	totalLengthHeaders = 11

	// totalLengthHeaders + the tail (which consists of CRC16) length
	totalLengthHeadersAndTail = totalLengthHeaders + 2
)

func ParseMessageFromReader(r io.Reader) (*Message, error) {
	var buf bytes.Buffer
	r = io.TeeReader(r, &buf)

	if err := expectToRead(r, MessageStartMagic); err != nil {
		return nil, fmt.Errorf("invalid beginning magic in Message: %w", err)
	}

	b, err := readBytes(r, 2)
	if err != nil {
		return nil, fmt.Errorf("unable to read the length/version: %w", err)
	}

	totalLength := uint16(b[0]) | (uint16(b[1]&0x03) << 8)
	if totalLength < totalLengthHeaders {
		return nil, fmt.Errorf("invalid total length in Message: %d < 13", totalLength)
	}

	version := b[1] >> 2
	if version != ProtocolVersion {
		return nil, fmt.Errorf("unexpected version: received:0x%02X expected:0x%02X", version, ProtocolVersion)
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

	r = io.LimitReader(r, int64(totalLength)-4)

	var msg Message

	if err := binary.Read(r, binary.BigEndian, &msg.Interface); err != nil {
		return nil, fmt.Errorf("unable to read the interface ID: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &msg.ID); err != nil {
		return nil, fmt.Errorf("unable to read the ID: %w", err)
	}

	if err := msg.Type.ParseFrom(r); err != nil {
		return nil, fmt.Errorf("unable to read the InterfaceMethod: %w", err)
	}

	payloadWithCRC, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("unable to read the payload: %w", err)
	}

	expectedLength := int(totalLength) - int(totalLengthHeaders)
	if len(payloadWithCRC) < expectedLength {
		return nil, fmt.Errorf("(%w) payload is too short: expected at least %d bytes, but received only %d", io.ErrUnexpectedEOF, expectedLength, len(payloadWithCRC))
	}
	if len(payloadWithCRC) > expectedLength {
		return nil, fmt.Errorf("payload is too long: expected %d bytes, but received %d", expectedLength, len(payloadWithCRC))
	}
	msg.Payload = payloadWithCRC[:len(payloadWithCRC)-2]
	msgCRC16 := BinaryOrder().Uint16(payloadWithCRC[len(payloadWithCRC)-2:])
	expectedMsgCRC16 := crc16(buf.Bytes()[:buf.Len()-2])
	if msgCRC16 != expectedMsgCRC16 {
		return nil, fmt.Errorf("the full Message CRC16 does not match: received:%04X, expected:%04X (payloadWithCRC:%X; totalLength:%d)", msgCRC16, expectedMsgCRC16, payloadWithCRC, totalLength)
	}

	return &msg, nil
}

func (msg *Message) Bytes() []byte {
	if len(msg.Payload) > 1023-totalLengthHeadersAndTail {
		panic(fmt.Errorf("the payload is too long: %d > %d", len(msg.Payload), 1023-totalLengthHeadersAndTail))
	}
	var buf bytes.Buffer
	must(buf.Write(MessageStartMagic))
	length := uint16(len(msg.Payload)) + uint16(totalLengthHeadersAndTail)
	must(buf.Write([]byte{byte(length & 0xff)}))
	must(buf.Write([]byte{(ProtocolVersion << 2) | byte(length>>8)&0x03}))
	must(buf.Write([]byte{crc8(buf.Bytes())}))
	cannotFail(binary.Write(&buf, binary.BigEndian, msg.Interface))
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

func (msg *Message) String() string {
	return fmt.Sprintf("Message{Interface:%s ID:%s Type:%s Payload:%X}", msg.Interface, msg.ID, msg.Type, msg.Payload)
}
