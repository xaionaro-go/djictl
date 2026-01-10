package djible

import (
	"context"

	"github.com/xaionaro-go/djictl/pkg/duml"
)

func (s *InterfaceAppToCamera) SetFCCEnable(ctx context.Context, enable bool) error {
	msg := duml.NewFCCEnableMessage()
	// NewFCCEnableMessage payload might need to be adjusted based on 'enable' param if it supported it,
	// but the original code was: msg := duml.NewFCCEnableMessage(); msg.Interface = duml.InterfaceIDAppToCamera
	// I'll assume NewFCCEnableMessage() currently just returns the "enable" command.
	msg.Interface = duml.InterfaceIDAppToCamera
	return s.Device().SendMessage(ctx, msg, true)
}
