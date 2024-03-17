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

    // Replace the previous action prompt with a loop that allows continuous command input
    reader := bufio.NewReader(os.Stdin)
    for {
        fmt.Print("Enter command (play, pause, next, previous, quit): ")
        command, _ := reader.ReadString('\n')
        command = strings.TrimSpace(command)

        if command == "quit" {
            fmt.Println("Exiting...")
            break
        }

        switch command {
        case "play", "pause", "next", "previous":
            controlMedia(conn, mediaPlayerPath, strings.Title(command))
        default:
            fmt.Println("Invalid command. Use 'play', 'pause', 'next', 'previous', or 'quit'")
        }
    }

    wg.Wait() // Wait for the listener goroutine to finish
}

// Updated controlMedia function as previously defined

// Updated listenForPropertiesChanged function with a WaitGroup parameter
func listenForPropertiesChanged(conn *dbus.Conn, mediaPlayerPath string, wg *sync.WaitGroup) {
    defer wg.Done()

    matchRule := fmt.Sprintf("type='signal',interface='org.freedesktop.DBus.Properties',path='%s'", mediaPlayerPath)
    conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)

    c := make(chan *dbus.Signal, 10)
    conn.Signal(c)

    for v := range c {
        // Process the signal here
        // Break loop or return if a "quit" condition is met, if necessary
    }
}
