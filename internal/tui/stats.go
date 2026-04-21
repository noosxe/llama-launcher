package tui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type MachineStats struct {
	CPU  string
	RAM  string
	GPU  string
	VRAM string
}

type statsMsg MachineStats
type tickStatsMsg time.Time

func tickStatsCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return tickStatsMsg(t)
	})
}

func fetchStatsCmd() tea.Msg {
	stats := MachineStats{
		CPU:  "N/A",
		RAM:  "N/A",
		GPU:  "N/A",
		VRAM: "N/A",
	}

	// Fetch GPU via nvidia-smi
	out, err := exec.Command("nvidia-smi", "--query-gpu=utilization.gpu,memory.used,memory.total", "--format=csv,noheader").Output()
	if err == nil {
		parts := strings.Split(strings.TrimSpace(string(out)), ",")
		if len(parts) >= 3 {
			stats.GPU = strings.TrimSpace(parts[0])
			stats.VRAM = fmt.Sprintf("%s / %s", strings.TrimSpace(parts[1]), strings.TrimSpace(parts[2]))
		}
	}

	// Fetch RAM via free
	out, err = exec.Command("free", "-h").Output()
	if err == nil {
		lines := strings.Split(string(out), "\n")
		if len(lines) > 1 {
			fields := strings.Fields(lines[1])
			if len(fields) >= 3 {
				stats.RAM = fmt.Sprintf("%s / %s", fields[2], fields[1])
			}
		}
	}

	// Fetch CPU via top
	out, err = exec.Command("sh", "-c", "top -bn1 | grep 'Cpu(s)' | awk '{print $2 + $4}' | awk '{printf \"%.1f%%\", $1}'").Output()
	if err == nil {
		cpuVal := strings.TrimSpace(string(out))
		if cpuVal != "" {
			stats.CPU = cpuVal
		}
	}

	return statsMsg(stats)
}
