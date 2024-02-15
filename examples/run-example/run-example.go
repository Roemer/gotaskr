package main

import (
	"os"

	"github.com/roemer/gotaskr"
	"github.com/roemer/gotaskr/goext"
	"github.com/roemer/gotaskr/log"
)

func init() {
	gotaskr.Task("Run-In-Directory", func() error {
		pwd, _ := os.Getwd()
		log.Informationf("Path before: %s", pwd)

		err := goext.RunInDirectory("subdir", func() error {
			pwd, _ = os.Getwd()
			log.Informationf("Path inside: %s", pwd)
			return nil
		})

		pwd, _ = os.Getwd()
		log.Informationf("Path after : %s", pwd)

		return err
	})

	gotaskr.Task("Run-With-Variables", func() error {
		log.Informationf("Variable before: %s", os.Getenv("TEST"))

		err := goext.RunWithEnvs(map[string]string{"TEST": "foo"}, func() error {
			log.Informationf("Variable inside: %s", os.Getenv("TEST"))
			return nil
		})

		log.Informationf("Variable after : %s", os.Getenv("TEST"))

		return err
	})

	gotaskr.Task("Run-With-Multiple-Options", func() error {
		return goext.RunWithOptions(func() error {
			return nil
		}, goext.RunOptionInDirectory("subdir"), goext.RunOptionWithEnvs(map[string]string{"TEST": "foo"}))
	})
}

func main() {
	os.Exit(gotaskr.Execute())
}
