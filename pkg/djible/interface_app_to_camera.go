package djible

import "github.com/xaionaro-go/djictl/pkg/duml"

type InterfaceAppToCamera Device

func (d *Device) AppToCamera() *InterfaceAppToCamera {
	return (*InterfaceAppToCamera)(d)
}

func (s *InterfaceAppToCamera) InterfaceID() duml.InterfaceID {
	return duml.InterfaceIDAppToCamera
}

func (s *InterfaceAppToCamera) Device() *Device {
	return (*Device)(s)
}
