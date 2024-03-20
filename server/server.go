package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"

	bt "github.com/QuiteLiterallyConnor/BluetoothManager"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

var clients = make(map[*websocket.Conn]bool)
var mutex = &sync.Mutex{}

var ControllerListener = func(event bt.Event) {
	fmt.Printf("CONTROLLER LISTENER | Device: %s, Event_Name: %s, Value: %v, Type: %v\n", event.Device, event.Category, event.Value, event.ValueType)

	var message string
	switch event.Category {
	case "Status":
		message = fmt.Sprintf(`{"type": "status", "value": "%s"}`, event.Value)
	case "Volume":
		message = fmt.Sprintf(`{"type": "volume", "value": %v}`, event.Value)
	case "Track":
		message = fmt.Sprintf(`{"type": "track", "value": %v}`, event.Value)
	}

	if message != "" {
		broadcastToClients(event.Json())
	}
}

var ScannerListener = func(device bt.Device) {
	fmt.Printf("Scanner LISTENER | Device: %s, Data: %+v, DataType: %v\n\n", device.MacAddress, device, reflect.TypeOf(device))
	broadcastToClients(fmt.Sprintf("Scanner Device - Mac: %s, Device: %+v\n", device.MacAddress, device))
}

var (
	controller *bt.BluetoothController
	scanner    *bt.BluetoothScanner
)

func main() {
	r := gin.Default()
	r.Static("public", "./public")
	r.LoadHTMLGlob("./public/*.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/ws", wshandler)

	initBluetooth()
	controller.StartController()
	scanner.StartScanner()

	r.Run(":8080")

}

func initBluetooth() {
	var err error
	controller, err = bt.NewBluetoothController(ControllerListener)

	if err != nil {
		fmt.Printf("Failed to initialize Bluetooth Controller: %v", err)
		os.Exit(1)
	}

	scanner, err = bt.NewBluetoothScanner(ScannerListener)

	if err != nil {
		fmt.Printf("Failed to initialize Bluetooth Scanner: %v", err)
		os.Exit(1)
	}
}

func wshandler(c *gin.Context) {
	w := c.Writer
	r := c.Request
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	defer conn.Close()

	// Add connection to the list of clients
	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	for {
		// You can also listen for messages from the client if needed
	}
}

// broadcastToClients sends a message to all connected WebSocket clients
func broadcastToClients(message string) {
	mutex.Lock()
	defer mutex.Unlock()
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Printf("WebSocket send error: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
