package bluetoothmanager

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/godbus/dbus/v5"
)

var (
	Debug = false
)

type Device struct {
	AdapterPath dbus.ObjectPath
	MacAddress  string
	Alias       string
	Blocked     bool
	Connected   bool
	Name        string
	Paired      bool
	RSSI        int16
	Trusted     bool
}

func EnableDebugging() {
	Debug = true
}

func PrintDebug(message string) {
	if Debug {
		fmt.Println(message)
	}
}

func parseDevice(path dbus.ObjectPath, props map[string]dbus.Variant) (device Device) {
	device.AdapterPath = path

	if address, ok := props["Address"]; ok {
		device.MacAddress = address.Value().(string)
	}

	if name, ok := props["Name"]; ok {
		device.Name = name.Value().(string)
	} else {
		device.Name = "Unknown"
	}

	if alias, ok := props["Alias"]; ok {
		device.Alias = alias.Value().(string)
		if device.Name == "Unknown" {
			device.Name = device.Alias
		}
	}

	if blocked, ok := props["Blocked"]; ok {
		device.Blocked = blocked.Value().(bool)
	}

	if connected, ok := props["Connected"]; ok {
		device.Connected = connected.Value().(bool)
	}

	if paired, ok := props["Paired"]; ok {
		device.Paired = paired.Value().(bool)
	}

	if rssi, ok := props["RSSI"]; ok {
		device.RSSI = rssi.Value().(int16)
	}

	if trusted, ok := props["Trusted"]; ok {
		device.Trusted = trusted.Value().(bool)
	}

	return
}

func extractMACAddress(input string) (match string) {
	pattern := `([0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2})`
	re, _ := regexp.Compile(pattern)
	if match = re.FindString(input); match == "" {
		return "Unknown"
	}
	return strings.ReplaceAll(match, "_", ":")
}
