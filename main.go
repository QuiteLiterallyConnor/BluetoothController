package main

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
)

type BluetoothController struct {
	DeviceName string
	DeviceMAC  string
	Conn       *dbus.Conn
}

func NewBluetoothController(deviceName, deviceMAC string) (*BluetoothController, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SystemBus: %w", err)
	}
	return &BluetoothController{
		DeviceName: deviceName,
		DeviceMAC:  deviceMAC,
		Conn:       conn,
	}, nil
}

func (bc *BluetoothController) ListenForPropertyChanges() {
	matchRule := "type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged'"
	bc.Conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)

	c := make(chan *dbus.Signal, 10)
	bc.Conn.Signal(c)

	fmt.Println("Listening for PropertiesChanged signals...")
	for v := range c {
		bc.onPropertiesChanged(v)
	}
}

func (bc *BluetoothController) onPropertiesChanged(signal *dbus.Signal) {
	if len(signal.Body) >= 3 {
		interfaceName := signal.Body[0].(string)
		changedProperties := signal.Body[1].(map[string]dbus.Variant)

		fmt.Println("PropertiesChanged on interface:", interfaceName)
		for propName, propValue := range changedProperties {
			fmt.Printf("Property %s changed to %v\n", propName, propValue)
		}
	}
}

func (bc *BluetoothController) ControlMedia(action string) {
	fmt.Printf("controlMedia: %s\n", action)
	deviceMACFormatted := strings.ToUpper(strings.Replace(bc.DeviceMAC, ":", "_", -1))
	mediaPlayerPath := fmt.Sprintf("/org/bluez/hci0/dev_%s/player0", deviceMACFormatted)

	mediaPlayer := bc.Conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	call := mediaPlayer.Call("org.bluez.MediaPlayer1."+strings.Title(action), 0)
	if call.Err != nil {
		fmt.Fprintf(os.Stderr, "Failed to %s: %s\n", strings.ToLower(action), call.Err)
		return
	}
	fmt.Printf("%s action executed for %s\n", action, bc.DeviceName)
}


func main() {
	bluetootheventlistener = BluetoothController{}
	bc, err := bluetootheventlistener.NewBluetoothController("Pixel_6", "0C:C4:13:12:67:62")
	if err != nil {
		fmt.Println("Error creating BluetoothController:", err)
		return
	}

	go bc.ListenForPropertyChanges()
	bc.ControlMedia("play")
}
