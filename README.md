# `djictl`

[![License: CC0-1.0](https://img.shields.io/badge/License-CC0%201.0-lightgrey.svg)](http://creativecommons.org/publicdomain/zero/1.0/)

An open-source CLI to manage your DJI Osmo device via BLE and without DJI MIMO.

Right now the only thing that it can do is to force your DJI device to stream to a specific RTMP destination. Thus work-in-progress (to add more capabilities).

## Quick start

```sh
make
./build/djictl-linux-amd64 --help
```

expected output:
```
streaming@void:~/go/src/github.com/xaionaro-go/djictl$ build/djictl-linux-amd64 --help
NAME:
   djictl - DJI Osmo devices control tool

USAGE:
   djictl [global options] command [command options]

COMMANDS:
   ble      BLE-based commands
   wifi     WiFi-based commands (UDP 9004)
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --log-level value           Log level (debug, info, warn, error, fatal, panic) (default: "info")
   --filter-device-addr value  Filter device by address
   --help, -h                  show help
```

```sh
./build/djictl-linux-amd64 ble --help
```
expected output:
```
streaming@void:~/go/src/github.com/xaionaro-go/djictl$ build/djictl-linux-amd64 ble --help
NAME:
   djictl ble - BLE-based commands

USAGE:
   djictl ble [command options]

COMMANDS:
   scan                              Scan for DJI devices
   connect-wifi-and-start-streaming  Connect device to WiFi and start RTMP streaming
   camera-ap-info                    Get camera AP SSID and Password [does not work, yet]
   fcc-enable                        Enable FCC mode [does not work, yet]
   set-goggles-mode                  Set Goggles mode [does not work, yet]
   remote-controller-simulator       Send Remote Controller simulator data [does not work, yet]
   rtmp-broadcast                    Configure RTMP broadcast [does not work, yet]
   battery-info                      Request battery information [does not work, yet]
   firmware-version                  Request firmware version [does not work, yet]
   help, h                           Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

Let's start a stream to our server:
```sh
sudo ./build/djictl-linux-amd64 ble connect-wifi-and-start-streaming --wifi-ssid '<MY-WIFI-SSID>' --wifi-psk '<MY-WIFI-PSK>' --rtmp-url 'rtmp://MY_HOST/live/stream'
```

If it does not work, create a ticket.

## Reverse engineering

The reverse engineering was done here:
* [github.com/xaionaro/reverse-engineering-dji](https://github.com/xaionaro/reverse-engineering-dji). Feel free to continue the research :)

## See also:

* A port to Qt (C++): [github.com/xaionaro/libdji](https://github.com/xaionaro/libdji)
