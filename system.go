package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/GcZuRi1886/eww-system-middleware/types"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

// ----- periodic system info -----
func sysInfoLoop() {
	for {
		// Time
		now := time.Now().Format("Mon 01 Jan 15:04:05")

		// CPU usage 
		cpuPercent, _ := cpu.Percent(0, true)
		cpuUsage := cpuPercent // per core
		avgPercent, _ := cpu.Percent(0, false)

		// Memory usage
		vm, _ := mem.VirtualMemory()
		totalMem := vm.Total
		usedMem := vm.Used

		// Update battery info
		batteryInfo := getBatteryInfo()

		networkinfo, err := getNetworkInfo()
		if err != nil {
			fmt.Println("Error getting network info:", err)
		}

		state.Mu.Lock()
		state.CurrentState.Time = now
		state.CurrentState.CPUPerCore = cpuUsage
		state.CurrentState.CPUAverage = avgPercent[0]
		state.CurrentState.MemoryUsed = int(usedMem)
		state.CurrentState.MemoryTotal = int(totalMem)
		state.CurrentState.Battery = *batteryInfo
		state.CurrentState.Network = *networkinfo
		state.Mu.Unlock()
		emit()

		time.Sleep(3 * time.Second)
	}
}

// ----- get battery info -----
func getBatteryInfo() *types.BatteryInfo {
	batteryBase := "/sys/class/power_supply/BAT0/"
	
	batteryInfo, err := os.Open(batteryBase + "uevent")
	if err != nil {
		return &types.BatteryInfo{}
	}
	defer batteryInfo.Close()

	batteryData := make(map[string]string)
	scanner := bufio.NewScanner(batteryInfo)
	
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			batteryData[parts[0]] = parts[1]
		}
	}
	

	if len(batteryData) > 0 {
		chargeRate, _ := strconv.ParseFloat(batteryData["POWER_SUPPLY_POWER_NOW"], 64)
		currentCapacity, _ := strconv.ParseFloat(batteryData["POWER_SUPPLY_ENERGY_NOW"], 64)
		fullCapacity, _ := strconv.ParseFloat(batteryData["POWER_SUPPLY_ENERGY_FULL"], 64)
		currentPercentage, _ := strconv.Atoi(batteryData["POWER_SUPPLY_CAPACITY"])
		status := batteryData["POWER_SUPPLY_STATUS"]
		if chargeRate == 0 {
			chargeRate = 1 // to avoid division by zero
		}
		battery := &types.BatteryInfo{
			Percentage:  currentPercentage,
			State:       status,
			TimeToEmpty: currentCapacity / chargeRate * 60, // in minutes
			TimeToFull:  (fullCapacity - currentCapacity) / chargeRate * 60, // in minutes
		}
		return battery
	}
	return &types.BatteryInfo{}
}
