package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    "sync"

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
    mediaPlayerPath := fmt.Sprintf("/org/bluez/hci0/dev_%s/player0", deviceMACFormatted)

    conn, err := dbus.SystemBus()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to connect to SystemBus: %s\n", err)
        os.Exit(1)
    }

    var wg sync.WaitGroup
    wg.Add(1)

    go listenForPropertiesChanged(conn, mediaPlayerPath, &wg)

    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter command (play, pause, next, previous, quit): ")
        command, _ := reader.ReadString('\n')
        command = strings.TrimSpace(command)

        if command == "quit" {
            fmt.Println("Exiting...")
            break
        }

        controlMedia(conn, mediaPlayerPath, strings.Title(command))
    }

    // Inform the properties changed listener to stop and wait for it
    wg.Done() // Mark the listener as done
    wg.Wait() // Wait for the listener goroutine to finish
}

func controlMedia(conn *dbus.Conn, mediaPlayerPath, method string) {
    mediaPlayer := conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
    call := mediaPlayer.Call("org.bluez.MediaPlayer1."+method, 0)
    if call.Err != nil {
        fmt.Fprintf(os.Stderr, "Failed to %s: %s\n", strings.ToLower(method), call.Err)
        return
    }
    fmt.Printf("%s action executed for %s\n", method, deviceName)
}

func listenForPropertiesChanged(conn *dbus.Conn, mediaPlayerPath string, wg *sync.WaitGroup) {
	defer wg.Done() // Mark the listener as done when the function returns

	mediaPlayer := conn.Object("org.bluez", dbus.ObjectPath(mediaPlayerPath))
	mediaPlayerIface := "org.freedesktop.DBus.Properties"

	// Add a match rule for PropertiesChanged signal
	matchRule := fmt.Sprintf("type='signal',interface='%s',path='%s'", mediaPlayerIface, mediaPlayerPath)
	err := conn.AddMatchSignal(matchRule)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to add match rule for PropertiesChanged signal: %s\n", err)
		return
	}
	defer conn.RemoveMatchSignal(matchRule)

	// Listen for PropertiesChanged signals
	ch := make(chan *dbus.Signal, 10)
	conn.Signal(ch)
	for signal := range ch {
		if signal.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" {
			onPropertiesChanged(signal)
		}
	}
}

func onPropertiesChanged(signal *dbus.Signal) {
    if len(signal.Body) >= 3 {
        interfaceName := signal.Body[0].(string)
        changedProperties := signal.Body[1].(map[string]dbus.Variant)
        // invalidatedProperties := signal.Body[2] // Depending on your needs

        fmt.Println("PropertiesChanged on interface:", interfaceName)
        for propName, propValue := range changedProperties {
            fmt.Printf("Property %s changed to %v\n", propName, propValue)
        }
    }
}