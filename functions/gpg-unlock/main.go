package main

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

func killByName(name string) {
	out, _ := run("ps", "aux")
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, name) && !strings.Contains(line, "grep") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				pid, err := strconv.Atoi(fields[1])
				if err == nil {
					exec.Command("kill", "-9", strconv.Itoa(pid)).Run()
				}
			}
		}
	}
}

func main() {
	killByName("gpg")
	killByName("keyboxd")
	locks, _ := filepath.Glob(filepath.Join(os.Getenv("HOME"), ".gnupg", "*.lock"))
	for _, lf := range locks {
		os.Remove(lf)
	}
}

