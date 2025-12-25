package winget

func buildWingetOptions(p WingetPackage, interactive bool) []string {
	parts := []string{"--exact", "--id", p.ID}

	if p.Version != "" {
		parts = append(parts, "--version", p.Version)
	}

	if p.InstallerType != "" {
		parts = append(parts, "--installer-type", p.InstallerType)
	}

	if interactive {
		parts = append(parts, "--interactive")
	} else {
		parts = append(parts, "--silent")
	}

	return parts
}

func BuildWingetInstallArguments(p WingetPackage) []string {
	parts := []string{"install", "--verbose", "--accept-package-agreements", "--accept-source-agreements", "--no-upgrade"}
	return append(parts, buildWingetOptions(p, p.InteractiveInstall)...)
}

func BuildWingetUpgradeArguments(p WingetPackage) []string {
	parts := []string{"upgrade", "--verbose", "--accept-package-agreements", "--accept-source-agreements"}
	return append(parts, buildWingetOptions(p, p.InteractiveUpgrade)...)
}
