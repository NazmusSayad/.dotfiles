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

	"gpg-unlock": {
		StartMenu: "GPG Unlock",
	},

	"symlink-setup": {
		StartMenu: "Symlink Setup",
	},

	"msys-init": {
		StartMenu: "MSYS2 Init",
	},

	"scoop-init": {
		StartMenu: "Scoop Init",
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

	"code-ext-sync": {
		StartMenu: "Code Extensions Sync",
		Exe:       "ces",
	},
}
