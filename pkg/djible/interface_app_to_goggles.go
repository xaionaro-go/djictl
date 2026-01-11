package djible

import (
	"context"

	"github.com/xaionaro-go/djictl/pkg/duml"
)

type InterfaceAppToGoggles Device

func (d *Device) AppToGoggles() *InterfaceAppToGoggles {
	return (*InterfaceAppToGoggles)(d)
}

func (s *InterfaceAppToGoggles) InterfaceID() duml.InterfaceID {
	return duml.InterfaceIDAppToGoggles
}

func (s *InterfaceAppToGoggles) Device() *Device {
	return (*Device)(s)
}

func (s *InterfaceAppToGoggles) SetMode(ctx context.Context, mode duml.GogglesMode) (*duml.Message, error) {
	msg := duml.NewGogglesModeMessage(mode)
	msg.Interface = s.InterfaceID()
	return s.Device().Request(ctx, msg, true)
}
