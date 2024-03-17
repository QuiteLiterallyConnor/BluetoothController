package main

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
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
	DeviceName       string
	DeviceMacAddress string
	Conn             *dbus.Conn
	Listener         func(string, string, interface{}, reflect.Type)
	Broadcaster      func(string, interface{}, reflect.Type)
}

// TODO: Add overloading such that NewBluetoothController can be called with or without the device_name and device_address parameters

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

	mac_address := bc.extractMACAddress(string(signal.Path))

	for event_name, prop := range signal.Body[1].(map[string]dbus.Variant) {
		value := prop.Value()
		typeof := reflect.TypeOf(prop.Value())

		fmt.Printf("sender address: %v		name: %+v	value: %+v	typeof: %+v\n", mac_address, event_name, value, typeof)
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

func (bc *BluetoothController) ControlMedia(action, mac_address string) error { // Adjusted to return an error
	mac_address = strings.Replace(mac_address, ":", "_", -1)
	fmt.Printf("action: %v 	mac_address: %v\n", action, mac_address)
	mediaPlayerPath := fmt.Sprintf("/org/bluez/hci0/dev_%s/player0", mac_address)

	fmt.Printf("mediaPlayerPath: %v\n", mediaPlayerPath)

	mediaPlayer := bc.Conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	call := mediaPlayer.Call("org.bluez.MediaPlayer1."+action, 0)
	if call.Err != nil {
		return fmt.Errorf("failed to %s: %w", strings.ToLower(action), call.Err) // Error wrapping for better handling
	}

	fmt.Printf("%s action executed for %s\n", action, bc.DeviceName)

	return nil
}

func functionListen(device_address, command_name string, value interface{}, typeof reflect.Type) {
	fmt.Printf("device_address: %v		command_name: %v		value: %v		typeof: %v\n", device_address, command_name, value, typeof)
}

func listenForCommand(bc *BluetoothController) {
	for {
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		parts := strings.Split(input, "::")
		if len(parts) != 2 {
			fmt.Fprintf(os.Stderr, "Invalid input format. Expected format is 'NAME:ADDRESS'.\n")
			continue
		}

		name := parts[0]
		address := parts[1]

		fmt.Printf("Command: %s, Address: %s\n", name, address)

		if err := bc.ControlMedia(name, address); err != nil {
			fmt.Fprintf(os.Stderr, "Error controlling media: %s\n", err)
		}
	}
}

func main() {
	bc, err := NewBluetoothController(functionListen)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing BluetoothController: %s\n", err)
		os.Exit(1)
	}

	bc.Start()

	listenForCommand(bc)

}
