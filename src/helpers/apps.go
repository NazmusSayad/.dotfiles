package helpers

type MsysAppConfig struct {
	ID   string
	Repo string
}

type ScoopAppConfig struct {
	ID     string
	Bucket string
}

type WingetAppConfig struct {
	ID            string
	Name          string
	Scope         string
	Version       string
	InstallerType string

	ForceAdminInstall bool
	ForceAdminUpgrade bool

	InteractiveInstall bool
	InteractiveUpgrade bool

	SkipInstall      bool
	SkipUpgrade      bool
	SkipDependencies bool
}

// func GetMsysApps() []MsysAppConfig {}

// func GetScoopApps() []ScoopAppConfig {}

// func GetWingetApps() []WingetAppConfig {}
