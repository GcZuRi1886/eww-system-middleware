package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var hyprlandCommandSocket net.Conn
var hyprlandCommandSockerMu sync.Mutex

func openHyprlandSocket(sockName string) (net.Conn, error) {
	sig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	runtimeDir := os.Getenv("XDG_RUNTIME_DIR")
	sock := filepath.Join(runtimeDir, "hypr", sig, sockName)
	addr := net.UnixAddr{Name: sock, Net: "unix"}

	conn, err := net.DialUnix("unix", nil, &addr)
	if err != nil {
		return nil, fmt.Errorf("cannot open Hyprland socket %s: %v", sockName, err)
	}
	return conn, nil
}


// ---- listen and update workspace state over hyprland socket ----
func openHyprlandCommandSocket() error {
	if hyprlandCommandSocket != nil {
		return nil
	}
	conn, err := openHyprlandSocket(".socket.sock")
	if err != nil {
		return err
	}
	hyprlandCommandSocket = conn
	return nil
}

func closeHyprlandCommandSocket() {
	if hyprlandCommandSocket != nil {
		hyprlandCommandSocket.Close()
		hyprlandCommandSocket = nil
	}
}

func sendHyprlandCommand(cmd string) ([]byte, error) {
	hyprlandCommandSockerMu.Lock()
	defer hyprlandCommandSockerMu.Unlock()
	
	err := openHyprlandCommandSocket()
	if err != nil {
		return nil, err
	}

	_, err = hyprlandCommandSocket.Write([]byte(cmd + "\n"))
	if err != nil {
		closeHyprlandCommandSocket()
		return nil, err
	}

	buf := make([]byte, 4096)
	n, err := hyprlandCommandSocket.Read(buf)
	if err != nil {
		closeHyprlandCommandSocket()
		return nil, err
	}
	
	return buf[:n], nil
}

// ----- get current workspace info -----
func getWorkspaceState() {
	println("Fetching workspace state...")
	out, err := sendHyprlandCommand("j/monitors")
	if err != nil {
		return
	}
	println(string(out))
	current := gjsonGet(out, "0.activeWorkspace.id")

	out2, err := sendHyprlandCommand("j/workspaces")
	if err != nil {
		return
	}
	ids := gjsonArrayInt(out2, "id")

	state.mu.Lock()
	state.CurrentState.Workspace.Current = current
	state.CurrentState.Workspace.List = ids
	state.mu.Unlock()
	emit()
}


// ----- listen to hyprland socket -----
func listenHyprlandEventSocket() {
	f, _ := openHyprlandSocket(".socket2.sock")
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "workspace>>") ||
			strings.HasPrefix(line, "createworkspace>>") ||
			strings.HasPrefix(line, "destroyworkspace>>") {
			getWorkspaceState()
		}
	}
}

