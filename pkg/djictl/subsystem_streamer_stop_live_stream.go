package djictl

import (
	"bytes"
	"context"
	"fmt"
)

func (s *SubsystemStreamer) StopLiveStream(
	ctx context.Context,
) error {
	err := s.SendMessageStopLiveStream(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the Message: %w", err)
	}
	_, err = s.ReceiveMessageStopLiveStreamResult(ctx)
	if err != nil {
		return fmt.Errorf("unable to receive a response: %w", err)
	}
	return nil
}

func (s *SubsystemStreamer) SendMessageStopLiveStream(
	ctx context.Context,
) error {
	msg := s.GetMessageStopLiveStream()
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemStreamer) GetMessageStopLiveStream() *Message {
	return &Message{
		Subsystem: s.SubsystemID(),
		ID:        MessageIDStopStreaming,
		Type:      MessageTypeStartStopStreaming,
		Payload:   s.GetMessagePayloadStopLiveStream(),
	}
}

func (s *SubsystemStreamer) GetMessagePayloadStopLiveStream() []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x01, 0x01, 0x1A, 0x00, 0x01, 0x02}))
	return buf.Bytes()
}

func (s *SubsystemStreamer) ReceiveMessageStopLiveStreamResult(
	ctx context.Context,
) (*Message, error) {
	return s.Device().ReceiveMessage(ctx, MessageTypeStartStopStreaming)
}
