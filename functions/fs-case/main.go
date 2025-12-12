package main

import "os/exec"

func main() {
	exec.Command("fsutil.exe", "file", "setCaseSensitiveInfo", ".", "enable", "recursive").Run()
}

