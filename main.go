package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/godbus/dbus/v5"
)

type Event struct {
	Name     string
	Value    interface{}
	TypeOf   reflect.Type
	Function func(value interface{})
}

type BluetoothController struct {
	DeviceName  string
	DeviceMAC   string
	Conn        *dbus.Conn
	Listener    func(string, interface{}, reflect.Type)
	Broadcaster func(string, interface{}, reflect.Type)
}

func NewBluetoothController(deviceName, deviceMAC string, listener func(string, interface{}, reflect.Type)) (*BluetoothController, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SystemBus: %w", err)
	}

	return &BluetoothController{
		DeviceName: deviceName,
		DeviceMAC:  deviceMAC,
		Conn:       conn,
		Listener:   listener,
	}, nil

}

func (bc *BluetoothController) Start() {
	go bc.ListenForPropertyChanges()
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
	if len(signal.Body) < 3 {
		return
	}

	fmt.Printf("signal.Sender: %+v\n", signal.Sender)

	changedProperties := signal.Body[1].(map[string]dbus.Variant)

	for name, prop := range changedProperties {
		value := prop.Value()
		typeof := reflect.TypeOf(prop.Value())

		fmt.Printf("name: %+v	value: %+v	typeof: %+v\n", name, value, typeof)
		bc.Listener(name, value, typeof)
	}
}

func (bc *BluetoothController) ControlMedia(action string) error { // Adjusted to return an error
	deviceMACFormatted := strings.ToUpper(strings.Replace(bc.DeviceMAC, ":", "_", -1))
	mediaPlayerPath := fmt.Sprintf("/org/bluez/hci0/dev_%s/player0", deviceMACFormatted)

	mediaPlayer := bc.Conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	call := mediaPlayer.Call("org.bluez.MediaPlayer1."+action, 0)
	if call.Err != nil {
		return fmt.Errorf("failed to %s: %w", strings.ToLower(action), call.Err) // Error wrapping for better handling
	}

	fmt.Printf("%s action executed for %s\n", action, bc.DeviceName)

	return nil
}

func functionListen(name string, value interface{}, typeof reflect.Type) {
	fmt.Printf("Listener - Name: %v		Value: %v		Typeof: %v\n", name, value, typeof)
}

func listenForCommand(bc *BluetoothController) {
	for {
		var action string
		fmt.Scanln(&action)
		action = strings.TrimSpace(action)

		if action == "exit" {
			break
		}

		if err := bc.ControlMedia(action); err != nil {
			fmt.Fprintf(os.Stderr, "Error controlling media: %s\n", err)
		}
	}
}

func main() {
	var deviceName, deviceMAC string
	flag.StringVar(&deviceName, "name", "", "Name of the Bluetooth device")
	flag.StringVar(&deviceMAC, "mac_address", "", "MAC address of the Bluetooth device")
	flag.Parse()

	if deviceName == "" || deviceMAC == "" {
		fmt.Println("Both -name and -mac_address flags must be specified.")
		flag.PrintDefaults()
		os.Exit(1)
	}

	bc, err := NewBluetoothController(deviceName, deviceMAC, functionListen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing BluetoothController: %s\n", err)
		os.Exit(1)
	}
	bc.Start()

	listenForCommand(bc)

}
