package djible

import (
	"context"

	"github.com/xaionaro-go/djictl/pkg/duml"
)

type InterfaceAppToBattery Device

func (d *Device) AppToBattery() *InterfaceAppToBattery {
	return (*InterfaceAppToBattery)(d)
}

func (s *InterfaceAppToBattery) InterfaceID() duml.InterfaceID {
	return duml.InterfaceIDAppToBattery
}

func (s *InterfaceAppToBattery) Device() *Device {
	return (*Device)(s)
}

func (s *InterfaceAppToBattery) GetInfo(ctx context.Context) error {
	msg := &duml.Message{
		Interface: s.InterfaceID(),
		Type:      duml.MessageTypeGetBatteryInfo,
	}
	return s.Device().SendMessage(ctx, msg, true)
}
