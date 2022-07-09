package semver

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dannrocha/czen/src/setup"
)

func (version *Version) IncrementVersion(level string) {

	increment := map[string]func(Version) Version {
		"MAJOR": func(v Version) Version {
			return Version {
				Major: v.Major + 1,
			}
		},

		"MINOR": func(v Version) Version {
			return Version {
				Major: v.Major,
				Minor: v.Minor + 1,
			}
		},

		"PATCH": func(v Version) Version {
			return Version {
				Major: v.Major,
				Minor: v.Minor,
				Path: v.Path + 1,
			}
		},
	}

	level = strings.ToUpper(level)

	*version = increment[level](*version)
}

func (version Version) ConvertToSemver() SemVer {
	config := setup.Configuration{}
	config.LoadConfigurationFile()

	envs := map[string] string {
		"$major": strconv.Itoa(version.Major),
		"$minor": strconv.Itoa(version.Minor),
		"$patch": strconv.Itoa(version.Path),
		"$version": fmt.Sprintf("%v.%v.%v", version.Major, version.Minor, version.Path),
	}

	profile, errProfile := config.FindCurrentProfileEnable()

	if errProfile != nil {
		panic(errProfile)
	}

	format := profile.TagFormat

	for match, content := range envs {
		format = strings.Replace(format, match, content, -1)
	}

	return SemVer{
		Version: format,
	}
}