package types

type Monitor struct {
	ID              int      `json:"id"`
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	Make            string   `json:"make"`
	Model           string   `json:"model"`
	Serial          string   `json:"serial"`
	Width           int      `json:"width"`
	Height          int      `json:"height"`
	RefreshRate     float64  `json:"refreshRate"`
	X               int      `json:"x"`
	Y               int      `json:"y"`
	ActiveWorkspace struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"activeWorkspace"`
	SpecialWorkspace struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"specialWorkspace"`
	Reserved          [4]int    `json:"reserved"`
	Scale             float64   `json:"scale"`
	Transform         int       `json:"transform"`
	Focused           bool      `json:"focused"`
	DpmsStatus        bool      `json:"dpmsStatus"`
	Vrr               bool      `json:"vrr"`
	Solitary          string    `json:"solitary"`
	ActivelyTearing   bool      `json:"activelyTearing"`
	DirectScanoutTo   string    `json:"directScanoutTo"`
	Disabled          bool      `json:"disabled"`
	CurrentFormat     string    `json:"currentFormat"`
	MirrorOf          string    `json:"mirrorOf"`
	AvailableModes    []string  `json:"availableModes"`
}


type Workspace struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Monitor          string `json:"monitor"`
	MonitorID        int    `json:"monitorID"`
	Windows          int    `json:"windows"`
	HasFullscreen    bool   `json:"hasfullscreen"`
	LastWindow       string `json:"lastwindow"`
	LastWindowTitle  string `json:"lastwindowtitle"`
	IsPersistent     bool   `json:"ispersistent"`
}

type WorkspaceInfo struct {
		Current int   `json:"current_workspace"`
		List    []int `json:"workspace_list"`
}

