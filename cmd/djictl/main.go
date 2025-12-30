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

			err := connectWiFiAndStartStreaming(ctx, dev, *wifiSSID, *wifiPSK, *rtmpURL)
			if err != nil {
				logger.Fatalf(ctx, "%v", err)
			}
		case err := <-errCh:
			logger.Fatalf(ctx, "%v", err)
		}
	}
}
