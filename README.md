# SteelSeries Siberia v2 headset LEDs control tool

## Why?

SteelSeries software supports Windows and MacOSX.
This tool is able control headset's LEDs in Linux.

## Supported hardware

I tested SteelSeries Siberia v2 Frost Blue only.

## Requirements

* Go compiler
* libusb >= 1.0

## Instalation

```
git clone git://github.com/antage/ssv2leds.git
cd ssv2leds
go get github.com/hanwen/go-mtpfs/usb
go build
sudo install -m 0644 ssv2leds /usr/local/bin/ssv2leds
```

## Usage

Turn off LEDs:
```
ssv2leds -i 0
```

Maximal brightness:
```
ssv2leds -i 255
```

Medium brightness:
```
ssv2leds -i 128
```

Set pulsation mode:
```
ssv2leds -p slow
```


Pulsation modes:

* steady - LEDs are on always
* slow - slow pulsation
* medium - medium pulastion
* fast - fast pulsation
* trigger - I don't know

