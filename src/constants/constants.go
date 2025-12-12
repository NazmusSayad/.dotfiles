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

	"file-sys-case": {
		Exe: "fscs",
	},

	"gpg-unlock": {
		StartMenuName: "GPG Unlock",
	},

	"clean-code-snippets": {
		StartMenuName: "Clean Code Snippets",
	},

	"symlink-setup": {
		StartMenuName: "Symlink Setup",
	},

	"slack-startup": {
		StartMenuName: "Slack Startup",
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
