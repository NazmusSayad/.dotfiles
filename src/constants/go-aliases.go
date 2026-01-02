package constants

var BIN_ALIASES = map[string][]string{
	"gc":  {"git", "checkout"},
	"gcn": {"git", "checkout", "-b"},
	"gds": {"git", "diff", "--stat"},

	"ghv": {"gh", "repo", "view"},
	"ghw": {"gh", "repo", "view", "--web"},
	"ghp": {"gh", "pr", "create", "-B"},

	"fsc": {"fsutil.exe", "file", "setCaseSensitiveInfo", ".", "enable", "recursive"},
}
