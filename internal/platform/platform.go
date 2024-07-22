package platform

import (
	"fmt"
	"runtime"
)

type Platform string

const (
	ExportTemplates Platform = "Export Templates"
	WindowsAmd64    Platform = "Windows Amd64"
	Windows386      Platform = "Windows 386"
	LinuxArm64      Platform = "Linux Arm64"
	LinuxAmd64      Platform = "Linux Amd64"
	LinuxArm        Platform = "Linux Arm"
	Linux386        Platform = "Linux 386"
	DarwinArm64     Platform = "Darwin Arm64"
	DarwinAmd64     Platform = "Darwin Amd64"
	Darwin386       Platform = "Darwin 386"
)

var Platforms = []Platform{
	ExportTemplates,
	WindowsAmd64,
	Windows386,
	LinuxArm64,
	LinuxAmd64,
	LinuxArm,
	Linux386,
	DarwinArm64,
	DarwinAmd64,
	Darwin386,
}

func Get() (Platform, error) {
	switch runtime.GOOS {
	case "windows":
		switch runtime.GOARCH {
		case "amd64":
			return WindowsAmd64, nil
		case "386":
			return Windows386, nil
		}
	case "linux":
		switch runtime.GOARCH {
		case "arm64":
			return LinuxArm64, nil
		case "amd64":
			return LinuxAmd64, nil
		case "arm":
			return LinuxArm, nil
		case "386":
			return Linux386, nil
		}
	case "darwin":
		switch runtime.GOARCH {
		case "arm64":
			return DarwinArm64, nil
		case "amd64":
			return DarwinAmd64, nil
		case "386":
			return Darwin386, nil
		}
	}

	return "", fmt.Errorf("invalid platform")
}
