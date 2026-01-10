package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/facebookincubator/go-belt/tool/logger"
	"github.com/urfave/cli/v2"
	"github.com/xaionaro-go/djictl/pkg/djible"
	"github.com/xaionaro-go/djictl/pkg/djiwifi"
	"github.com/xaionaro-go/djictl/pkg/duml"
)

var errDone = fmt.Errorf("done")

func main() {
	app := &cli.App{
		Name:  "djictl",
		Usage: "DJI Osmo devices control tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "log-level",
				Value: "info",
				Usage: "Log level (debug, info, warn, error, fatal, panic)",
			},
			&cli.StringFlag{
				Name:  "filter-device-addr",
				Value: "",
				Usage: "Filter device by address",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "ble",
				Usage: "BLE-based commands",
				Subcommands: []*cli.Command{
					{
						Name:  "scan",
						Usage: "Scan for DJI devices",
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								fmt.Printf("Found device: %s\n", dev)
								return nil
							})
						},
					},
					{
						Name:  "connect-wifi-and-start-streaming",
						Usage: "Connect device to WiFi and start RTMP streaming",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "wifi-ssid",
								Usage:    "WiFi SSID",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "wifi-psk",
								Usage:    "WiFi Password",
								Required: true,
							},
							&cli.StringFlag{
								Name:     "rtmp-url",
								Usage:    "RTMP URL",
								Required: true,
							},
							&cli.StringFlag{
								Name:  "resolution",
								Usage: "Video resolution (allowed values: 480p, 720p, 1080p)",
								Value: "1080p",
							},
							&cli.UintFlag{
								Name:  "bitrate-kbps",
								Usage: "bitrate in Kbps",
								Value: 6000,
							},
							&cli.UintFlag{
								Name:  "fps",
								Usage: "frames per second (allowed values: 25, 30)",
								Value: 30,
							},
						},
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								resolution := duml.ResolutionFromString(c.String("resolution"))
								if resolution == duml.UndefinedResolution {
									return fmt.Errorf("invalid resolution value %q", c.String("resolution"))
								}
								fps := duml.FPSFromUint(uint(c.Uint("fps")))
								if fps == duml.UndefinedFPS {
									return fmt.Errorf("invalid fps value %d", c.Uint("fps"))
								}
								return connectWiFiAndStartStreaming(
									ctx,
									dev,
									c.String("wifi-ssid"),
									c.String("wifi-psk"),
									c.String("rtmp-url"),
									resolution,
									uint16(c.Uint("bitrate-kbps")),
									fps,
								)
							})
						},
					},
					{
						Name:  "camera-ap-info",
						Usage: "Get camera AP SSID and Password [does not work, yet]",
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								err := dev.Init(ctx)
								if err != nil {
									return fmt.Errorf("unable to initialize: %w", err)
								}
								err = dev.AppToWiFiGroundStation().Pair(ctx)
								if err != nil {
									return fmt.Errorf("unable to pair: %w", err)
								}
								ssid, psk, err := dev.AppToWiFiGroundStation().CameraAPInfo(ctx)
								if err != nil {
									return fmt.Errorf("unable to get camera AP info: %w", err)
								}
								fmt.Printf("SSID: %s\nPSK: %s\n", ssid, psk)
								return errDone
							})
						},
					},
					{
						Name:  "fcc-enable",
						Usage: "Enable FCC mode [does not work, yet]",
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								err := dev.Init(ctx)
								if err != nil {
									return err
								}
								return dev.AppToCamera().SetFCCEnable(ctx, true)
							})
						},
					},
					{
						Name:  "set-goggles-mode",
						Usage: "Set Goggles mode [does not work, yet]",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "mode",
								Value: "usb",
								Usage: "Mode: usb or normal",
							},
						},
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								err := dev.Init(ctx)
								if err != nil {
									return err
								}
								var mode duml.GogglesMode
								switch strings.ToLower(c.String("mode")) {
								case "usb":
									mode = duml.GogglesModeUSB
								case "normal":
									mode = duml.GogglesModeNormal
								default:
									return fmt.Errorf("invalid mode: %s", c.String("mode"))
								}
								return dev.AppToGoggles().SetMode(ctx, mode)
							})
						},
					},
					{
						Name:  "remote-controller-simulator",
						Usage: "Send Remote Controller simulator data [does not work, yet]",
						Flags: []cli.Flag{
							&cli.IntFlag{Name: "right-h", Value: 1024},
							&cli.IntFlag{Name: "right-v", Value: 1024},
							&cli.IntFlag{Name: "left-v", Value: 1024},
							&cli.IntFlag{Name: "left-h", Value: 1024},
						},
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								err := dev.Init(ctx)
								if err != nil {
									return err
								}
								data := duml.RemoteControllerSimulatorData{
									RightStickHorizontal: uint16(c.Int("right-h")),
									RightStickVertical:   uint16(c.Int("right-v")),
									LeftStickVertical:    uint16(c.Int("left-v")),
									LeftStickHorizontal:  uint16(c.Int("left-h")),
								}
								return dev.AppToRemoteController().SendData(ctx, data)
							})
						},
					},
					{
						Name:  "rtmp-broadcast",
						Usage: "Configure RTMP broadcast [does not work, yet]",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "url", Usage: "RTMP URL", Required: true},
							&cli.BoolFlag{Name: "disable", Usage: "Disable (instead of enable)"},
						},
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								err := dev.Init(ctx)
								if err != nil {
									return err
								}
								return dev.AppToVideoTransmission().ConfigureRTMP(ctx, c.String("url"), !c.Bool("disable"))
							})
						},
					},
					{
						Name:  "battery-info",
						Usage: "Request battery information",
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								err := dev.Init(ctx)
								if err != nil {
									return err
								}
								status, err := dev.AppToBattery().GetInfo(ctx)
								if err != nil {
									return err
								}
								fmt.Printf("Battery capacity: %s\n", status.Capacity)
								return nil
							})
						},
					},
					{
						Name:  "firmware-version",
						Usage: "Request firmware version [does not work, yet]",
						Action: func(c *cli.Context) error {
							return runOnBLE(c, func(ctx context.Context, dev *djible.Device) error {
								err := dev.Init(ctx)
								if err != nil {
									return err
								}
								return dev.AppToCamera().GetVersion(ctx)
							})
						},
					},
				},
			},
			{
				Name:  "wifi",
				Usage: "WiFi-based commands (UDP 9004)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "addr",
						Value: "192.168.2.1:9004",
						Usage: "Device UDP address",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "start-video",
						Usage: "Start video via WiFi [does not work, yet]",
						Action: func(c *cli.Context) error {
							return runOnWiFi(c, func(ctx context.Context, ctrl *djiwifi.Controller) error {
								if err := ctrl.SendHandshake(ctx); err != nil {
									return err
								}
								return ctrl.SendVideoHandshake(ctx)
							})
						},
					},
					{
						Name:  "fcc-enable",
						Usage: "Enable FCC mode via WiFi [does not work, yet]",
						Action: func(c *cli.Context) error {
							return runOnWiFi(c, func(ctx context.Context, ctrl *djiwifi.Controller) error {
								return ctrl.SendFCCEnable(ctx, true)
							})
						},
					},
					{
						Name:  "rtmp-broadcast",
						Usage: "Configure RTMP broadcast via WiFi [does not work, yet]",
						Flags: []cli.Flag{
							&cli.StringFlag{Name: "url", Usage: "RTMP URL", Required: true},
							&cli.BoolFlag{Name: "disable", Usage: "Disable (instead of enable)"},
						},
						Action: func(c *cli.Context) error {
							return runOnWiFi(c, func(ctx context.Context, ctrl *djiwifi.Controller) error {
								return ctrl.SendConfigureBroadcast(ctx, c.String("url"), !c.Bool("disable"))
							})
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		if err == errDone {
			return
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runOnBLE(c *cli.Context, action func(ctx context.Context, dev *djible.Device) error) error {
	var loggerLevel logger.Level
	if err := loggerLevel.Set(c.String("log-level")); err != nil {
		return fmt.Errorf("invalid log level '%s': %w", c.String("log-level"), err)
	}

	ctx := getContext(loggerLevel, false, "")
	logger.Debugf(ctx, "log level: %s (raw value: '%s')", loggerLevel, c.String("log-level"))
	filterDeviceAddr := c.String("filter-device-addr")

	devCh, errCh, err := djible.Scan(ctx)
	if err != nil {
		return fmt.Errorf("unable to start scanning: %w", err)
	}

	for {
		select {
		case dev := <-devCh:
			logger.Debugf(ctx, "found device %s", dev)
			if !strings.Contains(strings.ToLower(dev.ID.String()), strings.ToLower(filterDeviceAddr)) {
				logger.Infof(ctx, "found device %s; but skipping, because it's address does not match filter '%s'...", dev, filterDeviceAddr)
				continue
			}

			if err := action(ctx, dev); err != nil {
				return err
			}
		case err := <-errCh:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func runOnWiFi(c *cli.Context, action func(ctx context.Context, ctrl *djiwifi.Controller) error) error {
	var loggerLevel logger.Level
	if err := loggerLevel.Set(c.String("log-level")); err != nil {
		return fmt.Errorf("invalid log level '%s': %w", c.String("log-level"), err)
	}

	ctx := getContext(loggerLevel, false, "")
	logger.Debugf(ctx, "log level: %s (raw value: '%s')", loggerLevel, c.String("log-level"))
	addr := c.String("addr")

	ctrl, err := djiwifi.NewController(ctx, addr)
	if err != nil {
		return err
	}
	defer ctrl.Close()

	return action(ctx, ctrl)
}
