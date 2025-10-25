package types


import "sync"

type State struct {
	CurrentState CurrentStateData
	Mu           sync.Mutex
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
	Battery struct {
		Percentage int `json:"percentage"`
		State	string    `json:"is_charging"`
		TimeToEmpty  float64     `json:"time_to_empty"`
		TimeToFull   float64     `json:"time_to_full"`
	} `json:"battery"`
}
