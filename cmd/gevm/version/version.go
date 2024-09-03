package version

import (
	"runtime"
	"strings"

	"github.com/bashidogames/gevm/internal/utils"
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

	utils.Printlnf(runtime.GOOS)
	utils.Printlnf(runtime.GOARCH)
	utils.Printlnf(semver)

	return nil
}
