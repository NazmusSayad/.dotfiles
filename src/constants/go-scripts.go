package constants

type BinScript struct {
	Exe           string
	StartMenuName string
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

	"gpg-unlock": {
		StartMenuName: "GPG Unlock",
	},

	"symlink-setup": {
		StartMenuName: "Symlink Setup",
	},

	"slack-status": {
		StartMenuName: "Slack Status",
		Exe:           "ss",
	},

	"msys-init": {
		StartMenuName: "MSYS2 Init",
	},

	"scoop-init": {
		StartMenuName: "Scoop Init",
	},

	"packages-sync": {
		StartMenuName: "Packages Sync",
		Exe:           "psy",
	},

	"winget-install": {
		StartMenuName: "WinGet Install",
		Exe:           "wgi",
	},

	"winget-upgrade": {
		StartMenuName: "WinGet Upgrade",
		Exe:           "wgu",
	},

	"clean-code-snippets": {
		StartMenuName: "Clean Code Snippets",
		Exe:           "csc",
	},

	"code-ext-sync": {
		StartMenuName: "Code Extensions Sync",
		Exe:           "ces",
	},
}
