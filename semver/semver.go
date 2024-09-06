package semver

import (
	"cmp"
	"errors"
	"fmt"
	"regexp"

	"github.com/bashidogames/gevm/internal/utils"
)

const VERSION_REGEX_PATTERN = "([1-9][0-9]*|0)[.]([1-9][0-9]*|0)([.]([1-9][0-9]*|0))?([.]([1-9][0-9]*|0))?"
const RELEASE_REGEX_PATTERN = "((dev|alpha|beta|rc)([1-9][0-9]*|0)|stable)([-_.](unofficial))?"
const RELVER_REGEX_PATTERN = "(" + VERSION_REGEX_PATTERN + ")[-_.](" + RELEASE_REGEX_PATTERN + ")"
const SEMVER_REGEX_PATTERN = RELVER_REGEX_PATTERN + "([-_.](mono))?"

var VersionRegex = regexp.MustCompile(VERSION_REGEX_PATTERN)
var ReleaseRegex = regexp.MustCompile(RELEASE_REGEX_PATTERN)
var RelverRegex = regexp.MustCompile(RELVER_REGEX_PATTERN)
var SemverRegex = regexp.MustCompile(SEMVER_REGEX_PATTERN)

var ErrRegexFailed = errors.New("regex failed")

var Labels = map[string]int{
	"dev":    1,
	"alpha":  2,
	"beta":   3,
	"rc":     4,
	"stable": 5,
}

type Version struct {
	Original string
	Major    int
	Minor    int
	Patch    int
	Build    int
}

func (v Version) IsValid() bool {
	return len(v.Original) > 0
}

func (a Version) Compare(b Version) int {
	if result := cmp.Compare(a.Major, b.Major); result != 0 {
		return result
	}

	if result := cmp.Compare(a.Minor, b.Minor); result != 0 {
		return result
	}

	if result := cmp.Compare(a.Patch, b.Patch); result != 0 {
		return result
	}

	if result := cmp.Compare(a.Build, b.Build); result != 0 {
		return result
	}

	return 0
}

func (a Version) GreaterOrEqual(b Version) bool {
	return a.Compare(b) >= 0
}

func (a Version) LessOrEqual(b Version) bool {
	return a.Compare(b) <= 0
}

func (a Version) Greater(b Version) bool {
	return a.Compare(b) > 0
}

func (a Version) Less(b Version) bool {
	return a.Compare(b) < 0
}

func (a Version) Equal(b Version) bool {
	return a.Compare(b) == 0
}

func (v Version) String() string {
	return v.Original
}

func ParseVersion(version string) (Version, error) {
	parts := VersionRegex.FindStringSubmatch(version)
	if parts == nil {
		return Version{}, fmt.Errorf("invalid version string: %w: %s", ErrRegexFailed, version)
	}

	major, err := utils.AtoiIfNotEmpty(parts[1])
	if err != nil {
		return Version{}, fmt.Errorf("invalid major: %w", err)
	}

	minor, err := utils.AtoiIfNotEmpty(parts[2])
	if err != nil {
		return Version{}, fmt.Errorf("invalid minor: %w", err)
	}

	patch, err := utils.AtoiIfNotEmpty(parts[4])
	if err != nil {
		return Version{}, fmt.Errorf("invalid patch: %w", err)
	}

	build, err := utils.AtoiIfNotEmpty(parts[6])
	if err != nil {
		return Version{}, fmt.Errorf("invalid build: %w", err)
	}

	return Version{
		Original: version,
		Major:    major,
		Minor:    minor,
		Patch:    patch,
		Build:    build,
	}, nil
}

type Release struct {
	Original string
	Label    string
	Digit    int
	Meta     string
}

func (r Release) IsValid() bool {
	return len(r.Original) > 0
}

func (a Release) Compare(b Release) int {
	if result := cmp.Compare(Labels[a.Label], Labels[b.Label]); result != 0 {
		return result
	}

	if result := cmp.Compare(a.Digit, b.Digit); result != 0 {
		return result
	}

	return 0
}

func (a Release) GreaterOrEqual(b Release) bool {
	return a.Compare(b) >= 0
}

func (a Release) LessOrEqual(b Release) bool {
	return a.Compare(b) <= 0
}

func (a Release) Greater(b Release) bool {
	return a.Compare(b) > 0
}

func (a Release) Less(b Release) bool {
	return a.Compare(b) < 0
}

func (a Release) Equal(b Release) bool {
	return a.Compare(b) == 0
}

func (r Release) String() string {
	return r.Original
}

