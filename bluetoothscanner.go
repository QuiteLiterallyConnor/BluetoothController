package bluetoothmanager

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

type BluetoothScanner struct {
	Conn    *dbus.Conn
	Devices map[string]Device
}

func NewBluetoothScanner() (*BluetoothScanner, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SystemBus: %w", err)
	}

	return &BluetoothScanner{
		Conn:    conn,
		Devices: make(map[string]Device),
	}, nil
}

func (bs *BluetoothScanner) GetDevices() map[string]Device {
	return bs.Devices
}

func (bs *BluetoothScanner) GetManagedDevices() error {
	obj := bs.Conn.Object("org.bluez", "/")
	managedObjects := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	call := obj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0)
	if call.Err != nil {
		return call.Err
	}
	call.Store(&managedObjects)

	for path, interfaces := range managedObjects {
		if _, ok := interfaces["org.bluez.Device1"]; !ok {
			continue
		}
		props, found := interfaces["org.bluez.Device1"]
		if !found {
			continue
		}
		device := parseDevice(path, props)
		bs.Devices[device.Name] = device
	}

	return nil
}

func (bs *BluetoothScanner) ListenForDevices() {
	PrintDebug("Listening for devices")
	matchRule := "type='signal',interface='org.freedesktop.DBus.ObjectManager',member='InterfacesAdded',sender='org.bluez'"
	call := bs.Conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)
	if call.Err != nil {
		fmt.Printf("Error adding D-Bus match rule: %s\n", call.Err)
		return
	}

	signalChan := make(chan *dbus.Signal, 10)
	bs.Conn.Signal(signalChan)

	go func() {
		for signal := range signalChan {
			fmt.Printf("signal: %+v\n", signal)

			if signal.Name != "org.freedesktop.DBus.ObjectManager.InterfacesAdded" {
				continue
			}
			if path, props, ok := parseInterfacesAddedSignal(signal); ok {
				device := parseDevice(path, props["org.bluez.Device1"])
				bs.Devices[device.Name] = device
			}
		}
	}()
}

func parseInterfacesAddedSignal(signal *dbus.Signal) (dbus.ObjectPath, map[string]map[string]dbus.Variant, bool) {
	if len(signal.Body) < 2 {
		return "", nil, false
	}

	path, ok1 := signal.Body[0].(dbus.ObjectPath)
	properties, ok2 := signal.Body[1].(map[string]map[string]dbus.Variant)

	if !ok1 || !ok2 {
		return "", nil, false
	}

	return path, properties, true
}
