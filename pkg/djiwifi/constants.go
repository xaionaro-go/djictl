package djiwifi

const (
	// HeaderLength is the size of the WiFi wrapper header in bytes.
	HeaderLength = 20

	// Signature1 is the first byte of the WiFi wrapper signature.
	Signature1 = 0x47
	// Signature2 is the second byte of the WiFi wrapper signature.
	Signature2 = 0xa8

	// MetadataSize is the size of the metadata field in the WiFi wrapper.
	MetadataSize = 16

	// DUMLMagic is the first byte of a DUML message.
	DUMLMagic = 0x55

	// DefaultUDPPort is the default UDP port used by DJI devices for WiFi communication.
	DefaultUDPPort = 9004

	// ReadBufferSize is the default size for the UDP read buffer.
	ReadBufferSize = 2048

	// OffsetPacketLength is the byte offset of the packet length field.
	// It is a 12-bit little-endian value in bytes 0 and 1.
	OffsetPacketLength = 0
	// OffsetMessageType is the byte offset of the message type field.
	OffsetMessageType = 1
	// OffsetSignature1 is the byte offset of the first signature byte.
	OffsetSignature1 = 2
	// OffsetSignature2 is the byte offset of the second signature byte.
	OffsetSignature2 = 3
	// OffsetMetadata is the byte offset where metadata starts.
	OffsetMetadata = 4
	// OffsetWhType is the byte offset of the WhType (Sub-protocol multiplexer).
	OffsetWhType = 6

	// StatusProductInfoSize is the size of the product info block in a status report.
	StatusProductInfoSize = 6
	// StatusStreamCapabilitySize is the size of each stream capability block in a status report.
	StatusStreamCapabilitySize = 11

	// OffsetStatusMTU is the byte offset of the MTU field in a StreamCapability block (Little Endian).
	// Common value: 0xc0 05 -> 1472 bytes.
	OffsetStatusMTU = 0
	// OffsetStatusFrameInterval is the byte offset of the frame interval field in a StreamCapability block (Little Endian).
	// Common value: 0x14 00 -> 20ms (roughly 50 FPS).
	OffsetStatusFrameInterval = 2
	// OffsetStatusQuality is the byte offset of the quality/bitrate field in a StreamCapability block (Little Endian).
	// Values from 1 to 100. Observed: 0x64 00 -> 100.
	OffsetStatusQuality = 4

	// ProtocolUDP is the network protocol string for UDP.
	ProtocolUDP = "udp"
)

var (
	// MetadataInitial is the metadata block used in the initial status packet (Camera to App).
	// Content analysis:
	// [0:4]   - Sequence number or timestamp (e.g., 0000005f)
	// [4:8]   - Unknown (e.g., 38426400)
	// [8:10]  - Link Quality or Bitrate (e.g., 6400 -> 100)
	// [10:12] - MTU (e.g., c005 -> 1472)
	// [12:14] - Frame/Packet Interval (e.g., 1400 -> 20ms)
	// [14:16] - Unknown (e.g., 0064)
	MetadataInitial = Metadata{0x00, 0x00, 0x00, 0x5f, 0x38, 0x42, 0x64, 0x00, 0x64, 0x00, 0xc0, 0x05, 0x14, 0x00, 0x00, 0x64}

	// PayloadInitial is the payload used in the initial status packet (Camera to App).
	// It is a status report containing product info and stream capabilities.
	// Product Info: 00 00 01 90 01
	// Stream 1: MTU=1472, Interval=20, Quality=100
	// Stream 2: MTU=1472, Interval=20, Quality=100
	PayloadInitial = []byte{
		0x00, 0x00, 0x01, 0x90, 0x01, 0xc0, 0x05, 0x14, 0x00, 0x00, 0x64, 0x00, 0x14, 0x00, 0x64, 0x00,
		0xc0, 0x05, 0x14, 0x00, 0x00, 0x64, 0x00, 0x01, 0x01, 0x04, 0x01, 0x02,
	}

	// MetadataApp is the metadata block used for App-initiated DUML commands.
	// Content analysis:
	// [0:8]   - App-side timestamps/session IDs
	// [12:14] - WiFi layer sequence number
	MetadataApp = Metadata{0x40, 0x42, 0x05, 0x47, 0x38, 0x42, 0x40, 0x42, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x00, 0x00}
	// PayloadAppIdentifier is the DUML payload for the "#sAPP" command.
	PayloadAppIdentifier = []byte{0x73, 0x41, 0x50, 0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x10}

	// PayloadHandshakeRMVT is the magic "RMVT" handshake used to trigger video streaming.
	PayloadHandshakeRMVT = []byte{0x52, 0x4d, 0x56, 0x54, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
)
