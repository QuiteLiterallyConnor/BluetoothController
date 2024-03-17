package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
)

var (
	bluezDest   = "org.bluez"
	deviceName  string
	deviceMAC   string
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
	var action string

	fmt.Print("Enter the action to perform (play or pause): ")
	fmt.Scanln(&action)

	deviceMACFormatted := strings.ToUpper(strings.Replace(deviceMAC, ":", "_", -1))
	mediaControlPath := fmt.Sprintf("/org/bluez/hci0/dev_%s/player0", deviceMACFormatted)

	conn, err := dbus.SystemBus()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to SystemBus: %s\n", err)
		os.Exit(1)
	}

	switch action {
	case "play":
		playMedia(conn, mediaControlPath)
	case "pause":
		pauseMedia(conn, mediaControlPath)
	default:
		fmt.Println("Invalid action. Use 'play' or 'pause'")
	}
}

func playMedia(conn *dbus.Conn, mediaControlPath string) {
	call := conn.Object(bluezDest, dbus.ObjectPath(mediaControlPath)).Call("org.bluez.MediaControl1.Play", 0)
	if call.Err != nil {
		fmt.Fprintf(os.Stderr, "Failed to play media: %s\n", call.Err)
		return
	}
	fmt.Println("Playback started for", deviceName)
}

func pauseMedia(conn *dbus.Conn, mediaControlPath string) {
	call := conn.Object(bluezDest, dbus.ObjectPath(mediaControlPath)).Call("org.bluez.MediaControl1.Pause", 0)
	if call.Err != nil {
		fmt.Fprintf(os.Stderr, "Failed to pause media: %s\n", call.Err)
		return
	}
	fmt.Println("Playback paused for", deviceName)
}
