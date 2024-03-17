package bluetooth_manager

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
	var deviceName, deviceMAC string
	flag.StringVar(&deviceName, "name", "", "Name of the Bluetooth device")
	flag.StringVar(&deviceMAC, "mac_address", "", "MAC address of the Bluetooth device")
	flag.Parse()

	// Check if both the device name and MAC address have been provided
	if deviceName == "" || deviceMAC == "" {
		fmt.Println("Both -name and -mac_address flags must be specified.")
		flag.PrintDefaults() // Print usage information
		os.Exit(1)
	}

	bc, err := NewBluetoothController(deviceName, deviceMAC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing BluetoothController: %s\n", err)
		return
	}

	go bc.ListenForPropertyChanges()

	// Example: Control the device with a hard-coded action
	fmt.Println("Enter 'play', 'pause', 'next', 'previous' to control the device, or 'exit' to quit:")
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