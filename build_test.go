package dotnetpublish_test

import (
	"io/ioutil"
	"os"
	"testing"

	dotnetpublish "github.com/paketo-buildpacks/dotnet-publish"
	"github.com/paketo-buildpacks/dotnet-publish/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBuild(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir     string
		rootManager    *fakes.RootManager
		publishProcess *fakes.PublishProcess
		build          packit.BuildFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		rootManager = &fakes.RootManager{}
		rootManager.SetupCall.Returns.Root = "some-root-dir"

		publishProcess = &fakes.PublishProcess{}

		os.Setenv("DOTNET_ROOT", "some-existing-root-dir")
		os.Setenv("SDK_LOCATION", "some-sdk-location")

		build = dotnetpublish.Build(rootManager, publishProcess)
	})

	it.After(func() {
		os.Unsetenv("DOTNET_ROOT")
		os.Unsetenv("SDK_LOCATION")

		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("returns a build result", func() {
		result, err := build(packit.BuildContext{
			WorkingDir: workingDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(packit.BuildResult{}))

		Expect(rootManager.SetupCall.Receives.ExistingRoot).To(Equal("some-existing-root-dir"))
		Expect(rootManager.SetupCall.Receives.SdkLocation).To(Equal("some-sdk-location"))

		Expect(publishProcess.ExecuteCall.Receives.WorkingDir).To(Equal(workingDir))
		Expect(publishProcess.ExecuteCall.Receives.RootDir).To(Equal("some-root-dir"))
	})
}
