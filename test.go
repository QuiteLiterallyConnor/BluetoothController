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

    // Setup the match rule for listening to PropertiesChanged signals
    // You might need to adjust the rule based on your specific needs,
    // for example, by specifying a particular object path or interface.
    matchRule := "type='signal',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged'"
    conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchRule)

    // Create a channel to receive signals
    c := make(chan *dbus.Signal, 10)
    conn.Signal(c)

    fmt.Println("Listening for PropertiesChanged signals...")
    for v := range c {
		fmt.Printf("Received signal: %+v\n", v)
        onPropertiesChanged(v)
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
