package main

import (
	"fmt"
	"reflect"

	bt "github.com/QuiteLiterallyConnor/BluetoothController"
)

var Listener = func(device, event string, value interface{}, valueType reflect.Type) {
	fmt.Printf("Device: %s, Event: %s, Value: %v, Type: %v\n", device, event, value, valueType)
}

func controller() {

	bc, err := bt.NewBluetoothController(Listener)
	if err != nil {
		fmt.Println("Error initializing Bluetooth Controller:", err)
		return
	}
	// bc.EnableDebugging()
	bc.StartController()

	// Example of controlling a media player
	err = bc.ControlMedia("Play", "0C:C4:13:12:67:62")
	if err != nil {
		fmt.Println("Error controlling media:", err)
	}

	for {
	}
}

func scanner() {
	bt.EnableDebugging()

	bs, err := bt.NewBluetoothScanner()
	if err != nil {
		fmt.Println("Error initializing Bluetooth Scanner:", err)
		return
	}

	if err = bs.GetManagedDevices(); err != nil {
		fmt.Printf("%w\n", err)
	}

	for name, device := range bs.Devices {
		fmt.Printf("Name: %s, Device: %+v\n", name, device)
	}

	bs.ListenForDevices()

	for {
	}
}

func main() {
	scanner()
	// controller()
}
