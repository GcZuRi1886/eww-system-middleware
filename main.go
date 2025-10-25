package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GcZuRi1886/eww-system-middleware/types"
)

var state types.State

// ----- emit updates to stdout -----
func emit() {
	state.Mu.Lock()
	defer state.Mu.Unlock()
	data := state.CurrentState
	dataJSON, _ := json.Marshal(data)
	fmt.Println(string(dataJSON))
}

// ----- main -----
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	getWorkspaceState()

	go listenHyprlandEventSocket()
	go sysInfoLoop()

	log.Println("Daemon started. Press Ctrl+C to exit.")
	<-ctx.Done()
	log.Println("Shutting down daemon.")
}

