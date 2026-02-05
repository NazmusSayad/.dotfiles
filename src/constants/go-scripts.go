package constants

type BinScript struct {
	Exe             string
	StartMenuName   string
	NoProxySimulate bool
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
		StartMenuName:   "Slack Status",
		NoProxySimulate: true,
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

	"code-ext-sync": {
		StartMenuName: "Code Extensions Sync",
	},

	"packages-sync": {
		StartMenuName: "Packages Sync",
		Exe:           "psy",
	},
}
