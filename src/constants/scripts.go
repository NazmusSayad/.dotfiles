package constants

type BinScript struct {
	Exe       string
	StartMenu string
}

var BIN_SCRIPTS = map[string]BinScript{
	"git-clone": {
		Exe: "c",
	},

	"git-checkout": {
		Exe: "gc",
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

	"github-pr-create": {
		Exe: "ghp",
	},

	"slack-status": {
		Exe: "ss",
	},

	"gpg-unlock": {
		StartMenu: "GPG Unlock",
	},

	"symlink-setup": {
		StartMenu: "Symlink Setup",
	},

	"code-init": {
		StartMenu: "Code Init",
	},

	"code-ext-sync": {
		StartMenu: "Code Extensions Sync",
	},

	"code-state-pull": {
		StartMenu: "Code UI State Pull",
	},

	"code-state-push": {
		StartMenu: "Code UI State Push",
	},

	"opencode-server": {
		StartMenu: "OpenCode Server",
	},

	"opencode-models": {
		StartMenu: "OpenCode Models Sync",
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
