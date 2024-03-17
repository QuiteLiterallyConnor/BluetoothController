package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
)
type BluetoothController struct {
	DeviceName 	string
	DeviceMAC  	string
	Conn       	*dbus.Conn
	Listeners 	map[string]Event
	Broadcaster map[string]Event
}

type Event struct {
	Name 		string
	Function 	func()
}

func NewBluetoothController(deviceName, deviceMAC string, listeners, receivers []Event) (*BluetoothController, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SystemBus: %w", err)
	}

	listenerMap := make(map[string]Event)
	for _, event := range listeners {
		listenerMap[event.Name] = event
	}

	broadcasterMap := make(map[string]Event)
	for _, event := range receivers {
		broadcasterMap[event.Name] = event
	}

	return &BluetoothController{
		DeviceName: 	deviceName,
		DeviceMAC:  	deviceMAC,
		Conn:       	conn,
		Listeners:  	listenerMap,
		Broadcaster:	broadcasterMap,
	}, nil

}

func (bs *BluetoothController) Start() {
	go bc.listenForControlMedia()
	bs.ListenForPropertyChanges()
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

			// bc.Listeners[propName].Function()
		}
	}
}

func (bc *BluetoothController) listenForControlMedia() {
	for {
		var action string
		fmt.Scanln(&action)
		action = strings.TrimSpace(action)

		if action == "exit" {
			break
		}

		if err := ControlMedia(action); err != nil {
			fmt.Fprintf(os.Stderr, "Error controlling media: %s\n", err)
		}
	}
}

func (bc *BluetoothController) ControlMedia(action string) error { // Adjusted to return an error
	deviceMACFormatted := strings.ToUpper(strings.Replace(bc.DeviceMAC, ":", "_", -1))
	mediaPlayerPath := fmt.Sprintf("/org/bluez/hci0/dev_%s/player0", deviceMACFormatted)

	mediaPlayer := bc.Conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	call := mediaPlayer.Call("org.bluez.MediaPlayer1."+strings.Title(action), 0)
	if call.Err != nil {
		return fmt.Errorf("failed to %s: %w", strings.ToLower(action), call.Err) // Error wrapping for better handling
	}

	fmt.Printf("%s action executed for %s\n", action, bc.DeviceName)



	return nil
}

func functionThatPlaysOnReceivePause() {
	fmt.Println("Receive Pausing...")
}

func functionThatPlaysOnReceiveStop() {
	fmt.Println("Receive Stopping...")
}

func functionThatPlaysOnReceivePlay() {
	fmt.Println("Receive Playing...")
}

func functionThatPlaysOnReceiveNext() {
	fmt.Println("Receive Next track...")
}

func functionThatPlaysOnReceivePrevious() {
	fmt.Println("Receive Previous track...")
}

func functionThatPlaysOnReceiveVolumeChange() {
	fmt.Println("Receive Volume changed...")
}

func functionThatPlaysOnReceiveTrack() {	
	fmt.Println("Receive Track changed...")
}



///


func functionThatPlaysOnSendPause() {
	fmt.Println("Sending Pausing...")
}

func functionThatPlaysOnSendStop() {
	fmt.Println("Sending Stopping...")
}

func functionThatPlaysOnSendPlay() {
	fmt.Println("Sending Playing...")
}

func functionThatPlaysOnSendNext() {
	fmt.Println("Sending Playing next track...")
}

func functionThatPlaysOnSendPrevious() {
	fmt.Println("Sending Playing previous track...")
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

    listeners := []Event{
        {"Pause", functionThatPlaysOnReceivePause},
        {"Stop", functionThatPlaysOnReceiveStop},
        {"Play", functionThatPlaysOnReceivePlay},
        {"Next", functionThatPlaysOnReceiveNext},
        {"Previous", functionThatPlaysOnReceivePrevious},
        {"VolumeChange", functionThatPlaysOnReceiveVolumeChange},
        {"Track", functionThatPlaysOnReceiveTrack},
    }

    broadcasters := []Event{
        {"Pause", functionThatPlaysOnSendPause},
        {"Stop", functionThatPlaysOnSendStop},
        {"Play", functionThatPlaysOnSendPlay},
        {"Next", functionThatPlaysOnSendNext},
        {"Previous", functionThatPlaysOnSendPrevious},
    }

    bc, err := NewBluetoothController(deviceName, deviceMAC, listeners, broadcasters)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error initializing BluetoothController: %s\n", err)
        os.Exit(1)
    }

    bc.Start()
}