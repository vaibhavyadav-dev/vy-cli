package feature

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func GetSysStatus() {
	// Get the system status
	getSystemInfo()

	// CPU usage using /proc/stat
	// getCPU()

	// Get number of running processes
	getProcesses()

	// Memory info using /proc/meminfo
	getMemory()

	// Disk usage using df command
	getStorage()
}

func getSystemInfo() {
	hostname, _ := exec.Command("hostname").Output()
	os, _ := exec.Command("uname", "-s").Output()
	platform, _ := exec.Command("uname", "-m").Output()

	// Define ANSI color codes
	headerColor := "\033[1;36m"    // Bright Cyan
	paramColor := "\033[1;33m"     // Bright Yellow
	borderColor := "\033[1;35m"    // Bright Magenta
	reset := "\033[0m"

	// Create formatted output
	fmt.Printf("%s %s  %-47s%s %s\n", borderColor, headerColor, "		System Information", borderColor, reset)
	border := fmt.Sprintf("%s+===========================================+%s\n", borderColor, reset)
	fmt.Print(border)	
	fmt.Printf("%s|%s %-12s%-30s%s|%s\n", 
		borderColor, headerColor, "Type", "Value", borderColor, reset)
	fmt.Print(border)
	
	fmt.Printf("%s|%s %-12s%-30s%s|%s\n",
		borderColor, paramColor, "Hostname", strings.TrimSpace(string(hostname)), borderColor, reset)
	fmt.Printf("%s|%s %-12s%-30s%s|%s\n",
		borderColor, paramColor, "OS", strings.TrimSpace(string(os)), borderColor, reset)
	fmt.Printf("%s|%s %-12s%-30s%s|%s\n",
		borderColor, paramColor, "Platform", strings.TrimSpace(string(platform)), borderColor, reset)
	fmt.Print(border)
}

// func getCPU(){
// 	cpuFile, _ := os.ReadFile("/proc/stat")
// 	cpuLines := strings.Split(string(cpuFile), "\n")
// 	if len(cpuLines) > 0 {
// 		cpuData := strings.Fields(cpuLines[0])
// 		if len(cpuData) > 1 {
// 			fmt.Printf("\nCPU Info: %s\n", cpuData[0])
// 		}
// 	}
// }

func getMemory() {
	memFile, _ := os.ReadFile("/proc/meminfo")
	memLines := strings.Split(string(memFile), "\n")
	
	// Define ANSI color codes
	headerColor := "\033[1;36m"    // Bright Cyan
	paramColor := "\033[1;33m"     // Bright Yellow
	borderColor := "\033[1;35m"    // Bright Magenta
	reset := "\033[0m"

	var memTotal, memFree float64
	
	// Parse memory information
	for _, line := range memLines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		valueKB, _ := strconv.ParseFloat(fields[1], 64)
		valueGB := valueKB / 1048576 // Convert KB to GB

		switch {
		case strings.HasPrefix(line, "MemTotal"):
			memTotal = valueGB
		case strings.HasPrefix(line, "MemFree"):
			memFree = valueGB
		// case strings.HasPrefix(line, "MemAvailable"):
			// memAvailable = valueGB
		}
	}

	// Create formatted output
	fmt.Printf("\n%s %s  %-47s%s %s\n", borderColor, headerColor, "		Memory Statistics", borderColor, reset)
	border := fmt.Sprintf("%s+=================================================+%s\n", borderColor, reset)
	fmt.Print(border)
	fmt.Printf("%s|%s %-12s%-12s%-12s%-12s%s|%s\n", 
		borderColor, headerColor, "Type", "Total", "Used", "Free", borderColor, reset)
	fmt.Print(border)
	
	usedMem := memTotal - memFree
	fmt.Printf("%s|%s %-12s%-12.2f%-12.2f%-12.2f%s|%s\n",
		borderColor, paramColor, "RAM", memTotal, usedMem, memFree, borderColor, reset)
	fmt.Print(border)
}

func getProcesses(){
	processes, _ := exec.Command("ps", "aux").Output()
	processLines := strings.Split(string(processes), "\n")
	fmt.Printf("\nProcess Info:\n")
	fmt.Printf("  Total Processes: %d\n", len(processLines)-1)

	// Sort processes by memory usage
	type Process struct {
		name   string
		memory float64
	}

	var topProcesses []Process
	for i, line := range processLines {
		if i == 0 { // Skip header
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 10 {
			continue
		}
		memPercent, _ := strconv.ParseFloat(fields[3], 64)
		topProcesses = append(topProcesses, Process{
			name:   fields[10],
			memory: memPercent,
		})
	}

	// overloading the sort interface
	sort.Slice(topProcesses, func(i, j int) bool {
		return topProcesses[i].memory > topProcesses[j].memory
	})


	// Define ANSI color codes
	headerColor := "\033[1;36m"    // Bright Cyan
	paramColor := "\033[1;33m"     // Bright Yellow
	borderColor := "\033[1;35m"    // Bright Magenta
	reset := "\033[0m"

	// Create a decorated border
	fmt.Printf("\n%s %s  %-47s%s %s\n", borderColor, headerColor, "	 Top 10 Memory Consuming Processes", borderColor, reset)
	border := fmt.Sprintf("%s+=================================================+%s\n", borderColor, reset)
	fmt.Print(border)
	fmt.Printf("%s|%s  %-32s%10s     %s|%s\n", borderColor, headerColor, "Process Name", "Memory %", borderColor, reset)
	fmt.Print(border)
	
	// Print process information
	for i := 0; i < 10 && i < len(topProcesses); i++ {
		fmt.Printf("%s|%s  %-29s%10.1f     %s   |%s\n", 
			borderColor, 
			paramColor, 
			topProcesses[i].name, 
			topProcesses[i].memory, 
			borderColor,
			reset)
	}
	fmt.Print(border)
}

func getStorage() {
	// Get disk usage using df command
	df, _ := exec.Command("df", "-h", "/").Output()
	lines := strings.Split(string(df), "\n")
	
	if len(lines) < 2 {
		return
	}
	
	// Define colors
	headerColor := "\033[1;36m"
	paramColor := "\033[1;33m"
	borderColor := "\033[1;35m"
	reset := "\033[0m"
	
	// Parse disk information (skip header)
	fields := strings.Fields(lines[1])
	if len(fields) < 6 {
		return
	}
	fmt.Printf("\n%s %s  %-47s%s %s\n", borderColor, headerColor, "	 	 Storage Stats", borderColor, reset)

	// Create formatted output
	border := fmt.Sprintf("%s+================================================+%s\n", borderColor, reset)
	fmt.Print(border)
	fmt.Printf("%s|%s  %-10s%-10s%-10s%-10s  %s    |%s\n", 
		borderColor, headerColor, "Total", "Used", "Free", "Use%", borderColor, reset)
	fmt.Print(border)
	fmt.Printf("%s|%s  %-10s%-10s%-10s%-10s  %s    |%s\n",
		borderColor, paramColor, fields[1], fields[2], fields[3], fields[4], borderColor, reset)
	fmt.Print(border)
}