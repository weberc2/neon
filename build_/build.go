package build

import "bitbucket.org/weberc2/neon/config"

var (
	buildErrors ErrorClass
)

func Build() error {
	buildErrors = NewErrorClass("build")

	conf, err := config.Load()
	if err != nil {
		return err
	}

	return BuildSite(conf)
}
