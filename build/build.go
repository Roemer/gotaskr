package main

import (
	"os"

	"github.com/roemer/gotaskr"
	"github.com/roemer/gotaskr/execr"
)

func main() {
	os.Exit(gotaskr.Execute())
}

////////////////////////////////////////////////////////////
// Task Registration
////////////////////////////////////////////////////////////

func init() {
	gotaskr.Task("UpdateDependencies", UpdateDependencies).Description("Updates the dependencies")
}

func UpdateDependencies() error {
	if err := execr.Run("go", []string{"get", "-u"}...); err != nil {
		return err
	}
	if err := execr.Run("go", []string{"run", "tidy"}...); err != nil {
		return err
	}
	return nil
}
