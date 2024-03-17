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
    defer wg.Done()

    matchRule := fmt.Sprintf("type='signal',interface='org.freedesktop.DBus.Properties',path='%s'", mediaPlayerPath)
    conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)

    c := make(chan *dbus.Signal, 10)
    conn.Signal(c)

    for range c { // Adjusted to ignore the variable 'v' since it's unused in this snippet.
        // Example processing or print statement could go here.
        // This loop will exit when the channel is closed.
    }
}
