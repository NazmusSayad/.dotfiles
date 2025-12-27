package constants

const SOURCE_DIR = "./src"
const BUILD_DIR = "./.build"

const SCRIPTS_SOURCE_DIR = SOURCE_DIR + "/scripts"

const BUILD_TEMP_DIR = BUILD_DIR + "/temp"
const BUILD_SCRIPTS_DIR = BUILD_DIR + "/bin"
const BUILD_TASKS_RUNNER_DIR = BUILD_DIR + "/tasks"

type BinScript struct {
	Exe           string
	StartMenuName string
}

var BIN_ALIASES = map[string][]string{
	"r":   {"nr"},
	"nid": {"ni -D"},

	"gp":  {"git", "pull"},
	"gds": {"git", "diff", "--stat"},

	"ghv": {"gh", "repo", "view"},
	"ghw": {"gh", "repo", "view", "--web"},
	"ghp": {"gh", "pr", "create", "-B"},

	"fsc": {"fsutil.exe", "file", "setCaseSensitiveInfo", ".", "enable", "recursive"},
}

var BIN_SCRIPTS = map[string]BinScript{
	"git-clone": {
		Exe: "c",
	},

	"git-pull-rebase": {
		Exe: "gpr",
	},

	"git-pull-merge": {
		Exe: "gpm",
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
