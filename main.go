package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
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
	CPUPerCore   []float64    `json:"cpu_per_core"`
	CPUAverage  float64      `json:"cpu_average"`
	MemoryUsed int   `json:"memory_used"`
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



// ----- periodic system info -----
func sysInfoLoop() {
	for {
		// Time
		now := time.Now().Format("Mon 01 Jan 15:04:05")

		// CPU usage 
		percent, _ := cpu.Percent(0, true)
		cpuUsage := percent // per core
		avgPercent, _ := cpu.Percent(0, false)

		// Memory usage
		vm, _ := mem.VirtualMemory()
		totalMem := vm.Total
		usedMem := vm.Used

		state.mu.Lock()
		state.CurrentState.Time = now
		state.CurrentState.CPUPerCore = cpuUsage
		state.CurrentState.CPUAverage = avgPercent[0]
		state.CurrentState.MemoryUsed = int(usedMem)
		state.CurrentState.MemoryTotal = int(totalMem)
		state.mu.Unlock()
		emit()

		time.Sleep(3 * time.Second)
	}
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

