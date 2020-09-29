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
	// TODO: figure out why this isn't working
	// var env []string
	// for _, variable := range os.Environ() {
	// 	if !strings.HasPrefix(variable, "PATH=") {
	// 		env = append(env, variable)
	// 	}
	// }
	// env = append(env, fmt.Sprintf("PATH=%s:%s", root, os.Getenv("PATH")))

	os.Setenv("PATH", fmt.Sprintf("%s:%s", root, os.Getenv("PATH")))

	err := p.executable.Execute(pexec.Execution{
		Args:   []string{"--info"},
		Dir:    workingDir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		panic(err)
	}

	err = p.executable.Execute(pexec.Execution{
		Args: []string{
			"publish", workingDir,
			"--configuration", "Release",
			"--runtime", "ubuntu.18.04-x64",
			"--self-contained", "false",
			"--output", workingDir,
		},
		Dir: workingDir,

		// TODO: remove, as these are only for debugging
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("failed to execute 'dotnet publish': %w", err)
	}

	return nil
}
