package bluetoothmanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/godbus/dbus/v5"
)

type BluetoothScanner struct {
	Conn     *dbus.Conn
	Devices  map[string]Device
	Listener func(Device)
}

func NewBluetoothScanner(listener func(Device)) (*BluetoothScanner, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SystemBus: %w", err)
	}

	return &BluetoothScanner{
		Conn:     conn,
		Devices:  make(map[string]Device),
		Listener: listener,
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
		var d Device
		if valid := d.ParseDevice(path, props); valid {
			bs.Listener(d)
			bs.Devices[d.MacAddress] = d
		}
	}

	return nil
}

func (bs *BluetoothScanner) StartScanner() (err error) {
	adapterPath := dbus.ObjectPath(filepath.Join(string(os.PathSeparator), "org", "bluez", "hci0"))
	signals := make(chan *dbus.Signal, 10)
	bs.Conn.Signal(signals)

	matchRule := "type='signal',interface='org.freedesktop.DBus.ObjectManager',member='InterfacesAdded'"
	if call := bs.Conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule); call.Err != nil {
		return call.Err
	}

	adapter := bs.Conn.Object("org.bluez", adapterPath)
	err = adapter.Call("org.bluez.Adapter1.StartDiscovery", 0).Store()
	if err != nil {
		return err
	}
	PrintDebug("Scanning for devices...")
	go bs.HandleSignals(signals)

	return nil
}

func (bs *BluetoothScanner) HandleSignals(signals chan *dbus.Signal) {
	for signal := range signals {
		bs.HandleSignal(signal)
	}
}

func (bs *BluetoothScanner) HandleSignal(signal *dbus.Signal) {
	if signal.Name != "org.freedesktop.DBus.ObjectManager.InterfacesAdded" {
		return
	}
	path, properties := signal.Body[0].(dbus.ObjectPath), signal.Body[1].(map[string]map[string]dbus.Variant)
	for _, props := range properties {
		var d Device
		if valid := d.ParseDevice(path, props); valid {
			bs.Listener(d)
			bs.Devices[d.MacAddress] = d
		}
	}
}
