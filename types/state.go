// Package types defines types for this project
package types

import "sync"

type State struct {
	CurrentState CurrentStateData
	Mu           sync.Mutex
}

type CurrentStateData struct {
	Workspace 	WorkspaceInfo `json:"workspace"`
	Time  			string 				`json:"time"`
	CPUPerCore  []float64    	`json:"cpu_per_core"`
	CPUAverage  float64      	`json:"cpu_average"`
	MemoryUsed 	int   				`json:"memory_used"`
	MemoryTotal int   				`json:"memory_total"`
	Battery 		BatteryInfo		`json:"battery"`
	Network 		NetworkInfo		`json:"network"`
}

type WorkspaceInfo struct {
		Current int   `json:"current"`
		List    []int `json:"list"`
}

type BatteryInfo struct {
		Percentage 	int 		`json:"percentage"`
		State 			string 	`json:"state"`
		TimeToEmpty float64	`json:"time_to_empty"`
		TimeToFull  float64	`json:"time_to_full"`
}

type NetworkInfo struct {
  Interface      string  `json:"interface"`
  IPAddress      string  `json:"ip_address"`
  IsWifi         bool    `json:"is_wifi"`
  SSID           string  `json:"ssid,omitempty"`
  IsConnected    bool    `json:"is_connected"`
  SignalStrength int     `json:"signal_strength"` // in percentage
  DownSpeed      float64 `json:"down_speed_kbps"`
  UpSpeed        float64 `json:"up_speed_kbps"`
  BytesSent      uint64  `json:"bytes_sent"`
  BytesRecv      uint64  `json:"bytes_recv"`
}

