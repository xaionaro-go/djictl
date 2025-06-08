package djictl

import (
	"bytes"
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
)

func (s *SubsystemStreamer) PrepareToLiveStream(
	ctx context.Context,
) (_err error) {
	logger.Tracef(ctx, "PrepareToLiveStream")
	defer func() { logger.Tracef(ctx, "/PrepareToLiveStream: %v", _err) }()

	err := s.SendMessagePrepareToLiveStreamStage1(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the message (stage1): %w", err)
	}

	logger.Debugf(ctx, "waiting for a streaming status")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-s.Device().getReceiveMessageChan(ctx, MessageTypePrepareToLiveStreamResult):
		logger.Debugf(ctx, "received a MessageTypePrepareToLiveStreamResult: %#+v", msg)
		if len(msg.Payload) != 1 {
			return fmt.Errorf("invalid payload size: %d", len(msg.Payload))
		}
		if msg.Payload[0] != 0x00 {
			return fmt.Errorf("expected the payload to be 0x00, but received 0x%X", msg.Payload)
		}
	}

	err = s.SendMessagePrepareToLiveStreamStage2(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the message (stage2): %w", err)
	}

	logger.Debugf(ctx, "waiting for the command result")
	select {
	case <-ctx.Done():
		return ctx.Err()
	case msg := <-s.Device().getReceiveMessageChan(ctx, MessageTypeStartStopStreamingResult):
		logger.Debugf(ctx, "received a command result: %#+v", msg)
	}

	return nil
}

func (s *SubsystemStreamer) SendMessagePrepareToLiveStreamStage1(
	ctx context.Context,
) error {
	msg := s.GetMessagePrepareToLiveStreamStage1()
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemStreamer) GetMessagePrepareToLiveStreamStage1() *Message {
	return &Message{
		Subsystem: s.SubsystemID(),
		ID:        MessageIDPrepareToLiveStreamStage1,
		Type:      MessageTypePrepareToLiveStream,
		Payload:   s.GetMessagePayloadPrepareToLiveStreamStage1(),
	}
}

func (s *SubsystemStreamer) GetMessagePayloadPrepareToLiveStreamStage1() []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x1A}))
	return buf.Bytes()
}

func (s *SubsystemStreamer) SendMessagePrepareToLiveStreamStage2(
	ctx context.Context,
) error {
	msg := s.GetMessagePrepareToLiveStreamStage2()
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemStreamer) GetMessagePrepareToLiveStreamStage2() *Message {
	return &Message{
		Subsystem: s.SubsystemID(),
		ID:        MessageIDStartStreaming,
		Type:      MessageTypeStartStopStreaming,
		Payload:   s.GetMessagePayloadPrepareToLiveStreamStage2(),
	}
}

func (s *SubsystemStreamer) GetMessagePayloadPrepareToLiveStreamStage2() []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x00, 0x01, 0x1C, 0x00}))
	return buf.Bytes()
}
