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

	"gh-pull-create": {
		Exe: "gpc",
	},

	"slack-status": {
		Exe: "ss",
	},

	"slack-startup": {
		Exe: "sst",
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

	"clean-code-snippets": {
		StartMenu: "Clean Code Snippets",
	},

	"opencode-models": {
		StartMenu: "OpenCode Models Sync",
	},

	"opencode-server": {
		StartMenu: "OpenCode Server",
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
