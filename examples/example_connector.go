package main

import (
	"fmt"

	bt "github.com/QuiteLiterallyConnor/BluetoothManager"
)

func connector() {
	bt.EnableDebugging()

	connector, err := bt.NewBluetoothConnector()
	if err != nil {
		fmt.Printf("err connector: %v\n", err)
		return
	}
	connector.StartConnector()

}

func main() {
	connector()
}
