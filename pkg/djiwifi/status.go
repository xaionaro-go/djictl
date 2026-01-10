package djiwifi

import (
	"encoding/binary"
	"fmt"
)

// StreamCapability represents one stream configuration announced by the device in 0x80 status packets.
type StreamCapability struct {
	// MTU is the maximum transmission unit for UDP packets.
	// Observed as 0xc0 05 -> 1472.
	MTU uint16

	// FrameInterval or PacketInterval.
	// Observed as 0x14 00 -> 20.
	FrameInterval uint16

	// Quality or Bitrate related setting.
	// Observed as 0x64 00 -> 100.
	Quality uint16

	// Raw is the original 11-byte block for further analysis.
	Raw [StatusStreamCapabilitySize]byte
}

// StatusReport contains the decoded information from a 0x80 status packet.
type StatusReport struct {
	// ProductInfo is the first 6 bytes of the status payload.
	ProductInfo [StatusProductInfoSize]byte

	// Streams are the configuration blocks (11 bytes each).
	Streams []StreamCapability

	// Tail is the last byte of the payload.
	Tail uint8
}

// ParseStatusReport parses the payload of a MessageTypeControl (0x80) packet
// that contains device status and stream capabilities.
//
// Payload Structure:
// [0:6]   - Product Information (Model IDs, etc.)
// [6:]    - Stream Capability blocks (11 bytes each) followed by 1-byte tail.
func ParseStatusReport(payload []byte) (*StatusReport, error) {
	if len(payload) < StatusProductInfoSize+1 {
		return nil, fmt.Errorf("status payload too short: %d", len(payload))
	}

	report := &StatusReport{}
	copy(report.ProductInfo[:], payload[0:StatusProductInfoSize])

	// Each stream block is 11 bytes.
	// The rest of the payload minus the 1-byte tail should be a multiple of 11.
	streamData := payload[StatusProductInfoSize : len(payload)-1]
	if len(streamData)%StatusStreamCapabilitySize != 0 {
		return nil, fmt.Errorf("invalid stream data length: %d (not a multiple of %d)", len(streamData), StatusStreamCapabilitySize)
	}

	for i := 0; i < len(streamData); i += StatusStreamCapabilitySize {
		block := streamData[i : i+StatusStreamCapabilitySize]
		cap := StreamCapability{
			MTU:           binary.LittleEndian.Uint16(block[OffsetStatusMTU : OffsetStatusMTU+2]),
			FrameInterval: binary.LittleEndian.Uint16(block[OffsetStatusFrameInterval : OffsetStatusFrameInterval+2]),
			Quality:       binary.LittleEndian.Uint16(block[OffsetStatusQuality : OffsetStatusQuality+2]), // based on research decomposition
		}
		copy(cap.Raw[:], block)
		report.Streams = append(report.Streams, cap)
	}

	report.Tail = payload[len(payload)-1]

	return report, nil
}
