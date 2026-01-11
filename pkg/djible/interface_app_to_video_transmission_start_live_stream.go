package djible

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

func (s *InterfaceAppToVideoTransmission) LiveStream(
	ctx context.Context,
	resolution duml.Resolution,
	bitrateKbps uint16,
	fps duml.FPS,
	rtmpURL string,
) (_err error) {
	logger.Tracef(ctx, "LiveStream(ctx, %v, %v, %v, %v)", resolution, bitrateKbps, fps, rtmpURL)
	defer func() {
		logger.Tracef(ctx, "/LiveStream(ctx, %v, %v, %v, %v): %v", resolution, bitrateKbps, fps, rtmpURL, _err)
	}()

	_, err := s.RequestConfigureLiveStream(ctx, resolution, bitrateKbps, fps, rtmpURL)
	if err != nil {
		return fmt.Errorf("unable to send the message to configure the live stream: %w", err)
	}

	_, err = s.RequestStartLiveStream(ctx)
	if err != nil {
		return fmt.Errorf("unable to send the message to start the live stream: %w", err)
	}

	for {
		msg, err := s.ReceiveMessageLiveStreamResult(ctx)
		if err != nil {
			return fmt.Errorf("unable to receive a response: %w", err)
		}
		status, err := duml.ParseBatteryStatus(ctx, msg.Payload)
		if err != nil {
			logger.Debugf(ctx, "unable to parse battery status: %v", err)
			continue
		}
		logger.Infof(ctx, "battery: %s", status.Capacity)
	}
}

func (s *InterfaceAppToVideoTransmission) RequestConfigureLiveStream(
	ctx context.Context,
	resolution duml.Resolution,
	bitrateKbps uint16,
	fps duml.FPS,
	rtmpURL string,
) (*duml.Message, error) {
	logger.Tracef(ctx, "RequestConfigureLiveStream")
	defer func() { logger.Tracef(ctx, "/RequestConfigureLiveStream") }()
	msg := s.GetMessageConfigureLiveStream(
		resolution, bitrateKbps, fps, rtmpURL,
	)
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToVideoTransmission) GetMessageConfigureLiveStream(
	resolution duml.Resolution,
	bitrateKbps uint16,
	fps duml.FPS,
	rtmpURL string,
) *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDConfigureStreaming,
		Type:      duml.MessageTypeConfigureStreaming,
		Payload: s.GetMessagePayloadConfigureLiveStream(
			resolution,
			bitrateKbps,
			fps,
			rtmpURL,
		),
	}
}

func (s *InterfaceAppToVideoTransmission) GetMessagePayloadConfigureLiveStream(
	resolution duml.Resolution,
	bitrateKbps uint16,
	fps duml.FPS,
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
	cannotFail(binary.Write(&buf, duml.BinaryOrder(), bitrateKbps))
	must(buf.Write([]byte{0x02, 0x00}))
	must(buf.Write(array1ToSlice(fps.BytesFixed())))
	must(buf.Write([]byte{0x00, 0x00, 0x00}))
	must(buf.Write(duml.PackURL(rtmpURL)))
	return buf.Bytes()
}

func (s *InterfaceAppToVideoTransmission) ReceiveMessageStartLiveStreamResult(
	ctx context.Context,
) (*duml.Message, error) {
	return s.Device().ReceiveMessage(ctx, duml.MessageTypeStartStopStreamingResult)
}

func (s *InterfaceAppToVideoTransmission) ReceiveMessageLiveStreamResult(
	ctx context.Context,
) (*duml.Message, error) {
	return s.Device().ReceiveMessage(ctx, duml.MessageTypeBatteryStatus)
}

func (s *InterfaceAppToVideoTransmission) RequestStartLiveStream(
	ctx context.Context,
) (*duml.Message, error) {
	msg := s.GetMessageStartLiveStream()
	return s.Device().Request(ctx, msg, true)
}

func (s *InterfaceAppToVideoTransmission) GetMessageStartLiveStream() *duml.Message {
	return &duml.Message{
		Interface: s.InterfaceID(),
		ID:        duml.MessageIDStartStreaming,
		Type:      duml.MessageTypeStartStopStreaming,
		Payload:   s.GetMessagePayloadStartLiveStream(),
	}
}

func (s *InterfaceAppToVideoTransmission) GetMessagePayloadStartLiveStream() []byte {
	var buf bytes.Buffer
	must(buf.Write([]byte{0x01, 0x01, 0x1A, 0x00, 0x01, 0x01}))
	return buf.Bytes()
}
