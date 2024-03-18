package main

import (
	"fmt"

	bt "github.com/QuiteLiterallyConnor/BluetoothController"
)

var ControllerListener = func(event bt.Event) {
	fmt.Printf("Device: %s, Event_Name: %s, Value: %v, Type: %v\n", event.Device, event.Category, event.Value, event.ValueType)
}

func controller() {
	bt.EnableDebugging()

	bc, err := bt.NewBluetoothController(ControllerListener)
	if err != nil {
		fmt.Println("Error initializing Bluetooth Controller:", err)
		return
	}
	bc.StartController()

	// Example of controlling a media player
	// err = bc.ControlMedia("Play", "MAC_ADDRESS")
	// if err != nil {
	// 	fmt.Println("Error controlling media:", err)
	// }

	for {
	}
}

func main() {
	controller()
}
