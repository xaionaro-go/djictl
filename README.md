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
xaionaro@void:/home/streaming/go/src/github.com/xaionaro-go/djictl$ ./build/djictl-linux-amd64 --help
Usage of ./build/djictl-linux-amd64:
      --filter-device-addr string   
      --log-level Level             Log level (default info)
      --rtmp-url string             
      --wifi-psk string             
      --wifi-ssid string            
pflag: help requested
```

Let's start a stream to our server:
```sh
sudo ./build/djictl-linux-amd64 --wifi-ssid '<MY-WIFI-SSID>' --wifi-psk '<MY-WIFI-PSK>' --rtmp-url 'rtmp://MY_HOST/live/stream'
```

If it does not work, create a ticket.

## Reverse engineering

The reverse engineering was done here:
* [github.com/xaionaro/reverse-engineering-dji](https://github.com/xaionaro/reverse-engineering-dji). Feel free to continue the research :)

## See also:

* A port to Qt (C++): [github.com/xaionaro/libdji](https://github.com/xaionaro/libdji)
