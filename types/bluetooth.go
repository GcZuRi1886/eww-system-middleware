package types

type BluetoothDevice struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	Connected bool   `json:"connected"`
	Paired    bool   `json:"paired"`
	Trusted   bool   `json:"trusted"`
	Adapter   string `json:"adapter"`
}

type BluetoothInfo struct {
	Powered bool              	`json:"powered"`
	Devices map[string]*BluetoothDevice 	`json:"devices"`
}
