package cli

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/dannrocha/czen/src/gitscm"
	"github.com/dannrocha/czen/src/semver"
	"github.com/dannrocha/czen/src/setup"
	"github.com/urfave/cli/v2"
)

func Bump(c *cli.Context) error {

	script := setup.Script{}
	script.LoadScript()

	for _, auto := range script.Automation {
		if auto.Bind == BUMP && auto.Enable {
			if auto.When == setup.BEFORE {
				auto.Run()
			} else {
				defer auto.Run()
			}
		}
	}

	git, err := gitscm.New()

	if err != nil {
		fmt.Println(err)
	}

	lastestTag, ok := git.LastestTag()

	if !ok {
		return nil
	}

	newVersion := incrementVersion(lastestTag)


	annonation := newVersion.ConvertToSemver().Version
	hash := md5.Sum([]byte(fmt.Sprintf("%v-%v", annonation, time.Now())))

	ok, err = gitscm.CreateTag(annonation, fmt.Sprintf("%x", hash))

	if !ok || err != nil {
		return nil
	}

	return nil
}

func incrementVersion(tag gitscm.GitTag) semver.Version {

	config := setup.Configuration{}
	config.LoadConfigurationFile()

	profile, errProfile := config.FindCurrentProfileEnable()

	if errProfile != nil {
		panic(errProfile)
	}

	commits, errCommit := gitscm.LoadCommitsFrom(tag.Commit.Hash)

	if errCommit != nil {
		panic(errCommit.Error())
	}

	currentSemever := semver.New(tag.Annotation)
	oldVersion, errVersion := currentSemever.FindVersion()

	if errVersion != nil {
		panic(errVersion.Error())
	}

	newVersion := semver.Version{
		Major: oldVersion.Major,
		Minor: oldVersion.Minor,
		Path:  oldVersion.Path,
	}

	for _, commit := range commits {

		for context, pattern := range profile.BumpMap {
			if strings.Contains(commit.Message, context) {
				newVersion.IncrementVersion(pattern)
				break
			}
		}
	}

	return newVersion
}