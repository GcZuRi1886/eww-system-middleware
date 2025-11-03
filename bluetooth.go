package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/GcZuRi1886/eww-system-middleware/types"
	"github.com/godbus/dbus/v5"
)

func listenForBluetoothChanges() {
	conn, err := dbus.SystemBus()
	if err != nil {
		log.Fatalf("Failed to connect to system bus: %v", err)
	}
	
	var jsonOutput types.Wrapper
	jsonOutput.Type = "bluetooth"

	jsonOutput.Data = &types.BluetoothInfo{
		Devices: make(map[string]*types.BluetoothDevice),
	}

	// Initial state
	loadInitialState(conn, jsonOutput.Data.(*types.BluetoothInfo))
	emit(jsonOutput)

	// Add a signal match rule for BlueZ Property changes
	rule := "type='signal',sender='org.bluez',interface='org.freedesktop.DBus.Properties',member='PropertiesChanged'"
	call := conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, rule)
	if call.Err != nil {
		log.Fatalf("Failed to add D-Bus match: %v", call.Err)
	}

	// Channel to receive D-Bus signals
	c := make(chan *dbus.Signal, 10)
	conn.Signal(c)

	// Handle system interrupts to exit cleanly
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case sig := <-sigc:
			log.Printf("Exiting on signal: %v", sig)
			return
		case signalMsg := <-c:
			handleSignal(signalMsg, jsonOutput.Data.(*types.BluetoothInfo))
			emit(jsonOutput)
		}
	}
}

// Load initial device + adapter state
func loadInitialState(conn *dbus.Conn, info *types.BluetoothInfo) {
	obj := conn.Object("org.bluez", dbus.ObjectPath("/"))
	var managed map[dbus.ObjectPath]map[string]map[string]dbus.Variant

	err := obj.Call("org.freedesktop.DBus.ObjectManager.GetManagedObjects", 0).Store(&managed)
	if err != nil {
		log.Fatalf("Failed to get managed objects: %v", err)
	}

	for path, ifaces := range managed {
		if adapter, ok := ifaces["org.bluez.Adapter1"]; ok {
			info.Powered = getVariantBool(adapter["Powered"])
		}

		if dev, ok := ifaces["org.bluez.Device1"]; ok {
			devPath := string(path)
			d := &types.BluetoothDevice{
				Name:      getVariantString(dev["Name"]),
				Address:   getVariantString(dev["Address"]),
				Connected: getVariantBool(dev["Connected"]),
				Paired:    getVariantBool(dev["Paired"]),
				Trusted:   getVariantBool(dev["Trusted"]),
				Adapter:   parseAdapterFromPath(devPath),
			}
			info.Devices[devPath] = d
		}
	}
}

// Handle BlueZ property change events
func handleSignal(signalMsg *dbus.Signal, info *types.BluetoothInfo) {
	if len(signalMsg.Body) < 3 {
		return
	}
	iface, ok := signalMsg.Body[0].(string)
	if !ok {
		return
	}

	changedProps, _ := signalMsg.Body[1].(map[string]dbus.Variant)
	path := string(signalMsg.Path)

	switch iface {
	case "org.bluez.Device1":
		dev, ok := info.Devices[path]
		if !ok {
			dev = &types.BluetoothDevice{Adapter: parseAdapterFromPath(path)}
			info.Devices[path] = dev
		}
		for k, v := range changedProps {
			switch k {
			case "Name":
				dev.Name = getVariantString(v)
			case "Address":
				dev.Address = getVariantString(v)
			case "Connected":
				dev.Connected = getVariantBool(v)
			case "Paired":
				dev.Paired = getVariantBool(v)
			case "Trusted":
				dev.Trusted = getVariantBool(v)
			}
		}
	case "org.bluez.Adapter1":
		if powered, ok := changedProps["Powered"]; ok {
			info.Powered = getVariantBool(powered)
		}
	}
}

func printJSON(info *types.BluetoothInfo) {
	jsonBytes, _ := json.Marshal(info)
	fmt.Println(string(jsonBytes))
}

func getVariantString(v dbus.Variant) string {
	if val, ok := v.Value().(string); ok {
		return val
	}
	return ""
}

func getVariantBool(v dbus.Variant) bool {
	if val, ok := v.Value().(bool); ok {
		return val
	}
	return false
}

func parseAdapterFromPath(path string) string {
	// e.g. /org/bluez/hci0/dev_XX_XX_XX_XX_XX_XX
	if len(path) < 10 {
		return ""
	}
	for i := 10; i < len(path); i++ {
		if path[i] == '/' {
			return path[10:i]
		}
	}
	return ""
}

