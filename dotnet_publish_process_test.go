package dotnetpublish_test

import (
	"errors"
	"testing"

	dotnetpublish "github.com/paketo-buildpacks/dotnet-publish"
	"github.com/paketo-buildpacks/packit/cargo/fakes"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDotnetPublishProcess(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		executable *fakes.Executable
		process    dotnetpublish.DotnetPublishProcess
	)

	it.Before(func() {
		executable = &fakes.Executable{}

		process = dotnetpublish.NewDotnetPublishProcess(executable)
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
		Expect(executable.ExecuteCall.Receives.Execution.Env).To(ContainElement("DOTNET_ROOT=some-dotnet-root-dir"))
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
