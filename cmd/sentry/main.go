package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/Pepegakac123/sentry-cli/pkg/monitor"
)

func main() {
	interval := flag.Duration("interval", 2*time.Second, "Check interval (e.g. 1s, 500ms)")
	flag.Parse()
	cpuMon, err := monitor.NewCPUMonitor()
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.Tick(*interval)
	for range ticker {
		ramUsage, err := monitor.GetRamUsage()
		if err != nil {
			fmt.Printf("Failed to get memory usage: %v\n", err)
			return
		}
		fmt.Printf("Total Memory: %d\t Available Memory: %d\t Memory Usage: %.2f%%\n", ramUsage.TotalMemory, ramUsage.AvailableMemory, ramUsage.MemoryUsage)
		cpuUsage, err := cpuMon.Update()
		if err != nil {
			log.Printf("Error reading CPU: %v", err)
			continue
		}
		fmt.Printf("CPU Usage: %.2f%%\n", cpuUsage)
	}

}
