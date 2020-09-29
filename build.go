package dotnetpublish

import (
	"os"

	"github.com/paketo-buildpacks/packit"
)

//go:generate faux --interface RootManager --output fakes/root_manager.go
type RootManager interface {
	Setup(existingRoot, sdkLocation string) (root string, err error)
}

//go:generate faux --interface PublishProcess --output fakes/publish_process.go
type PublishProcess interface {
	Execute(workingDir, rootDir string) error
}

func Build(rootManager RootManager, publishProcess PublishProcess) packit.BuildFunc {
	return func(context packit.BuildContext) (packit.BuildResult, error) {
		root, err := rootManager.Setup(os.Getenv("DOTNET_ROOT"), os.Getenv("SDK_LOCATION"))
		if err != nil {
			panic(err)
		}

		err = publishProcess.Execute(context.WorkingDir, root)
		if err != nil {
			panic(err)
		}

		return packit.BuildResult{}, nil
	}
}
