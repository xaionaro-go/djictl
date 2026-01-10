package duml

import (
	"encoding/binary"
)

var binaryOrder = binary.LittleEndian

func BinaryOrder() binary.ByteOrder {
	return binaryOrder
}
