package windows_admin

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertExeToRunAsAdmin(exe string) error {
	manifest := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <assemblyIdentity version="1.0.0.0" processorArchitecture="*" name="converted.app" type="win32"/>
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
</assembly>`

	manifestPath := filepath.Join(filepath.Dir(exe), "admin.manifest")

	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		return err
	}

	cmd := exec.Command("mt.exe",
		"-manifest", manifestPath,
		"-outputresource:"+exe+";#1",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
