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
	if err := execr.RunO("go", execr.WithArgsSplitted("get -u"), execr.WithConsoleOutput(true)); err != nil {
		return err
	}
	if err := execr.RunO("go", execr.WithArgsSplitted("mod tidy"), execr.WithConsoleOutput(true)); err != nil {
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
	if err := execr.RunO("go", execr.WithArgsSplitted("install github.com/jstemmer/go-junit-report/v2@latest"), execr.WithConsoleOutput(true)); err != nil {
		return err
	}

	stdout, _, err := execr.RunOGetOutput("go", execr.WithArgsSplitted("test -v ./... "), execr.WithConsoleOutput(true))
	if err != nil {
		return nil
	}
	if err := os.WriteFile(goTestReport, []byte(stdout), os.ModePerm); err != nil {
		return err
	}

	if err := execr.RunO(path.Join(build.Default.GOPATH, "bin/go-junit-report"), execr.WithArgs("-in", goTestReport, "-set-exit-code", "-out", junitTestReport), execr.WithConsoleOutput(true)); err != nil {
		return err
	}
	return nil
}
