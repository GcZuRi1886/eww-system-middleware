package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var socketPath = "/tmp/eww-system-middleware.sock"

// ----- emitToConsole updates to stdout -----
func emitToConsole(dataType string, data any) {
	dataJSON, _ := json.Marshal(data)
	fmt.Printf("\r%s", string(dataJSON))
}

// ----- main -----
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	args := os.Args
	if len(args) != 2 {
		log.Fatalf("Usage: %s <data_type>", args[0])
	}
	requestedData := args[1]

	switch requestedData {
		case "hyprland":
			go listenHyprlandEventSocket(emitToConsole)
		case "system":
			go sysInfoLoop(emitToConsole)
		case "bluetooth":
			go listenForBluetoothChanges(emitToConsole)
		case "socket":
			_, err := connectToSocket(socketPath)
			if err != nil {
				log.Fatalf("Failed to connect to socket: %v", err)
			}
			go sysInfoLoop(broadcast)
			go listenHyprlandEventSocket(broadcast)
			go listenForBluetoothChanges(broadcast)
		default:
			log.Fatalf("Unknown requested data type: %s", requestedData)
	}
	log.Println("Daemon started. Press Ctrl+C to exit.")
	<-ctx.Done()
	log.Println("Shutting down daemon.")
}

