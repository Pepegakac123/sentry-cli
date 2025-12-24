package monitor

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

type MemoryUsage struct {
	TotalMemory     int64
	AvailableMemory int64
	MemoryUsage     float64
}
type CPUStats struct {
	Idle  uint64
	Total uint64
}

func GetRamUsage() (*MemoryUsage, error) {
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
			if err != nil {
				return nil, err
			}
		}
		if bytes.HasPrefix(line, prefixMemAvailable) {
			availableMem, err = parseValue(line)
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

	return &MemoryUsage{
		TotalMemory:     totalMem,
		AvailableMemory: availableMem,
		MemoryUsage:     memUsagePrcent,
	}, nil

}
func parseValue(line []byte) (int64, error) {
	fields := bytes.Fields(line)
	if len(fields) < 2 {
		return 0, fmt.Errorf("invalid line format")
	}
	return strconv.ParseInt(string(fields[1]), 10, 64)
}

func GetCpuUsage() (float64, error) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, err
	}
	defer f.Close()

	s1, err := getCpuSnapshot(f)
	if err != nil {
		return 0, err
	}
	time.Sleep(1 * time.Second)
	_, err = f.Seek(0, 0)
	if err != nil {
		return 0, err
	}
	s2, err := getCpuSnapshot(f)
	if err != nil {
		return 0, err
	}

	return calculateCpuUsage(s1, s2), nil
}
func getCpuSnapshot(f *os.File) (*CPUStats, error) {
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
func calculateCpuUsage(s1, s2 *CPUStats) float64 {
	if s1.Total == 0 || s2.Total == 0 {
		return 0
	}
	totalDelta := s2.Total - s1.Total
	idleDelta := s2.Idle - s1.Idle
	return float64(totalDelta-idleDelta) / float64(totalDelta) * 100
}
