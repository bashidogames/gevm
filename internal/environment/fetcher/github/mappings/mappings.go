package mappings

import "github.com/bashidogames/gevm/internal/platform"

type Mapping struct {
	System []string
	Arch   []string
}

var Mappings = map[platform.Platform]Mapping{
	platform.ExportTemplates: {
		System: []string{"export"},
		Arch:   []string{"templates"},
	},
	platform.DarwinArm64: {
		System: []string{"macos", "osx"},
		Arch:   []string{"universal"},
	},
	platform.DarwinAmd64: {
		System: []string{"macos", "osx"},
		Arch:   []string{"universal", "fat", "64"},
	},
	platform.Darwin386: {
		System: []string{"macos", "osx"},
		Arch:   []string{"universal", "fat", "32"},
	},
	platform.WindowsAmd64: {
		System: []string{"win"},
		Arch:   []string{"64"},
	},
	platform.Windows386: {
		System: []string{"win"},
		Arch:   []string{"32"},
	},
	platform.LinuxArm64: {
		System: []string{"linux", "x11"},
		Arch:   []string{"arm64"},
	},
	platform.LinuxArm: {
		System: []string{"linux", "x11"},
		Arch:   []string{"arm32"},
	},
	platform.LinuxAmd64: {
		System: []string{"linux", "x11"},
		Arch:   []string{"64"},
	},
	platform.Linux386: {
		System: []string{"linux", "x11"},
		Arch:   []string{"32"},
	},
}