func ParseRelease(release string) (Release, error) {
	parts := ReleaseRegex.FindStringSubmatch(release)
	if parts == nil {
		return Release{}, fmt.Errorf("invalid release string: %w: %s", ErrRegexFailed, release)
	}

	label, err := utils.SelectFirstNotEmpty(parts[2], parts[1])
	if err != nil {
		return Release{}, fmt.Errorf("invalid label: %w", err)
	}

	digit, err := utils.AtoiIfNotEmpty(parts[3])
	if err != nil {
		return Release{}, fmt.Errorf("invalid digit: %w", err)
	}

	meta := parts[5]

	return Release{
		Original: release,
		Label:    label,
		Digit:    digit,
		Meta:     meta,
	}, nil
}

type Relver struct {
	Version Version
	Release Release
}

func (r Relver) IsValid() bool {
	return r.Version.IsValid() && r.Release.IsValid()
}

func (a Relver) Compare(b Relver) int {
	if result := a.Version.Compare(b.Version); result != 0 {
		return result
	}

	if result := a.Release.Compare(b.Release); result != 0 {
		return result
	}

	return 0
}

func (a Relver) GreaterOrEqual(b Relver) bool {
	return a.Compare(b) >= 0
}

func (a Relver) LessOrEqual(b Relver) bool {
	return a.Compare(b) <= 0
}

func (a Relver) Greater(b Relver) bool {
	return a.Compare(b) > 0
}

func (a Relver) Less(b Relver) bool {
	return a.Compare(b) < 0
}

func (a Relver) Equal(b Relver) bool {
	return a.Compare(b) == 0
}

func (vr Relver) ExportTemplatesString() string {
	return fmt.Sprintf("%s.%s", vr.Version, vr.Release)
}

func (vr Relver) GodotString() string {
	return fmt.Sprintf("%s-%s", vr.Version, vr.Release)
}

func ParseRelver(relver string) (Relver, error) {
	parts := SemverRegex.FindStringSubmatch(relver)
	if parts == nil {
		return Relver{}, fmt.Errorf("invalid relver string: %w: %s", ErrRegexFailed, relver)
	}

	version := parts[1]
	release := parts[8]

	return NewRelver(version, release)
}

func NewRelver(version string, release string) (Relver, error) {
	v, err := ParseVersion(version)
	if err != nil {
		return Relver{}, fmt.Errorf("invalid version: %w", err)
	}

	r, err := ParseRelease(release)
	if err != nil {
		return Relver{}, fmt.Errorf("invalid release: %w", err)
	}

	return Relver{
		Version: v,
		Release: r,
	}, nil
}

func MaybeRelver(version string, release string) Relver {
	relver, err := NewRelver(version, release)
	if err != nil {
		return Relver{}
	}

	return relver
}

type Semver struct {
	Relver Relver
	Mono   bool
}

func (s Semver) IsValid() bool {
	return s.Relver.IsValid()
}

func (a Semver) Compare(b Semver) int {
	return a.Relver.Compare(b.Relver)
}

func (a Semver) GreaterOrEqual(b Semver) bool {
	return a.Compare(b) >= 0
}

func (a Semver) LessOrEqual(b Semver) bool {
	return a.Compare(b) <= 0
}

func (a Semver) Greater(b Semver) bool {
	return a.Compare(b) > 0
}

func (a Semver) Less(b Semver) bool {
	return a.Compare(b) < 0
}

func (a Semver) Equal(b Semver) bool {
	return a.Compare(b) == 0
}

func (s Semver) ExportTemplatesString() string {
	if s.Mono {
		return fmt.Sprintf("%s.mono", s.Relver.ExportTemplatesString())
	} else {
		return s.Relver.ExportTemplatesString()
	}
}

func (s Semver) GodotString() string {
	if s.Mono {
		return fmt.Sprintf("%s-mono", s.Relver.GodotString())
	} else {
		return s.Relver.GodotString()
	}
}

func Parse(semver string) (Semver, error) {
	parts := SemverRegex.FindStringSubmatch(semver)
	if parts == nil {
		return Semver{}, fmt.Errorf("invalid semver string: %w: %s", ErrRegexFailed, semver)
	}

	mono := len(parts[15]) > 0
	version := parts[1]
	release := parts[8]

	return New(version, release, mono)
}

func New(version string, release string, mono bool) (Semver, error) {
	relver, err := NewRelver(version, release)
	if err != nil {
		return Semver{}, fmt.Errorf("invalid relver: %w", err)
	}

	return Semver{
		Relver: relver,
		Mono:   mono,
	}, nil
}

func Maybe(version string, release string, mono bool) Semver {
	semver, err := New(version, release, mono)
	if err != nil {
		return Semver{}
	}

	return semver
}
