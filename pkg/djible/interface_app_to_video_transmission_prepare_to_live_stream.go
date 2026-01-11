package djible

import (
	"bytes"
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

func (s *InterfaceAppToVideoTransmission) PrepareToLiveStream(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "PrepareToLiveStream")
	defer func() { logger.Tracef(ctx, "/PrepareToLiveStream: %v", _err) }()

	msg, err := s.RequestPrepareToLiveStreamStage1(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the message (stage1): %w", err)
	}

	logger.Debugf(ctx, "received a duml.MessageTypePrepareToLiveStreamResult: %#+v", msg)
	if len(msg.Payload) != 1 {
		return fmt.Errorf("invalid payload size: %d", len(msg.Payload))
	}
	if msg.Payload[0] != 0x00 {
		return fmt.Errorf("expected the payload to be 0x00, but received 0x%X", msg.Payload)
	}

	msg, err = s.RequestPrepareToLiveStreamStage2(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the message (stage2): %w", err)
	}

	logger.Debugf(ctx, "received a command result: %#+v", msg)

	return nil
}

func (s *InterfaceAppToVideoTransmission) RequestPrepareToLiveStreamStage1(
	ctx context.Context,
) (*duml.Message, error) {
	msg := s.GetMessagePrepareToLiveStreamStage1()
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToVideoTransmission) GetMessagePrepareToLiveStreamStage1() *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDPrepareToLiveStreamStage1,
		Type:      duml.MessageTypePrepareToLiveStream,
		Payload:   s.GetMessagePayloadPrepareToLiveStreamStage1(),
	}
}

func (s *InterfaceAppToVideoTransmission) GetMessagePayloadPrepareToLiveStreamStage1() []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x1A}))
	return buf.Bytes()
}

func (s *InterfaceAppToVideoTransmission) RequestPrepareToLiveStreamStage2(
	ctx context.Context,
) (*duml.Message, error) {
	msg := s.GetMessagePrepareToLiveStreamStage2()
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToVideoTransmission) GetMessagePrepareToLiveStreamStage2() *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDStartStreaming,
		Type:      duml.MessageTypeStartStopStreaming,
		Payload:   s.GetMessagePayloadPrepareToLiveStreamStage2(),
	}
}

func (s *InterfaceAppToVideoTransmission) GetMessagePayloadPrepareToLiveStreamStage2() []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x00, 0x01, 0x1C, 0x00}))
	return buf.Bytes()
}
