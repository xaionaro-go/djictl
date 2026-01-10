package djible

import (
	"github.com/xaionaro-go/djictl/pkg/duml"
)

type InterfaceAppToVideoTransmission Device

func (d *Device) AppToVideoTransmission() *InterfaceAppToVideoTransmission {
	return (*InterfaceAppToVideoTransmission)(d)
}

func (s *InterfaceAppToVideoTransmission) Device() *Device {
	return (*Device)(s)
}

func (s *InterfaceAppToVideoTransmission) InterfaceID() duml.InterfaceID {
	return duml.InterfaceIDAppToVideoTransmission
}
