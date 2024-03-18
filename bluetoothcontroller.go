package bluetoothmanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/godbus/dbus/v5"
)

type BluetoothController struct {
	Adapter      dbus.BusObject
	AdapterPath  dbus.ObjectPath
	Conn         *dbus.Conn
	Listener     func(Event)
	ActiveDevice Device
}

func NewBluetoothController(listener func(Event)) (*BluetoothController, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SystemBus: %w", err)
	}
	adapterPath := filepath.Join(string(os.PathSeparator), "org", "bluez", "hci0")
	return &BluetoothController{
		Adapter:     conn.Object("org.bluez", dbus.ObjectPath(adapterPath)),
		AdapterPath: dbus.ObjectPath(adapterPath),
		Conn:        conn,
		Listener:    listener,
	}, nil
}

func (bc *BluetoothController) GetActiveDevice() Device {
	return bc.ActiveDevice
}

func (bc *BluetoothController) StartController() {
	go bc.ListenForPropertyChanges()
}

func (bc *BluetoothController) ListenForPropertyChanges() {
	PrintDebug("Listening for property changes")
	matchRule := "type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged'"
	bc.Conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)
	c := make(chan *dbus.Signal, 10)
	bc.Conn.Signal(c)
	for v := range c {
		bc.onPropertiesChanged(v)
	}
}

func (bc *BluetoothController) onPropertiesChanged(signal *dbus.Signal) {
	bc.UpdateActiveDevice()

	if len(signal.Body) < 3 {
		return
	}
	mac_address := extractMACAddress(string(signal.Path))
	for event_name, prop := range signal.Body[1].(map[string]dbus.Variant) {
		var e Event
		e.ParseEvent(event_name, mac_address, prop)
		bc.Listener(e)
	}
}

func (bc *BluetoothController) ControlMedia(action, mac_address string) error {
	mac_address = strings.Replace(mac_address, ":", "_", -1)
	mediaPlayerPath := filepath.Join(string(bc.AdapterPath), fmt.Sprintf("dev_%s", mac_address), "player0")
	PrintDebug(fmt.Sprintf("mediaPlayerPath: %v\n", mediaPlayerPath))
	mediaPlayer := bc.Conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	PrintDebug(fmt.Sprintf("Calling %s on %s", action, mediaPlayerPath))
	if call := mediaPlayer.Call("org.bluez.MediaPlayer1."+action, 0); call.Err != nil {
		return fmt.Errorf("failed to %s: %w", strings.ToLower(action), call.Err)
	}
	return nil
}

func (bc *BluetoothController) UpdateActiveDevice() error {
	devicePaths, err := bc.getConnectedDevices()
	if err != nil {
		return err
	}

	for _, devicePath := range devicePaths {
		deviceProps, err := bc.getDeviceProperties(devicePath)
		if err != nil {
			PrintDebug(fmt.Sprintf("Error getting properties for device %s: %v", devicePath, err))
			continue
		}

		var device Device
		if !(device.ParseDevice(devicePath, deviceProps) && device.Connected) {
			continue
		}

		mediaPlayerProps, err := bc.getMediaPlayerProperties(devicePath)
		if err != nil {
			continue
		}
		status, exists := mediaPlayerProps["Status"]
		if exists && status.Value().(string) == "playing" {
			bc.ActiveDevice = device
			PrintDebug("Active device updated: " + device.MacAddress)
			return nil
		}
	}

	return fmt.Errorf("no active playing device found")
}

func (bc *BluetoothController) getConnectedDevices() ([]dbus.ObjectPath, error) {
	var connectedDevices []dbus.ObjectPath
	managedObjects := make(map[dbus.ObjectPath]map[string]map[string]dbus.Variant)
	if err := bc.Conn.Object("org.bluez", "/").Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&managedObjects); err != nil {
		return nil, fmt.Errorf("failed to get managed objects: %w", err)
	}
	for path, interfaces := range managedObjects {
		deviceProps, exists := interfaces["org.bluez.Device1"]
		if !exists {
			continue
		}

		if connected, ok := deviceProps["Connected"]; ok && connected.Value().(bool) {
			connectedDevices = append(connectedDevices, path)
		}
	}
	return connectedDevices, nil
}

func (bc *BluetoothController) getDeviceProperties(devicePath dbus.ObjectPath) (map[string]dbus.Variant, error) {
	device := bc.Conn.Object("org.bluez", devicePath)
	var properties map[string]dbus.Variant
	err := device.Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.bluez.Device1").Store(&properties)
	return properties, err
}

func (bc *BluetoothController) getMediaPlayerProperties(devicePath dbus.ObjectPath) (map[string]dbus.Variant, error) {
	mediaPlayerPath := filepath.Join(string(devicePath), "player0")
	mediaPlayer := bc.Conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	var properties map[string]dbus.Variant
	err := mediaPlayer.Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.bluez.MediaPlayer1").Store(&properties)
	return properties, err
}

func (bc *BluetoothController) ConnectToDevice(device Device) error {
	conn, err := dbus.SystemBus()
	if err != nil {
		return fmt.Errorf("connecting to D-Bus system bus failed: %w", err)
	}

	obj := conn.Object("org.bluez", dbus.ObjectPath(device.AdapterPath))

	call := obj.Call("org.bluez.Device1.Connect", 0)
	if call.Err != nil {
		return fmt.Errorf("connecting to the Bluetooth device failed: %w", call.Err)
	}

	PrintDebug("Connected successfully")
	return nil
}
