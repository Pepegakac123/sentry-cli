package monitor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
)

type MemoryUsage struct {
	TotalMemory     int64
	AvailableMemory int64
	MemoryUsage     float64
}
type MemoryMonitor struct {
	Stats *MemoryUsage
}

func NewMemoryUsageMonitor() (*MemoryMonitor, error) {
	stats, err := getRamUsage()
	if err != nil {
		return nil, err
	}
	return &MemoryMonitor{Stats: stats}, nil
}

func (m *MemoryMonitor) Update() (*MemoryUsage, error) {
	stats, err := getRamUsage()
	if err != nil {
		return nil, err
	}
	m.Stats = stats
	return stats, nil
}

func getRamUsage() (*MemoryUsage, error) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var totalMem, availableMem int64
	prefixMemTotal := []byte("MemTotal:")
	prefixMemAvailable := []byte("MemAvailable:")

	for scanner.Scan() {
		line := scanner.Bytes()

		if bytes.HasPrefix(line, prefixMemTotal) {
			totalMem, err = parseValue(line)
			totalMem = totalMem * 1024
			if err != nil {
				return nil, err
			}
		}
		if bytes.HasPrefix(line, prefixMemAvailable) {
			availableMem, err = parseValue(line)
			availableMem = availableMem * 1024
			if err != nil {
				return nil, err
			}
		}

		if totalMem > 0 && availableMem > 0 {
			break
		}
	}
	if totalMem == 0 || availableMem == 0 {
		return nil, fmt.Errorf("failed to get memory usage")
	}
	memUsagePrcent := float64(totalMem-availableMem) / float64(totalMem) * 100
	stats := &MemoryUsage{
		TotalMemory:     totalMem,
		AvailableMemory: availableMem,
		MemoryUsage:     memUsagePrcent,
	}
	return stats, nil

}
func parseValue(line []byte) (int64, error) {
	fields := bytes.Fields(line)
	if len(fields) < 2 {
		return 0, fmt.Errorf("invalid line format")
	}
	return strconv.ParseInt(string(fields[1]), 10, 64)
}
