package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/logrusorgru/aurora/v4"
)

func main() {
	if os.Getenv("OPENCODE_SERVER_PASSWORD") == "" {
		fmt.Println(aurora.Red("OPENCODE_SERVER_PASSWORD is not set"))
		os.Exit(1)
	}

	processes := []*exec.Cmd{
		exec.Command("opencode", "serve", "--port", "4747", "--cors", "https://oc.sayad.dev"),
		exec.Command("cloudflared", "tunnel", "run", "opencode"),
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	exited := make(chan string, len(processes))
	started := 0
	for _, cmd := range processes {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			fmt.Println(aurora.Red("failed to start "+cmd.Args[0]+":"), err)
			stopAll(processes, exited, started)
			os.Exit(1)
		}
		started++
		go func() {
			cmd.Wait()
			exited <- cmd.Args[0]
		}()
	}

	fmt.Println(aurora.Green("opencode server running at"), aurora.Cyan("https://oc.sayad.dev"))
	fmt.Println(aurora.Faint("> Press ^C to stop"))

	select {
	case name := <-exited:
		fmt.Println(aurora.Yellow(name + " exited, stopping the rest..."))
		stopAll(processes, exited, started-1)
		os.Exit(1)
	case <-sig:
		fmt.Println(aurora.Yellow("Stopping..."))
		stopAll(processes, exited, started)
	}
}

func stopAll(processes []*exec.Cmd, exited chan string, running int) {
	for _, cmd := range processes {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}
	for range running {
		<-exited
	}
}
