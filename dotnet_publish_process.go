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
	err := p.executable.Execute(pexec.Execution{
		Args: []string{
			"publish", workingDir,
			"--configuration", "Release",
			"--runtime", "ubuntu.18.04-x64",
			"--self-contained", "false",
			"--output", workingDir,
		},
		Dir: workingDir,
		Env: append(os.Environ(), fmt.Sprintf("DOTNET_ROOT=%s", root)),

		// TODO: remove, as these are only for debugging
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("failed to execute 'dotnet publish': %w", err)
	}

	return nil
}
