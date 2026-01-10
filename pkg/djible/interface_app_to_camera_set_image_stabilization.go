package djible

import (
	"bytes"
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

func (s *InterfaceAppToCamera) SetImageStabilization(
	ctx context.Context,
	v duml.ImageStabilization,
) error {
	err := s.SendMessageSetImageStabilization(ctx, v)
	if err != nil {
		return fmt.Errorf("unable to send the duml.Message: %w", err)
	}
	msg, err := s.ReceiveMessageSetImageStabilizationResult(ctx)
	if err != nil {
		return fmt.Errorf("unable to receive a response: %w", err)
	}
	logger.Debugf(ctx, "got set image stabilization result payload: %X", msg.Payload)
	return nil
}

func (s *InterfaceAppToCamera) SendMessageSetImageStabilization(
	ctx context.Context,
	v duml.ImageStabilization,
) error {
	msg := s.GetMessageSetImageStabilization(v)
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *InterfaceAppToCamera) GetMessageSetImageStabilization(
	v duml.ImageStabilization,
) *duml.Message {
	panic("not implemented")
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        0, //duml.MessageIDSetImageStabilization,
		Type:      duml.MessageTypeStartStopStreaming,
		Payload: s.GetMessagePayloadSetImageStabilization(
			v,
		),
	}
}

func (s *InterfaceAppToCamera) GetMessagePayloadSetImageStabilization(
	v duml.ImageStabilization,
) []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x01, 0x01}))
	must(buf.Write(array1ToSlice(s.Device().Type.BytesFixedSetImageStabilization())))
	must(buf.Write([]byte{0x00, 0x01}))
	must(buf.Write(array1ToSlice(v.BytesFixed())))
	return buf.Bytes()
}

func (s *InterfaceAppToCamera) ReceiveMessageSetImageStabilizationResult(
	ctx context.Context,
) (*duml.Message, error) {
	panic("not implemented")
	return s.Device().ReceiveMessage(ctx, duml.MessageTypeStartStopStreaming)
}
