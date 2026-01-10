package djible

import "github.com/xaionaro-go/djictl/pkg/duml"

type InterfaceAppToWiFiGroundStation Device

func (d *Device) AppToWiFiGroundStation() *InterfaceAppToWiFiGroundStation {
	return (*InterfaceAppToWiFiGroundStation)(d)
}

func (s *InterfaceAppToWiFiGroundStation) InterfaceID() duml.InterfaceID {
	return duml.InterfaceIDAppToWiFiGroundStation
}

func (s *InterfaceAppToWiFiGroundStation) Device() *Device {
	return (*Device)(s)
}
