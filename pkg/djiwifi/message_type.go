package djiwifi

type MessageType uint8

const (
	// MessageTypeStandard is the most common message type (0x80).
	MessageTypeStandard MessageType = 0x80
	// MessageTypeControl is used in some contexts (0x30).
	MessageTypeControl MessageType = 0x30
)

func (t MessageType) String() string {
	switch t {
	case MessageTypeStandard:
		return "Standard"
	case MessageTypeControl:
		return "Control"
	default:
		return "Unknown"
	}
}

// WhType (WheeType) is a sub-protocol multiplexer / secondary message type.
type WhType uint8

const (
	// WhTypeHandshake is the WhType for handshake packets.
	WhTypeHandshake WhType = 0x00
	// WhTypeDroneCmd1 is the WhType for drone command type 1.
	WhTypeDroneCmd1 WhType = 0x01
	// WhTypeVideo is the WhType for H.264/H.265 video packets.
	WhTypeVideo WhType = 0x02
	// WhTypeDroneCmd2 is the WhType for drone command type 2.
	WhTypeDroneCmd2 WhType = 0x03
	// WhTypeOperatorCmd1 is the WhType for operator command type 1.
	WhTypeOperatorCmd1 WhType = 0x04
	// WhTypeOperatorCmd2 is the WhType for operator command type 2.
	WhTypeOperatorCmd2 WhType = 0x05
	// WhTypeOperatorCmd3 is the WhType for operator command type 3.
	WhTypeOperatorCmd3 WhType = 0x06
)

func (t WhType) String() string {
	switch t {
	case WhTypeHandshake:
		return "Handshake"
	case WhTypeDroneCmd1:
		return "DroneCmd1"
	case WhTypeVideo:
		return "Video"
	case WhTypeDroneCmd2:
		return "DroneCmd2"
	case WhTypeOperatorCmd1:
		return "OperatorCmd1"
	case WhTypeOperatorCmd2:
		return "OperatorCmd2"
	case WhTypeOperatorCmd3:
		return "OperatorCmd3"
	default:
		return "Unknown"
	}
}
