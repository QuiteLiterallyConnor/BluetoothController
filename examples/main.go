package main

import (
	"fmt"
	"reflect"

	bt "github.com/QuiteLiterallyConnor/BluetoothController"
)

func main() {
	listener := func(device, event string, value interface{}, valueType reflect.Type) {
		fmt.Printf("Device: %s, Event: %s, Value: %v, Type: %v\n", device, event, value, valueType)
	}

	bc, err := bt.NewBluetoothController(listener)
	if err != nil {
		fmt.Println("Error initializing Bluetooth Controller:", err)
		return
	}
	// bc.EnableDebugging()
	bc.Start()

	// Example of controlling a media player
	// err = bc.ControlMedia("Play", "MAC_ADDRESS")
	// if err != nil {
	// 	fmt.Println("Error controlling media:", err)
	// }

	for {
	}
}
