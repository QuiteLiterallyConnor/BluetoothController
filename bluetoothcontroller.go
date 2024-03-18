package bluetoothmanager

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/godbus/dbus/v5"
)

type BluetoothController struct {
	Adapter     dbus.BusObject
	AdapterPath dbus.ObjectPath
	Conn        *dbus.Conn
	Listener    func(string, string, interface{}, reflect.Type)
}

func NewBluetoothController(listener func(string, string, interface{}, reflect.Type)) (*BluetoothController, error) {
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
	PrintDebug(fmt.Sprintf("Received signal: %v", signal))
	if len(signal.Body) < 3 {
		return
	}
	mac_address := extractMACAddress(string(signal.Path))
	for event_name, prop := range signal.Body[1].(map[string]dbus.Variant) {
		value := prop.Value()
		typeof := reflect.TypeOf(prop.Value())
		PrintDebug(fmt.Sprintf("MAC: %s, Event: %s, Value: %v, Type: %v", mac_address, event_name, value, typeof))
		bc.Listener(mac_address, event_name, value, typeof)
	}
}

func (bc *BluetoothController) ControlMedia(action, mac_address string) error {
	mac_address = strings.Replace(mac_address, ":", "_", -1)
	mediaPlayerPath := filepath.Join(string(bc.AdapterPath), fmt.Sprintf("dev_%s", mac_address), "player0")
	fmt.Printf("mediaPlayerPath: %v\n", mediaPlayerPath)
	mediaPlayer := bc.Conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	PrintDebug(fmt.Sprintf("Calling %s on %s", action, mediaPlayerPath))
	if call := mediaPlayer.Call("org.bluez.MediaPlayer1."+action, 0); call.Err != nil {
		return fmt.Errorf("failed to %s: %w", strings.ToLower(action), call.Err)
	}
	return nil
}
