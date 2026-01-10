package main

import (
	"context"
	"fmt"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/xaionaro-go/djictl/pkg/djible"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

func connectWiFiAndStartStreaming(ctx context.Context, dev *djible.Device, wifiSSID, wifiPSK, rtmpURL string) error {
	logger.Infof(ctx, "found device %s; initializing...", dev)

	err := dev.Init(ctx)
	if err != nil {
		return fmt.Errorf("unable to initialize the connection to the device: %w", err)
	}
	logger.Infof(ctx, "requesting to pair")
	err = dev.AppToWiFiGroundStation().Pair(ctx)
	if err != nil {
		return fmt.Errorf("unable to pair: %w", err)
	}

	logger.Infof(ctx, "prepare to live stream")
	err = dev.AppToVideoTransmission().PrepareToLiveStream(ctx)
	if err != nil {
		return fmt.Errorf("unable to request the device to prepare to live stream: %w", err)
	}
	logger.Infof(ctx, "requesting to connect to WiFi")
	err = dev.AppToWiFiGroundStation().ConnectToWiFi(ctx, wifiSSID, wifiPSK)
	if err != nil {
		return fmt.Errorf("unable to make the device connect to our WiFi: %w", err)
	}
	switch dev.Type {
	case duml.DeviceTypeOsmoAction4, duml.DeviceTypeOsmoAction5Pro:
		logger.Infof(ctx, "set image stabilization")
		err = dev.AppToCamera().SetImageStabilization(ctx, duml.ImageStabilizationRockSteadyPlus)
		if err != nil {
			return fmt.Errorf("unable to set image stabilization to RockSteadyPlus: %w", err)
		}
	}
	logger.Infof(ctx, "start live stream")
	err = dev.AppToVideoTransmission().LiveStream(ctx, duml.Resolution1080p, 5000, duml.FPS25, rtmpURL)
	if err != nil {
		return fmt.Errorf("unable to make the device to stream: %w", err)
	}
	logger.Infof(ctx, "done")
	return nil
}
