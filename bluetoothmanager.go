package bluetoothmanager

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/godbus/dbus/v5"
)

type BluetoothController struct {
	Conn     *dbus.Conn
	Listener func(string, string, interface{}, reflect.Type)
	Debug    bool
}

func NewBluetoothController(listener func(string, string, interface{}, reflect.Type)) (*BluetoothController, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SystemBus: %w", err)
	}
	return &BluetoothController{
		Conn:     conn,
		Listener: listener,
	}, nil
}

func (bc *BluetoothController) Start() {
	go bc.ListenForPropertyChanges()
}

func (bc *BluetoothController) EnableDebugging() {
	bc.Debug = true
}

func (bc *BluetoothController) PrintDebug(message string) {
	if bc.Debug {
		fmt.Println(message)
	}
}

func (bc *BluetoothController) ListenForPropertyChanges() {
	bc.PrintDebug("Listening for property changes")
	matchRule := "type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged'"
	bc.Conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)
	c := make(chan *dbus.Signal, 10)
	bc.Conn.Signal(c)
	for v := range c {
		bc.onPropertiesChanged(v)
	}
}

func (bc *BluetoothController) onPropertiesChanged(signal *dbus.Signal) {
	bc.PrintDebug(fmt.Sprintf("Received signal: %v", signal))
	if len(signal.Body) < 3 {
		return
	}
	mac_address := bc.extractMACAddress(string(signal.Path))
	for event_name, prop := range signal.Body[1].(map[string]dbus.Variant) {
		value := prop.Value()
		typeof := reflect.TypeOf(prop.Value())
		bc.PrintDebug(fmt.Sprintf("MAC: %s, Event: %s, Value: %v, Type: %v", mac_address, event_name, value, typeof))
		bc.Listener(mac_address, event_name, value, typeof)
	}
}

func (bc *BluetoothController) extractMACAddress(input string) (match string) {
	pattern := `([0-9A-Fa-f]{2}_[0-9A-Fa-f]{2}_[0-9A-Fa-f]{2}_[0-9A-Fa-f]{2}_[0-9A-Fa-f]{2}_[0-9A-Fa-f]{2})`
	re, _ := regexp.Compile(pattern)
	if match = re.FindString(input); match == "" {
		return "Unknown"
	}
	return strings.Replace(match, "_", ":", -1)
}

func (bc *BluetoothController) ControlMedia(action, mac_address string) error {
	mac_address = strings.Replace(mac_address, ":", "_", -1)
	mediaPlayerPath := fmt.Sprintf("/org/bluez/hci0/dev_%s/player0", mac_address)
	mediaPlayer := bc.Conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	bc.PrintDebug(fmt.Sprintf("Calling %s on %s", action, mediaPlayerPath))
	if call := mediaPlayer.Call("org.bluez.MediaPlayer1."+action, 0); call.Err != nil {
		return fmt.Errorf("failed to %s: %w", strings.ToLower(action), call.Err)
	}
	return nil
}
