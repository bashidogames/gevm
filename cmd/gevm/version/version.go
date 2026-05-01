package version

import (
	"runtime"
	"strings"

	"github.com/bashmills/gevm/internal/utils"
)

var version string

type Version struct{}

func (c *Version) Run() error {
	var semver string
	if len(strings.TrimSpace(version)) > 0 {
		semver = version
	} else {
		semver = "dev"
	}

	utils.Println(runtime.GOOS)
	utils.Println(runtime.GOARCH)
	utils.Println(semver)

	return nil
}
