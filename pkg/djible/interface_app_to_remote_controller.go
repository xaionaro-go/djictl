package djible

import (
	"context"

	"github.com/xaionaro-go/djictl/pkg/duml"
)

type InterfaceAppToRemoteController Device

func (d *Device) AppToRemoteController() *InterfaceAppToRemoteController {
	return (*InterfaceAppToRemoteController)(d)
}

func (s *InterfaceAppToRemoteController) InterfaceID() duml.InterfaceID {
	return duml.InterfaceIDAppToRemoteController
}

func (s *InterfaceAppToRemoteController) Device() *Device {
	return (*Device)(s)
}

func (s *InterfaceAppToRemoteController) SendData(ctx context.Context, data duml.RemoteControllerSimulatorData) error {
	msg := duml.NewRemoteControllerSimulatorMessage(data)
	msg.Interface = s.InterfaceID()
	return s.Device().SendMessage(ctx, msg, true)
}
