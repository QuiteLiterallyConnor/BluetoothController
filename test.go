package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/godbus/dbus/v5"
)

var (
	deviceName string
	deviceMAC  string
)

func init() {
	flag.StringVar(&deviceName, "name", "", "Name of the Bluetooth device")
	flag.StringVar(&deviceMAC, "mac_address", "", "MAC address of the Bluetooth device")
	flag.Parse()

	if deviceName == "" || deviceMAC == "" {
		fmt.Println("Both device name and MAC address must be specified using the -name and -mac_address flags.")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func main() {
	deviceMACFormatted := strings.ToUpper(strings.Replace(deviceMAC, ":", "_", -1))
	mediaPlayerPath := fmt.Sprintf("/org/bluez/hci0/dev_%s", deviceMACFormatted)

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to SystemBus: %s\n", err)
		os.Exit(1)
	}

	// Listen for properties changed signal
	conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',interface='org.freedesktop.DBus.Properties',path='"+mediaPlayerPath+"',member='PropertiesChanged'")

	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	fmt.Println("Listening for property changes...")
	for v := range c {
		handlePropertiesChanged(v)
	}
}

func handlePropertiesChanged(signal *dbus.Signal) {
	if len(signal.Body) < 3 {
		return
	}

	interfaceName, ok := signal.Body[0].(string)
	if !ok || interfaceName != "org.bluez.MediaPlayer1" {
		return
	}

	changedProperties, ok := signal.Body[1].(map[string]dbus.Variant)
	if !ok {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	for property, value := range changedProperties {
		switch property {
		case "Volume":
			fmt.Printf("[%s] %s Volume Changed: %v\n", timestamp, deviceName, value.Value())
		case "Status":
			fmt.Printf("[%s] %s Playback Status Changed: %v\n", timestamp, deviceName, value.Value())
		case "Track":
			trackInfo, ok := value.Value().(map[string]dbus.Variant)
			if !ok {
				return
			}
			title, _ := trackInfo["Title"].Value().(string)
			artist, _ := trackInfo["Artist"].Value().(string)
			fmt.Printf("[%s] %s Now Playing: %s by %s\n", timestamp, deviceName, title, artist)
		}
	}
}
