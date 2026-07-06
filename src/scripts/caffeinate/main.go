package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	ES_CONTINUOUS       = 0x80000000
	ES_SYSTEM_REQUIRED  = 0x00000001
	ES_DISPLAY_REQUIRED = 0x00000002
)

func main() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	setThreadExecutionState := kernel32.NewProc("SetThreadExecutionState")

	setThreadExecutionState.Call(
		uintptr(ES_CONTINUOUS | ES_SYSTEM_REQUIRED | ES_DISPLAY_REQUIRED),
	)

	start := time.Now()
	fmt.Printf("[%s] caffeinate on: sleep and display blocked\n", start.Format("15:04:05"))
	fmt.Println("> Press ^C to stop")

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case now := <-ticker.C:
			fmt.Printf("[%s] awake for %s\n", now.Format("15:04:05"), now.Sub(start).Round(time.Second))
		case <-sig:
			setThreadExecutionState.Call(uintptr(ES_CONTINUOUS))
			fmt.Printf("[%s] caffeinate off: sleep restored after %s\n", time.Now().Format("15:04:05"), time.Since(start).Round(time.Second))
			return
		}
	}
}
