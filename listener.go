package main

import (
    "fmt"
    "os"

    "github.com/godbus/dbus/v5"
)

func main() {
    conn, err := dbus.SystemBus()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to connect to SystemBus: %s\n", err)
        os.Exit(1)
    }

    matchRule := "type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged'"
    conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)

    c := make(chan *dbus.Signal, 10)
    conn.Signal(c)

    fmt.Println("Listening for PropertiesChanged signals...")
    for v := range c {
        onPropertiesChanged(v)
    }
}

func onPropertiesChanged(signal *dbus.Signal) {
    if len(signal.Body) >= 3 {
        interfaceName := signal.Body[0].(string)
        changedProperties := signal.Body[1].(map[string]dbus.Variant)

        fmt.Println("PropertiesChanged on interface:", interfaceName)
        for propName, propValue := range changedProperties {
            fmt.Printf("Property %s changed to %v\n", propName, propValue)
        }
    }
}
