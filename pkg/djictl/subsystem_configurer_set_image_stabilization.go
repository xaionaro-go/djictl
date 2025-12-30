package djictl

import (
	"bytes"
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
)

func (s *SubsystemConfigurer) SetImageStabilization(
	ctx context.Context,
	v ImageStabilization,
) error {
	err := s.SendMessageSetImageStabilization(ctx, v)
	if err != nil {
		return fmt.Errorf("unable to send the Message: %w", err)
	}
	msg, err := s.ReceiveMessageSetImageStabilizationResult(ctx)
	if err != nil {
		return fmt.Errorf("unable to receive a response: %w", err)
	}
	logger.Debugf(ctx, "got set image stabilization result payload: %X", msg.Payload)
	return nil
}

func (s *SubsystemConfigurer) SendMessageSetImageStabilization(
	ctx context.Context,
	v ImageStabilization,
) error {
	msg := s.GetMessageSetImageStabilization(v)
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemConfigurer) GetMessageSetImageStabilization(
	v ImageStabilization,
) *Message {
	panic("not implemented")
	return &Message{
		Subsystem: s.SubsystemID(),
		ID:        0, //MessageIDSetImageStabilization,
		Type:      MessageTypeConfigure,
		Payload: s.GetMessagePayloadSetImageStabilization(
			v,
		),
	}
}

func (s *SubsystemConfigurer) GetMessagePayloadSetImageStabilization(
	v ImageStabilization,
) []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x01, 0x01}))
	must(buf.Write(array1ToSlice(s.Type.BytesFixedSetImageStabilization())))
	must(buf.Write([]byte{0x00, 0x01}))
	must(buf.Write(array1ToSlice(v.BytesFixed())))
	return buf.Bytes()
}

func (s *SubsystemConfigurer) ReceiveMessageSetImageStabilizationResult(
	ctx context.Context,
) (*Message, error) {
	panic("not implemented")
	return s.Device().ReceiveMessage(ctx, MessageTypeConfigure)
}
