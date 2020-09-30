package dotnetpublish

import (
	"fmt"
	"os"

	"github.com/paketo-buildpacks/packit/pexec"
)

type Executable interface {
	Execute(pexec.Execution) error
}

type DotnetPublishProcess struct {
	executable Executable
}

func NewDotnetPublishProcess(executable Executable) DotnetPublishProcess {
	return DotnetPublishProcess{
		executable: executable,
	}
}

func (p DotnetPublishProcess) Execute(workingDir, root string) error {
	// TODO: remove when https://github.com/paketo-buildpacks/packit/pull/73 is merged and released
	os.Setenv("PATH", fmt.Sprintf("%s:%s", root, os.Getenv("PATH")))

	err := p.executable.Execute(pexec.Execution{
		Args: []string{
			"publish", workingDir,
			"--configuration", "Release",
			"--runtime", "ubuntu.18.04-x64",
			"--self-contained", "false",
			"--output", workingDir,
		},
		Dir: workingDir,

		// TODO: uncomment when https://github.com/paketo-buildpacks/packit/pull/73 is merged and released
		// Env: append(os.Environ(), fmt.Sprintf("PATH=%s:%s", root, os.Getenv("PATH"))),

		// TODO: remove, as these are only for debugging
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("failed to execute 'dotnet publish': %w", err)
	}

	return nil
}
