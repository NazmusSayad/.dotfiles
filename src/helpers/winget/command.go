package winget

func buildWingetOptions(p WingetPackage, interactive bool) []string {
	parts := []string{"--exact", "--id", p.ID, "--verbose", "--accept-package-agreements", "--accept-source-agreements"}

	if p.Scope != "" {
		parts = append(parts, "--scope", p.Scope)
	}

	if p.Version != "" {
		parts = append(parts, "--version", p.Version)
	}

	if p.InstallerType != "" {
		parts = append(parts, "--installer-type", p.InstallerType)
	}

	if p.SkipDependencies {
		parts = append(parts, "--skip-dependencies")
	}

	if interactive {
		parts = append(parts, "--interactive")
	} else {
		parts = append(parts, "--silent")
	}

	return parts
}

func BuildWingetInstallArguments(p WingetPackage) []string {
	parts := []string{"install", "--no-upgrade"}
	return append(parts, buildWingetOptions(p, p.InteractiveInstall)...)
}

func BuildWingetUpgradeArguments(p WingetPackage) []string {
	parts := []string{"upgrade"}
	return append(parts, buildWingetOptions(p, p.InteractiveUpgrade)...)
}
