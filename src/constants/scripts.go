package constants

type BinScript struct {
	Exe       string
	StartMenu string
}

var BIN_SCRIPTS = map[string]BinScript{
	"git-clone": {
		Exe: "gcl",
	},

	"git-pull": {
		Exe: "gp",
	},

	"git-pull-all": {
		Exe: "gpa",
	},

	"git-pull-rebase": {
		Exe: "gpr",
	},

	"git-pull-merge": {
		Exe: "gpm",
	},

	"github-pr": {
		Exe: "ghp",
	},

	"slack-status": {
		Exe: "ss",
	},

	"gpg-unlock": {
		StartMenu: "GPG Unlock",
	},

	"symlink-init": {
		StartMenu: "Symlink Init",
	},

	"clean-code-snippets": {
		StartMenu: "Clean Code Snippets",
	},

	"opencode-server": {
		StartMenu: "OpenCode Server",
	},

	"opencode-configure": {
		StartMenu: "OpenCode Configure",
		Exe:       "oconfig",
	},

	"claude-pull": {
		StartMenu: "Claude Pull",
		Exe:       "ccp",
	},

	"claude-switch": {
		StartMenu: "Claude Switch",
		Exe:       "ccs",
	},

	"packages-sync": {
		StartMenu: "Packages Sync",
		Exe:       "psy",
	},

	"winget-install": {
		StartMenu: "WinGet Install",
		Exe:       "wgi",
	},

	"winget-upgrade": {
		StartMenu: "WinGet Upgrade",
		Exe:       "wgu",
	},
}
