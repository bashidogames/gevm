package platform

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
