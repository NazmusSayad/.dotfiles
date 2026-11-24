package helpers

import (
	"errors"
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"charm.land/huh/v2"
	"github.com/charmbracelet/x/term"
	"github.com/logrusorgru/aurora/v4"
)

const machineFile = "@/.local/machine"

type Machine struct {
	ID          string   `yaml:"-"`
	OS          string   `yaml:"OS"`
	WSL         bool     `yaml:"WSL"`
	Name        string   `yaml:"Name"`
	Workstation bool     `yaml:"Workstation"`
	Offline     bool     `yaml:"Offline"`
	SSH         []string `yaml:"SSH"`
}

func ReadMachines() []Machine {
	var machines []Machine
	for id, machine := range ReadConfig[map[string]Machine]("#/machines.yml") {
		if !machine.Offline {
			machine.ID = id
			machines = append(machines, machine)
		}
	}
	return machines
}

func CurrentMachine(machines []Machine) Machine {
	path := ResolvePath(machineFile)
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		exitWithError("failed to read " + path + ": " + err.Error())
	}
	byID := map[string]Machine{}
	for _, machine := range machines {
		byID[machine.ID] = machine
	}
	saved := strings.TrimSpace(string(data))
	if machine, found := byID[saved]; found {
		return machine
	}

	if saved != "" {
		fmt.Println(aurora.Yellow(path + " says " + saved + ", which is not an online machine in machines.yml"))
	}
	if !term.IsTerminal(os.Stdin.Fd()) {
		exitWithError("no machine picked in " + path + ". Run any dotfiles script in a terminal on this machine to pick it")
	}

	ids := slices.Sorted(maps.Keys(byID))
	var picked string
	err = huh.NewSelect[string]().
		Title("Which machine is this? Offline machines are not listed").
		Options(huh.NewOptions(ids...)...).
		Value(&picked).
		WithTheme(HuhTheme()).
		Run()
	if errors.Is(err, huh.ErrUserAborted) {
		exitWithError("no machine picked")
	}
	if err != nil {
		exitWithError("failed to pick a machine: " + err.Error())
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		exitWithError("failed to create " + filepath.Dir(path) + ": " + err.Error())
	}
	if err := os.WriteFile(path, []byte(picked+"\n"), 0o644); err != nil {
		exitWithError("failed to save " + path + ": " + err.Error())
	}
	fmt.Println(aurora.Faint("Saved " + picked + " to " + path))
	return byID[picked]
}

func WorkstationMachine(machines []Machine) Machine {
	return onlyMachine(machines, "Workstation", func(machine Machine) bool {
		return machine.Workstation
	})
}

func onlyMachine(machines []Machine, description string, match func(Machine) bool) Machine {
	var found []Machine
	for _, machine := range machines {
		if match(machine) {
			found = append(found, machine)
		}
	}
	if len(found) != 1 {
		exitWithError(fmt.Sprintf("machines.yml has %d online machines for %s, expected 1", len(found), description))
	}
	return found[0]
}

func exitWithError(message string) {
	fmt.Println(aurora.Red("Error: " + message))
	os.Exit(1)
}
