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
		Exe: "ghp",
	},

	"gpg-unlock": {
		StartMenuName: "GPG Unlock",
	},

	"clean-code-snippets": {
		StartMenuName: "Clean Code Snippets",
	},

	"msys-setup": {
		StartMenuName: "MSYS2 Setup",
	},

	"symlink-setup": {
		StartMenuName: "Symlink Setup",
	},

	"slack-status": {
		StartMenuName: "Slack Status",
	},

	"scoop-init": {
		StartMenuName: "Scoop Init",
	},

	"winget-install": {
		StartMenuName: "WinGet Install",
	},

	"winget-upgrade": {
		StartMenuName: "WinGet Upgrade",
	},

	"packages-sync": {
		StartMenuName: "Packages Sync",
	},

	"sync-code-ext": {
		StartMenuName: "Sync Code Extensions",
	},
}
