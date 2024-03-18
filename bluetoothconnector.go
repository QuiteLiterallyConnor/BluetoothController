package bluetoothmanager

import (
	"fmt"
)

type BluetoothConnector struct {
	Scanner *BluetoothScanner
}

func deviceListener(device Device) {
	PrintDebug(fmt.Sprintf("Discovered Device: %s, Paired: %t\n", device.MacAddress, device.Paired))
	if !device.Paired {
		return
	}
	PrintDebug(fmt.Sprintf("Attempting to connect to paired device: %s", device.MacAddress))
	err := device.Connect()
	if err != nil {
		PrintDebug(fmt.Sprintf("Failed to connect to device %s: %s\n", device.MacAddress, err))
	} else {
		PrintDebug(fmt.Sprintf("Successfully connected to device: %s\n", device.MacAddress))
	}
}

func NewBluetoothConnector() (*BluetoothConnector, error) {
	scanner, err := NewBluetoothScanner(deviceListener)
	if err != nil {
		return nil, err
	}

	return &BluetoothConnector{
		Scanner: scanner,
	}, nil
}

func (bc *BluetoothConnector) StartConnector() {
	err := bc.Scanner.GetManagedDevices()
	if err != nil {
		fmt.Printf("Error getting managed devices: %s\n", err)
		return
	}

	bc.Scanner.StartScanner()
}
