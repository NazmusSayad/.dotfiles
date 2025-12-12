package winsdk

import (
	"os"
	"os/exec"
	"path/filepath"
)

func ConvertExeToRunAsAdmin(exe string) error {
	manifest := `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="requireAdministrator"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <!--The ID below indicates application support for Windows Vista -->
      <supportedOS Id="{e2011457-1546-43c5-a5fe-008deee3d3f0}"/>
      <!--The ID below indicates application support for Windows 7 -->
      <supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"/>
      <!--The ID below indicates application support for Windows 8 -->
      <supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
      <!--The ID below indicates application support for Windows 8.1 -->
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/> 
      <!--The ID below indicates application support for Windows 10 -->
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/> 
    </application>
  </compatibility>
</assembly>
`

	manifestPath := filepath.Join(filepath.Dir(exe), "admin.manifest")
	if err := os.WriteFile(manifestPath, []byte(manifest), 0644); err != nil {
		return err
	}

	cmd := exec.Command("mt.exe",
		"-manifest", manifestPath,
		"-outputresource:"+exe+";#1",
	)

	err := cmd.Run()
	os.Remove(manifestPath)
	return err
}
