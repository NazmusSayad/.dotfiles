package constants

const SOURCE_DIR = "./src"
const BUILD_DIR = "./.build"

const SCRIPTS_SOURCE_DIR = SOURCE_DIR + "/scripts"
const SCRIPTS_BUILD_BIN_DIR = BUILD_DIR + "/bin"

type Script struct {
	Exe           string
	StartMenuName string
}

var SCRIPTS_MAP = map[string]Script{
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

	"gh-repo-view": {
		Exe: "ghv",
	},

	"gh-repo-view-web": {
		Exe: "ghw",
	},

	"gh-pr-create": {
		Exe: "ghp",
	},

	"file-sys-case": {
		Exe: "fscs",
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

	"winget-install": {
		StartMenuName: "WinGet Install",
	},

	"winget-upgrade": {
		StartMenuName: "WinGet Upgrade",
	},
}
