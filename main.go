package main

import (
	"Car-Bluetooth/bluetooth_manager" // Replace with the actual module path
	"fmt"
)

func main() {
	bc, err := bluetootheventlistener.NewBluetoothController("Pixel_6", "0C:C4:13:12:67:62")
	if err != nil {
		fmt.Println("Error creating BluetoothController:", err)
		return
	}

	go bc.ListenForPropertyChanges()
	bc.ControlMedia("play")
}
