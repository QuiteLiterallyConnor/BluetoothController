package bluetoothmanager

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

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

func (d *Device) ParseDevice(path dbus.ObjectPath, props map[string]dbus.Variant) (valid bool) {
	valid = true
	d.AdapterPath = path

	if address, ok := props["Address"]; ok {
		d.MacAddress = address.Value().(string)
	} else {
		valid = false
	}

	if alias, ok := props["Alias"]; ok {
		d.Alias = alias.Value().(string)
	}

	if name, ok := props["Name"]; ok {
		d.Name = name.Value().(string)
	} else {
		d.Name = d.MacAddress
	}

	if blocked, ok := props["Blocked"]; ok {
		d.Blocked = blocked.Value().(bool)
	}

	if connected, ok := props["Connected"]; ok {
		d.Connected = connected.Value().(bool)
	}

	if paired, ok := props["Paired"]; ok {
		d.Paired = paired.Value().(bool)
	}

	if rssi, ok := props["RSSI"]; ok {
		d.RSSI = rssi.Value().(int16)
	}

	if trusted, ok := props["Trusted"]; ok {
		d.Trusted = trusted.Value().(bool)
	}

	return
}

func (d *Device) Connect() error {
	if d.Connected {
		return nil
	}

	if !d.Paired {
		return nil
	}

	if err := d.ConnectToDevice(); err != nil {
		return fmt.Errorf("failed to establish connection: %w", err)
	}

	timeout := time.After(10 * time.Second)
	tick := time.Tick(500 * time.Millisecond)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("connection timeout")
		case <-tick:
			if d.Connected {
				fmt.Println("Connected successfully.")
				return nil
			}
		}
	}
}

func (d *Device) ConnectToDevice() error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return fmt.Errorf("connecting to D-Bus system bus failed: %w", err)
	}

	obj := conn.Object("org.bluez", dbus.ObjectPath(d.AdapterPath))

	call := obj.Call("org.bluez.Device1.Connect", 0)
	if call.Err != nil {
		return fmt.Errorf("connecting to the Bluetooth device failed: %w", call.Err)
	}

	PrintDebug("Connected successfully")
	return nil
}

type Event struct {
	Device    string
	Category  string
	Value     string
	ValueType reflect.Type
}

func (e *Event) ParseEvent(event_name, address string, prop dbus.Variant) {
	e.Device = extractMACAddress(address)
	e.Category = event_name
	e.ValueType = reflect.TypeOf(prop.Value())
	e.Value = prop.String()
}

func (e *Event) Json() string {
	jsonString, _ := json.Marshal(e)
	return string(jsonString)
}

func EnableDebugging() {
	Debug = true
}

func PrintDebug(message string) {
	if Debug {
		fmt.Println(message)
	}
}

func extractMACAddress(input string) (match string) {
	pattern := `([0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2}[:_][0-9A-Fa-f]{2})`
	re, _ := regexp.Compile(pattern)
	if match = re.FindString(input); match == "" {
		return "Unknown"
	}
	return strings.ReplaceAll(match, "_", ":")
}
