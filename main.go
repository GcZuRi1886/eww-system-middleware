package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)
type State struct {
	CurrentState CurrentStateData `json:"current_state"`
	mu 	sync.Mutex
}

type CurrentStateData struct {
	Workspace struct {
		Current int   `json:"current"`
		List    []int `json:"list"`
	} `json:"workspace"`
	Time  string `json:"time"`
	CPU   []float64    `json:"cpu"`
	MemoryUsed int   `json:"memory"`
	MemoryTotal int   `json:"memory_total"`
}

var state State

// ----- emit updates to stdout -----
func emit() {
	state.mu.Lock()
	defer state.mu.Unlock()
	data := state.CurrentState
	dataJSON, _ := json.Marshal(data)
	fmt.Println(string(dataJSON))
}



// minimal gjson-like helpers (no external deps)
func gjsonGet(jsonData []byte, key string) int {
	// very naive extractor for small arrays like monitors JSON
	s := string(jsonData)
	idx := strings.Index(s, key)
	if idx == -1 {
		return 0
	}
	var val int
	fmt.Sscanf(s[idx+len(key)+2:], "%d", &val)
	return val
}
func gjsonArrayInt(jsonData []byte, key string) []int {
	s := string(jsonData)
	lines := strings.Split(s, key)
	var ids []int
	for i := 1; i < len(lines); i++ {
		var val int
		fmt.Sscanf(lines[i][2:], "%d", &val)
		ids = append(ids, val)
	}
	sort.Ints(ids)
	return ids
}


// ----- periodic system info -----
func sysInfoLoop() {
	for {
		// Time
		now := time.Now().Format("15:04")

		// CPU usage 
		percent, _ := cpu.Percent(0, true)
		cpuUsage := percent 

		// Memory usage
		vm, _ := mem.VirtualMemory()
		totalMem := vm.Total
		usedMem := vm.Used

		state.mu.Lock()
		state.CurrentState.Time = now
		state.CurrentState.CPU = cpuUsage
		state.CurrentState.MemoryUsed = int(usedMem)
		state.CurrentState.MemoryTotal = int(totalMem)
		state.mu.Unlock()
		emit()

		time.Sleep(3 * time.Second)
	}
}

// ----- main -----
func main() {
	openHyprlandCommandSocket()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	defer closeHyprlandCommandSocket()

	go listenHyprlandEventSocket()
	go sysInfoLoop()

	log.Println("Daemon started. Press Ctrl+C to exit.")
	<-ctx.Done()
	log.Println("Shutting down daemon.")
}

