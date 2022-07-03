package semver

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func New(Version string) SemVer {
	return SemVer{
		Version,
	}
}

func (sem SemVer) FindVersion() (Version, error) {

	if !sem.IsVersionValid() {
		return Version{}, errors.New("Semantic version is invalid")
	}

	regex := regexp.MustCompile(`(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)`)

	version := strings.Split(regex.FindString(sem.Version), ".")
	major, majorErr := strconv.Atoi(version[MAJOR])
	minor, minorErr := strconv.Atoi(version[MINOR])
	path, pathErr := strconv.Atoi(version[PATCH])

	if majorErr != nil || minorErr != nil || pathErr != nil {
		return Version{}, errors.New("Error converting version to integer")
	}

	return Version{
		Major: major,
		Minor: minor,
		Path:  path,
	}, nil
}

func (sem SemVer) IsVersionValid() bool {
	return regexp.MustCompile(SEMVER_REGEX).MatchString(sem.Version)
}
