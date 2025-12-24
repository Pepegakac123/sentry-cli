package monitor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

// CPUStats trzyma surowe liczniki z /proc/stat
type CPUStats struct {
	Idle  uint64
	Total uint64
}

// CPUMonitor to nasz obiekt stanowy. Pamięta tylko ostatni odczyt.
type CPUMonitor struct {
	lastStats *CPUStats
}

// NewCPUMonitor inicjalizuje monitor i robi pierwszy odczyt
func NewCPUMonitor() (*CPUMonitor, error) {
	stats, err := getCpuSnapshot()
	if err != nil {
		return nil, err
	}
	return &CPUMonitor{
		lastStats: stats,
	}, nil
}

func (m *CPUMonitor) Update() (float64, error) {
	currentStats, err := getCpuSnapshot()
	if err != nil {
		return 0, err
	}
	totalDelta := currentStats.Total - m.lastStats.Total
	idleDelta := currentStats.Idle - m.lastStats.Idle

	m.lastStats = currentStats

	if totalDelta == 0 {
		return 0, nil
	}

	usedDelta := totalDelta - idleDelta
	usage := float64(usedDelta) / float64(totalDelta) * 100
	return usage, nil
}
func getCpuSnapshot() (*CPUStats, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var prefixCpu = []byte("cpu ")
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, prefixCpu) {
			fields := bytes.Fields(line)
			if len(fields) < 5 {
				return nil, fmt.Errorf("invalid line format")
			}
			var total uint64
			for i := 1; i < len(fields); i++ {
				val, err := strconv.ParseUint(string(fields[i]), 10, 64)
				if err != nil {
					return nil, err
				}
				total += val
			}
			idle, _ := strconv.ParseUint(string(fields[4]), 10, 64)

			return &CPUStats{
				Idle:  idle,
				Total: total,
			}, nil
		}
	}
	return nil, scanner.Err()
}
