package djiwifi

import (
	"fmt"

	"github.com/xaionaro-go/djictl/pkg/duml"
)

// Packet represents a DJI WiFi UDP packet wrapper.
// This wrapper is used to communicate with DJI Osmo/Pocket devices over WiFi (typically UDP port 9004).
//
// Protocol structure (observed from pcap):
// Offset | Size | Description
// -------|------|------------
// 0      | 1    | Packet Length
// 1      | 1    | Message Type (0x80 for control)
// 2      | 2    | Signature (0x47 0xa8)
// 4      | 16   | Metadata (timestamps, sequence numbers, etc.)
// 20     | n    | Payload (usually starts with DUML magic 0x55)

// Metadata represents the 16-byte metadata field in the WiFi wrapper.
type Metadata [MetadataSize]byte

type Packet struct {
	Length   uint16
	Type     MessageType
	WhType   WhType
	Metadata Metadata
	Payload  []byte
}

// ParsePacket parses a WiFi wrapper packet from bytes.
func ParsePacket(b []byte) (*Packet, error) {
	if len(b) < 1 {
		return nil, fmt.Errorf("empty packet")
	}

	// If it starts with DUML magic (0x55), it's likely a raw DUML message without the WiFi wrapper.
	if b[0] == duml.MessageStartMagicByte {
		if _, err := duml.ParseMessage(b); err == nil {
			return &Packet{
				Length:  uint16(len(b)),
				Type:    MessageTypeControl,
				Payload: b,
			}, nil
		}
	}

	if len(b) < HeaderLength {
		return nil, fmt.Errorf("packet too short: %d < %d", len(b), HeaderLength)
	}

	// 12-bit length in bytes 0 and 1
	length := uint16(b[0]) | (uint16(b[1]&0x0F) << 8)

	if b[OffsetSignature1] != Signature1 || b[OffsetSignature2] != Signature2 {
		return nil, fmt.Errorf("invalid signature: %02x%02x", b[OffsetSignature1], b[OffsetSignature2])
	}

	p := &Packet{
		Length: length,
		Type:   MessageType(b[OffsetMessageType] & 0xF0), // High nibble is MessageType
		WhType: WhType(b[OffsetWhType]),
	}
	copy(p.Metadata[:], b[OffsetMetadata:HeaderLength])
	p.Payload = b[HeaderLength:]

	return p, nil
}

// Bytes returns the serialized bytes of the packet.
func (p *Packet) Bytes() []byte {
	length := uint16(HeaderLength + len(p.Payload))
	b := make([]byte, HeaderLength+len(p.Payload))

	// 12-bit length encoding
	b[0] = uint8(length & 0xFF)
	b[1] = (uint8(p.Type) & 0xF0) | uint8((length>>8)&0x0F)

	b[OffsetSignature1] = Signature1
	b[OffsetSignature2] = Signature2
	copy(b[OffsetMetadata:HeaderLength], p.Metadata[:])
	b[OffsetWhType] = uint8(p.WhType)
	copy(b[HeaderLength:], p.Payload)
	return b
}

// DUMLMessage parses the payload as a DUML message if it starts with 0x55.
func (p *Packet) DUMLMessage() (*duml.Message, error) {
	if len(p.Payload) == 0 || p.Payload[0] != DUMLMagic {
		return nil, fmt.Errorf("payload is not a DUML message")
	}
	return duml.ParseMessage(p.Payload)
}

// NewDUMLPacket creates a new WiFi packet wrapping a DUML message.
func NewDUMLPacket(msg *duml.Message, metadata Metadata) *Packet {
	return &Packet{
		Type:     MessageTypeControl,
		Metadata: metadata,
		Payload:  msg.Bytes(),
	}
}
