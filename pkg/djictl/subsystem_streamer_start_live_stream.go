package djictl

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
)

func (s *SubsystemStreamer) LiveStream(
	ctx context.Context,
	resolution Resolution,
	bitrateKbps uint16,
	fps FPS,
	rtmpURL string,
) (_err error) {
	logger.Tracef(ctx, "LiveStream(ctx, %v, %v, %v, %v)", resolution, bitrateKbps, fps, rtmpURL)
	defer func() {
		logger.Tracef(ctx, "/LiveStream(ctx, %v, %v, %v, %v): %v", resolution, bitrateKbps, fps, rtmpURL, _err)
	}()

	err := s.SendMessageConfigureLiveStream(ctx, resolution, bitrateKbps, fps, rtmpURL)
	if err != nil {
		return fmt.Errorf("unable to send the message to configure the live stream: %w", err)
	}

	err = s.SendMessageStartLiveStream(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the message to start the live stream: %w", err)
	}

	for {
		msg, err := s.ReceiveMessageStartLiveStreamResult(ctx)
		if err != nil {
			return fmt.Errorf("unable to receive a response: %w", err)
		}
		if msg.ID != MessageIDStartStreaming {
			logger.Debugf(ctx, "received an unexpected Message, ID:%X", msg.ID)
			continue
		}
		break
	}
	for {
		msg, err := s.ReceiveMessageLiveStreamResult(ctx)
		if err != nil {
			return fmt.Errorf("unable to receive a response: %w", err)
		}
		if len(msg.Payload) < 21 {
			continue
		}
		batteryPercentage := msg.Payload[20]
		logger.Infof(ctx, "battery: %d%%", batteryPercentage)
	}
}

func (s *SubsystemStreamer) SendMessageConfigureLiveStream(
	ctx context.Context,
	resolution Resolution,
	bitrateKbps uint16,
	fps FPS,
	rtmpURL string,
) error {
	msg := s.GetMessageConfigureLiveStream(
		resolution, bitrateKbps, fps, rtmpURL,
	)
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemStreamer) GetMessageConfigureLiveStream(
	resolution Resolution,
	bitrateKbps uint16,
	fps FPS,
	rtmpURL string,
) *Message {
	return &Message{
		Subsystem: s.SubsystemID(),
		ID:        MessageIDConfigureStreaming,
		Type:      MessageTypeConfigureStreaming,
		Payload: s.GetMessagePayloadConfigureLiveStream(
			resolution,
			bitrateKbps,
			fps,
			rtmpURL,
		),
	}
}

func (s *SubsystemStreamer) GetMessagePayloadConfigureLiveStream(
	resolution Resolution,
	bitrateKbps uint16,
	fps FPS,
	rtmpURL string,
) []byte {
	// = packet example =
	// hdr: 55 42 04 b0 0208 b3bb 400878
	// payload: 00 32 00 0a 7017 0200 03 000000 270072746d703a2f2f3139322e3136382e302e3133313a313934362f746573742f73747265616d302f995c
	var buf bytes.Buffer
	must(buf.Write([]byte{0x00}))
	must(buf.Write(array1ToSlice(s.Type.BytesFixedStartStreaming())))
	must(buf.Write([]byte{0x00}))
	must(buf.Write(array1ToSlice(resolution.BytesFixed())))
	cannotFail(binary.Write(&buf, BinaryOrder(), bitrateKbps))
	must(buf.Write([]byte{0x02, 0x00}))
	must(buf.Write(array1ToSlice(fps.BytesFixed())))
	must(buf.Write([]byte{0x00, 0x00, 0x00}))
	must(buf.Write(packURL(rtmpURL)))
	return buf.Bytes()
}

func (s *SubsystemStreamer) ReceiveMessageStartLiveStreamResult(
	ctx context.Context,
) (*Message, error) {
	return s.Device().ReceiveMessage(ctx, MessageTypeStartStopStreamingResult)
}

func (s *SubsystemStreamer) ReceiveMessageLiveStreamResult(
	ctx context.Context,
) (*Message, error) {
	return s.Device().ReceiveMessage(ctx, MessageTypeStreamingStatus)
}

func (s *SubsystemStreamer) SendMessageStartLiveStream(
	ctx context.Context,
) error {
	msg := s.GetMessageStartLiveStream()
	return s.Device().SendMessage(ctx, msg, true)
}

func (s *SubsystemStreamer) GetMessageStartLiveStream() *Message {
	return &Message{
		Subsystem: s.SubsystemID(),
		ID:        MessageIDStartStreaming,
		Type:      MessageTypeStartStopStreaming,
		Payload:   s.GetMessagePayloadStartLiveStream(),
	}
}

func (s *SubsystemStreamer) GetMessagePayloadStartLiveStream() []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x01, 0x01, 0x1A, 0x00, 0x01, 0x01}))
	return buf.Bytes()
}
