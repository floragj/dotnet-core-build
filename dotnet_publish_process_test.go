package dotnetpublish_test

import (
	"errors"
	"os"
	"testing"

	dotnetpublish "github.com/paketo-buildpacks/dotnet-publish"
	"github.com/paketo-buildpacks/packit/cargo/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDotnetPublishProcess(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		path       string
		executable *fakes.Executable
		process    dotnetpublish.DotnetPublishProcess
	)

	it.Before(func() {
		path = os.Getenv("PATH")
		Expect(os.Setenv("PATH", "some-path")).To(Succeed())

		executable = &fakes.Executable{}

		process = dotnetpublish.NewDotnetPublishProcess(executable)
	})

	it.After(func() {
		Expect(os.Setenv("PATH", path)).To(Succeed())
	})

	it("executes the dotnet publish process", func() {
		err := process.Execute("some-working-dir", "some-dotnet-root-dir")
		Expect(err).NotTo(HaveOccurred())

		Expect(executable.ExecuteCall.Receives.Execution.Args).To(Equal([]string{
			"publish", "some-working-dir",
			"--configuration", "Release",
			"--runtime", "ubuntu.18.04-x64",
			"--self-contained", "false",
			"--output", "some-working-dir",
		}))

		Expect(executable.ExecuteCall.Receives.Execution.Dir).To(Equal("some-working-dir"))

		// TODO: uncomment when https://github.com/paketo-buildpacks/packit/pull/73 is merged and released
		// Expect(executable.ExecuteCall.Receives.Execution.Env).To(ContainElement("PATH=some-dotnet-root-dir:some-path"))
	})

	context("failure cases", func() {
		context("when the dotnet publish executable errors", func() {
			it.Before(func() {
				executable.ExecuteCall.Returns.Error = errors.New("execution error")
			})

			it("returns an error", func() {
				err := process.Execute("some-working-dir", "some-dotnet-root-dir")
				Expect(err).To(MatchError("failed to execute 'dotnet publish': execution error"))
			})
		})
	})
}
