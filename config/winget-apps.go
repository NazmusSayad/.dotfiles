package config

type Package struct {
	ID      string
	Name    string
	Version string
}

var Packages = []Package{
	// Basic Apps
	{ID: "MSYS2.MSYS2", Name: "MSYS2"},
	{ID: "ShareX.ShareX", Name: "ShareX"},
	{ID: "RARLab.WinRAR", Name: "WinRAR"},
	{ID: "Daum.PotPlayer", Name: "PotPlayer"},
	{ID: "9MSMLRH6LZF3", Name: "Notepad"},
	{ID: "9NBLGGH68TW4", Name: "Pictureflect"},
	{ID: "Piriform.CCleaner", Name: "CCleaner"},
	{ID: "ThioJoe.SvgThumbnailExtension", Name: "SVG Thumbnail"},

	// Productivity
	{ID: "Notion.Notion", Name: "Notion"},
	{ID: "OBSProject.OBSStudio", Name: "OBS Studio"},
	{ID: "PowerSoftware.AnyBurn", Name: "AnyBurn"},
	{ID: "UnifiedIntents.UnifiedRemote", Name: "Unified Remote"},

	// Browsers
	{ID: "Google.Chrome", Name: "Google Chrome"},
	{ID: "Mozilla.Firefox.DeveloperEdition", Name: "Firefox Developer Edition"},

	// Code Editors
	{ID: "Python.PythonInstallManager", Name: "Python Installer"},
	{ID: "CoreyButler.NVMforWindows", Name: "nvm for Windows"},
	{ID: "Microsoft.VisualStudioCode", Name: "VS Code"},
	{ID: "Anysphere.Cursor", Name: "Cursor"},

	// Programming Tools
	{ID: "Git.Git", Name: "Git"},
	{ID: "GitHub.cli", Name: "GitHub CLI"},
	{ID: "GoLang.Go", Name: "Go"},

	// Other
	{ID: "Postman.Postman", Name: "Postman"},
	{ID: "ApacheFriends.Xampp.8.2", Name: "XAMPP"},
	{ID: "PostgreSQL.PostgreSQL.17", Name: "PostgreSQL"},

	// Communication
	{ID: "SlackTechnologies.Slack", Name: "Slack"},
}
