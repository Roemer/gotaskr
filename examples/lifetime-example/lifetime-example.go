package main

import (
	"os"

	"github.com/roemer/gotaskr"
	"github.com/roemer/gotaskr/log"
)

func main() {
	os.Exit(gotaskr.Execute())
}

func init() {
	gotaskr.Task("Build", func() error {
		log.Information("Build...")
		return nil
	}).Then("Test")

	gotaskr.Task("Test", func() error {
		log.Information("Test...")
		return nil
	})

	gotaskr.Setup(func() error {
		log.Information("Setup...")
		return nil
	})

	gotaskr.Teardown(func() error {
		log.Information("Teardown...")
		return nil
	})

	gotaskr.TaskSetup(func() error {
		log.Information("TaskSetup...")
		return nil
	})

	gotaskr.TaskTeardown(func() error {
		log.Information("TaskTeardown...")
		return nil
	})
}
