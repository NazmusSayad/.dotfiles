package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func run(cmd string, args ...string) (string, int) {
	var b bytes.Buffer
	c := exec.Command(cmd, args...)
	c.Stdout = &b
	c.Stderr = &b
	err := c.Run()
	if err != nil {
		if e, ok := err.(*exec.ExitError); ok {
			return b.String(), e.ExitCode()
		}
		return b.String(), 1
	}
	return b.String(), 0
}

func main() {
	fmt.Println("Restore and clean?")
	fmt.Print("Press [Enter] to confirm, or any other key to cancel: ")
	r := bufio.NewReader(os.Stdin)
	input, _ := r.ReadString('\n')
	if input == "\n" {
		run("git", "restore", ".")
		run("git", "clean", "-fd")
		return
	}
	fmt.Println("❌ Aborted.")
}

