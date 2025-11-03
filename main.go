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

// ----- emit updates to stdout -----
func emit(data any) {
	dataJSON, _ := json.Marshal(data)
	fmt.Println(string(dataJSON))
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
			go listenHyprlandEventSocket()
		case "system":
			go sysInfoLoop()
		case "bluetooth":
			go listenForBluetoothChanges()
		default:
			log.Fatalf("Unknown requested data type: %s", requestedData)
	}
	log.Println("Daemon started. Press Ctrl+C to exit.")
	<-ctx.Done()
	log.Println("Shutting down daemon.")
}

