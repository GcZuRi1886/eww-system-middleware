package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/eww-system-middleware/types"
)

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
func openHyprlandCommandSocket() (net.Conn, error) {
	conn, err := openHyprlandSocket(".socket.sock")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func sendHyprlandCommand(cmd string) ([]byte, error) {
	conn, err := openHyprlandCommandSocket()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	
	return buf[:n], nil
}

// ----- get current workspace info -----
func getWorkspaceState() {
	println("Fetching workspace state...")
	out, err := sendHyprlandCommand("j/monitors")
	if err != nil {
		log.Printf("Error getting monitors: %v", err)
		return
	}
	current := readHyprlandWorkspaceCurrent(out)

	out2, err := sendHyprlandCommand("j/workspaces")
	if err != nil {
		log.Printf("Error getting workspaces: %v", err)
		return
	}
	ids := readHyprlandWorkspaceIDs(out2)

	state.mu.Lock()
	state.CurrentState.Workspace.Current = current
	state.CurrentState.Workspace.List = ids
	state.mu.Unlock()
	emit()
}

func readHyprlandWorkspaceIDs(workspacesJSON []byte) []int {
	var wokspaces []types.Workspace
	
	if err := json.Unmarshal(workspacesJSON, &wokspaces); err != nil {
		return nil
	}
	
	var ids []int
	for _, ws := range wokspaces {
		ids = append(ids, ws.ID)
	}
	slices.Sort(ids)
	return ids
}

func readHyprlandWorkspaceCurrent(workspacesJSON []byte) int {
	var monitors []types.Monitor
	if err := json.Unmarshal(workspacesJSON, &monitors); err != nil {
		return 0
	}

	if len(monitors) == 0 {
		return 0
	}

	return monitors[0].ActiveWorkspace.ID
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

