package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// Groups of clients subscribed to each info type
var subscribers = struct {
	sync.RWMutex
	m map[string]map[net.Conn]bool // type → set of connections
}{m: make(map[string]map[net.Conn]bool)}

// Broadcast message of a specific type to all subscribers of that type
func broadcast(infoType string, data any) {
	msg, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
		return
	}
	msg = append(msg, '\n')

	subscribers.RLock()
	defer subscribers.RUnlock()

	conns, ok := subscribers.m[strings.ToUpper(infoType)]
	if !ok {
		return
	}

	for c := range conns {
		_, err := c.Write([]byte(msg))
		if err != nil {
			fmt.Printf("Write error, removing client: %v\n", err)
			removeClientFromAllTypes(c)
			c.Close()
		}
	}
}

// Subscribe a client to a data type
func subscribe(c net.Conn, infoType string) {
	subscribers.Lock()
	defer subscribers.Unlock()

	if subscribers.m[infoType] == nil {
		subscribers.m[infoType] = make(map[net.Conn]bool)
	}

	subscribers.m[infoType][c] = true
}

// Remove client from all subscription lists
func removeClientFromAllTypes(c net.Conn) {
	subscribers.Lock()
	defer subscribers.Unlock()

	for _, conns := range subscribers.m {
		delete(conns, c)
	}
}

// Handle an individual client session
func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected")
			removeClientFromAllTypes(conn)
			return
		}

		cmd := strings.TrimSpace(line)

		if cmd == "" {
			continue
		}

		// Example: "SUB CPU" or "SUB MEMORY"
		parts := strings.SplitN(cmd, " ", 2)

		if len(parts) == 2 && strings.ToUpper(parts[0]) == "SUB" {
			infoType := strings.ToUpper(parts[1])
			subscribe(conn, infoType)
			conn.Write([]byte("OK subscribed to " + infoType + "\n"))
			continue
		}

		conn.Write([]byte("ERROR unknown command\n"))
	}
}

// connectToSocket creates and listens on the Unix socket
func connectToSocket(socketPath string) (net.Listener, error) {

	// Remove stale socket file
	if _, err := os.Stat(socketPath); err == nil {
		if err := os.Remove(socketPath); err != nil {
			return nil, fmt.Errorf("failed removing stale socket: %w", err)
		}
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on unix socket: %w", err)
	}

	fmt.Printf("Server listening on %s\n", socketPath)

	// Accept clients
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Printf("Accept error: %v\n", err)
				continue
			}

			fmt.Println("Client connected")
			go handleClient(conn)
		}
	}()

	return listener, nil
}

// ---- Your system info loops ----
func startSystemInfoLoops() {
	// Example loop: send CPU info
	go func() {
		for {
			cpuData := "CPU: 14%\n"
			broadcast("CPU", cpuData)
		}
	}()

	// Example loop: send memory info
	go func() {
		for {
			memData := "MEMORY: 2.3GB used\n"
			broadcast("MEMORY", memData)
		}
	}()
}

