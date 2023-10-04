package main

import (
	"go/build"
	"os"
	"path"

	"github.com/roemer/gotaskr"
	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gotaskr/log"
)

func main() {
	os.Exit(gotaskr.Execute())
}

////////////////////////////////////////////////////////////
// Variables
////////////////////////////////////////////////////////////

var reportPath = "reports"

////////////////////////////////////////////////////////////
// Task Registration
////////////////////////////////////////////////////////////

func init() {
	gotaskr.Task("Maintenance:Update-Dependencies", UpdateDependencies).
		Description("Updates the dependencies" + log.Newline + "Updates and tidies")

	gotaskr.Task("Run-Tests", RunTests)
}

////////////////////////////////////////////////////////////
// Tasks
////////////////////////////////////////////////////////////

func UpdateDependencies() error {
	if err := execr.Run(true, "go", "get", "-u"); err != nil {
		return err
	}
	if err := execr.Run(true, "go", "mod", "tidy"); err != nil {
		return err
	}
	return nil
}

func RunTests() error {
	if err := os.MkdirAll(reportPath, os.ModePerm); err != nil {
		return err
	}
	goTestReport := path.Join(reportPath, "go-test-report.txt")
	junitTestReport := path.Join(reportPath, "junit-test-report.xml")
	if err := execr.Run(true, "go", "install", "github.com/jstemmer/go-junit-report/v2@latest"); err != nil {
		return err
	}

	stdout, _, err := execr.RunGetOutput(false, "go", execr.SplitArgumentString("test -v ./... ")...)
	if err != nil {
		return nil
	}
	if err := os.WriteFile(goTestReport, []byte(stdout), os.ModePerm); err != nil {
		return err
	}

	if err := execr.Run(true, path.Join(build.Default.GOPATH, "bin/go-junit-report"), "-in", goTestReport, "-set-exit-code", "-out", junitTestReport); err != nil {
		return err
	}
	return nil
}
