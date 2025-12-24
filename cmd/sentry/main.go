package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Pepegakac123/sentry-cli/pkg/monitor"
)

var version = "dev"

func main() {
	var color string
	interval := flag.Duration("interval", 2*time.Second, "Check interval (e.g. 1s, 500ms)")
	flag.Parse()
	fmt.Printf("Sentry CLI version: %s\nInterval: %s\n", version, interval)
	cpuMon, err := monitor.NewCPUMonitor()
	if err != nil {
		log.Fatal(err)
	}
	ramMon, err := monitor.NewMemoryUsageMonitor()
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.Tick(*interval)
	for range ticker {

		stats, err := ramMon.Update()
		if err != nil {
			log.Printf("Failed to get memory usage: %v\n", err)

		}
		color = setColorByUsage(stats.MemoryUsage)
		fmt.Printf("Total Memory: %s\t Available Memory: %s\t Memory Usage: %s%.2f%%%s\n", normalizeBytes(float64(stats.TotalMemory)), normalizeBytes(float64(stats.AvailableMemory)), color, stats.MemoryUsage, ColorReset)
		cpuUsage, err := cpuMon.Update()
		if err != nil {
			log.Printf("Error reading CPU: %v", err)
		}
		color = setColorByUsage(cpuUsage)
		fmt.Printf("CPU Usage: %s%.2f%%%s\n", color, cpuUsage, ColorReset)
	}

}
