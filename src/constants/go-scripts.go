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

	"code-ext-sync": {
		StartMenu: "Code Extensions Sync",
	},

	"code-state-pull": {
		StartMenu: "Code UI State Pull",
	},

	"code-state-push": {
		StartMenu: "Code UI State Push",
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

	"clean-code-snippets": {
		StartMenu: "Clean Code Snippets",
		Exe:       "csc",
	},
}
