package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	introPath := filepath.Join(cwd, "./config/fish-config/__intro__.txt")

	cwdBase := filepath.Base(cwd)
	cpuInfo := getCPUInfo()
	ramInfo := getRAMInfo()
	diskInfo := getDiskInfo()
	currentTime := time.Now().Format("15:04:05")

	greeting := fmt.Sprintf(`📍 %s | 🕐 %s
🖥️  %s | 💾 %s
💿 %s
`, cwdBase, currentTime, cpuInfo, ramInfo, diskInfo)

	if err := os.WriteFile(introPath, []byte(greeting), 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Greeting updated: %s\n", introPath)
}

func getCPUInfo() string {
	cmd := exec.Command("powershell", "-c", "(Get-WmiObject Win32_Processor).Name")
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(output))
}

func getRAMInfo() string {
	cmd := exec.Command("powershell", "-c", "[math]::Round((Get-WmiObject Win32_ComputerSystem).TotalPhysicalMemory / 1GB, 0)")
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(output)) + "GB"
}

func getDiskInfo() string {
	cmd := exec.Command("powershell", "-c", "Get-WmiObject Win32_LogicalDisk -Filter 'DriveType=3' | ForEach-Object { $free=[math]::Round($_.FreeSpace/1GB,0); $total=[math]::Round($_.Size/1GB,0); Write-Output \"$($_.DeviceID) $free/$total GB\" }")
	output, err := cmd.Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(output))
}
