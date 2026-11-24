package helpers

import (
	"fmt"
	"os"
	"runtime"

	"github.com/logrusorgru/aurora/v4"
)

type Machine struct {
	OS          string   `yaml:"OS"`
	WSL         bool     `yaml:"WSL"`
	Name        string   `yaml:"Name"`
	Alias       string   `yaml:"Alias"`
	Workstation bool     `yaml:"Workstation"`
	Offline     bool     `yaml:"Offline"`
	SSH         []string `yaml:"SSH"`
}

func ReadMachines() []Machine {
	var machines []Machine
	for _, machine := range ReadConfig[map[string]Machine]("#/machines.yml") {
		if !machine.Offline {
			machines = append(machines, machine)
		}
	}
	return machines
}

func CurrentMachine(machines []Machine) Machine {
	wsl := IsWSL()
	return onlyMachine(machines, "this "+runtime.GOOS+" system", func(machine Machine) bool {
		return machine.OS == runtime.GOOS && machine.WSL == wsl
	})
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
		fmt.Println(aurora.Red(fmt.Sprintf("Error: machines.yml has %d online machines for %s, expected 1", len(found), description)))
		os.Exit(1)
	}
	return found[0]
}
