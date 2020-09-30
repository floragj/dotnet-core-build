package main

import (
	dotnetpublish "github.com/paketo-buildpacks/dotnet-publish"
	"github.com/paketo-buildpacks/packit"
	"github.com/paketo-buildpacks/packit/pexec"
)

func main() {
	packit.Run(
		dotnetpublish.Detect(
			dotnetpublish.NewProjectFileParser(),
			dotnetpublish.NewBuildpackYMLParser(),
		),
		dotnetpublish.Build(
			dotnetpublish.NewDotnetRootManager(),
			dotnetpublish.NewDotnetPublishProcess(
				pexec.NewExecutable("dotnet"),
			),
		),
	)
}
