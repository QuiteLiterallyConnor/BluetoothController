# Bluetooth Controller

This package provides a Go-based solution for controlling Bluetooth devices, specifically focusing on media control functionalities such as Play, Pause, and Stop.

## Features

- Listen for `PropertiesChanged` signals from Bluetooth devices
- Control media playback on Bluetooth devices that support the `org.bluez.MediaPlayer1` interface
- Display RSSI (Received Signal Strength Indicator) values for connected devices

## Installation

Ensure you have Go installed on your machine. Then, clone this repository and navigate to the project directory:

```bash
git clone https://github.com/QuiteLiterallyConnor/BluetoothController.git
cd BluetoothController
```

## Usage

An example implementation can be found in `examples`

## Troubleshooting

Q: Does my device support dbus/gdbus?
    A: Using dbus-send: dbus-send --system --print-reply --dest=org.bluez /org/bluez/hci0/dev_YOUR_DEVICE_MAC_ADDRESS org.freedesktop.DBus.Introspectable.Introspect
       Using gdbus:     gdbus introspect --system --dest org.bluez --object-path /org/bluez/hci0/dev_YOUR_DEVICE_MAC_ADDRESS