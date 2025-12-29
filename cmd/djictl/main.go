package main

import (
	"strings"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/spf13/pflag"
	"github.com/xaionaro-go/djictl/pkg/djictl"
)

func main() {
	loggerLevel := logger.LevelInfo
	pflag.Var(&loggerLevel, "log-level", "Log level")
	filterDeviceAddr := pflag.String("filter-device-addr", "", "")
	wifiSSID := pflag.String("wifi-ssid", "", "")
	wifiPSK := pflag.String("wifi-psk", "", "")
	rtmpURL := pflag.String("rtmp-url", "", "")
	pflag.Parse()

	ctx := getContext(loggerLevel, false, "")

	if *wifiSSID == "" {
		logger.Fatalf(ctx, "please set wifi SSID")
	}
	if *wifiPSK == "" {
		logger.Fatalf(ctx, "please set wifi PSK")
	}
	if *rtmpURL == "" {
		logger.Fatalf(ctx, "please set the RTMP URL")
	}

	devCh, errCh, err := djictl.Scan(ctx)
	if err != nil {
		logger.Fatalf(ctx, "%v", err)
	}

	for {
		select {
		case dev := <-devCh:
			logger.Debugf(ctx, "found device %s", dev)
			if !strings.Contains(strings.ToLower(dev.ID.String()), strings.ToLower(*filterDeviceAddr)) {
				logger.Infof(ctx, "found device %s; but skipping, because it's address does not match filter '%s'...", dev, dev.ID)
				continue
			}
			logger.Infof(ctx, "found device %s; initializing...", dev)

			err := dev.Init(ctx)
			if err != nil {
				logger.Fatalf(ctx, "unable to initialize the connection to the device: %v", err)
			}
			logger.Infof(ctx, "requesting to pair")
			err = dev.Pairer().Pair(ctx)
			if err != nil {
				logger.Fatalf(ctx, "unable to pair: %v", err)
			}
			/*
				logger.Infof(ctx, "just in case sending to stop streaming")
				err = dev.Streamer().StopLiveStream(ctx)
				if err != nil {
					logger.Fatalf(ctx, "unable to request the device to stop live stream: %v", err)
				}
			*/
			logger.Infof(ctx, "prepare to live stream")
			err = dev.Streamer().PrepareToLiveStream(ctx)
			if err != nil {
				logger.Fatalf(ctx, "unable to request the device to prepare to live stream: %v", err)
			}
			logger.Infof(ctx, "requesting to connect to WiFi")
			err = dev.Pairer().ConnectToWiFi(ctx, *wifiSSID, *wifiPSK)
			if err != nil {
				logger.Fatalf(ctx, "unable to make the device connect to our WiFi: %v", err)
			}
			switch dev.Type {
			case djictl.DeviceTypeOsmoAction4, djictl.DeviceTypeOsmoAction5Pro:
				logger.Infof(ctx, "set image stabilization")
				err = dev.Configurer().SetImageStabilization(ctx, djictl.ImageStabilizationRockSteadyPlus)
				if err != nil {
					logger.Fatalf(ctx, "unable to set image stabilization to RockSteadyPlus: %v", err)
				}
			}
			logger.Infof(ctx, "start live stream")
			err = dev.Streamer().LiveStream(ctx, djictl.Resolution1080p, 5000, djictl.FPS25, *rtmpURL)
			if err != nil {
				logger.Fatalf(ctx, "unable to make the device to stream: %v", err)
			}
			logger.Infof(ctx, "done")
		case err := <-errCh:
			logger.Fatalf(ctx, "%v", err)
		}
	}
}
