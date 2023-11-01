package main

import (
	"os"

	"github.com/roemer/gotaskr"
	"github.com/roemer/gotaskr/gttools"
	"github.com/roemer/gotaskr/log"
)

func main() {
	os.Exit(gotaskr.Execute())
}

func init() {
	initNxTasks()
}

func initNxTasks() {
	gotaskr.Task("Nx:List-Projects", func() error {
		nxSettings := gttools.NxShowProjectsSettings{}
		projects, err := gotaskr.Tools.Nx.ShowProjects(gttools.NxRunTypeNpx, nxSettings)
		if err != nil {
			return err
		}
		log.Information("Found the following projects:")
		for _, project := range projects {
			log.Information(project)
		}
		return nil
	}).Description("Lists all nx projects")
}
