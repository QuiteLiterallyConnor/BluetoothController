package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "strings"
    // "time"

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
        fmt.Println("Both device name and MAC address must be specified.")
        flag.PrintDefaults()
        os.Exit(1)
    }
}

func listenForPropertiesChanged(conn *dbus.Conn, devicePath string) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    matchSignal := fmt.Sprintf("type='signal',interface='org.freedesktop.DBus.Properties',path='%s'", devicePath)
    conn.BusObject().CallWithContext(ctx, "org.freedesktop.DBus.AddMatch", 0, matchSignal)

    c := make(chan *dbus.Signal, 10)
    conn.Signal(c)

    for {
        select {
        case v := <-c:
            handlePropertiesChanged(v)
        case <-ctx.Done():
            return
        }
    }
}

func handlePropertiesChanged(signal *dbus.Signal) {
    // Implementation similar to the Python version's on_properties_changed
    // Extract and process properties from the signal.Body
    fmt.Println("PropertiesChanged signal received:", signal)
    // Note: You'll need to adjust the extraction logic based on the signal's actual structure.
}

func main() {
    deviceMACFormatted := strings.ToUpper(strings.Replace(deviceMAC, ":", "_", -1))
    devicePath := fmt.Sprintf("/org/bluez/hci0/dev_%s", deviceMACFormatted)

    conn, err := dbus.SystemBus()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to connect to SystemBus: %s\n", err)
        os.Exit(1)
    }

    go listenForPropertiesChanged(conn, devicePath)

    // Keep the main goroutine running to listen for signals
    fmt.Println("Listening for property changes. Press CTRL+C to exit.")
    <-make(chan struct{})
}
