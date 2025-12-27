package main

import "fmt"

type StructName struct {
	Name      string
	Arguments []string
}

var ALIASES = []StructName{
	{
		Name:      "fsc",
		Arguments: []string{"fsutil.exe", "file", "setCaseSensitiveInfo", ".", "enable", "recursive"},
	},

	{
		Name:      "ghp",
		Arguments: []string{"gh", "pr", "create", "-B"},
	},

	{
		Name:      "ghv",
		Arguments: []string{"gh", "repo", "view"},
	},

	{
		Name:      "ghw",
		Arguments: []string{"gh", "repo", "view", "--web"},
	},

	{
		Name:      "gds",
		Arguments: []string{"git", "diff", "--stat"},
	},

	{
		Name:      "gp",
		Arguments: []string{"git", "pull"},
	},
}

func main() {
	for _, alias := range ALIASES {
		fmt.Println(alias.Name)
	}
}
