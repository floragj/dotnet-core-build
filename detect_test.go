package dotnetpublish_test

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	dotnetpublish "github.com/paketo-buildpacks/dotnet-publish"
	"github.com/paketo-buildpacks/dotnet-publish/fakes"
	"github.com/paketo-buildpacks/packit"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		versionParser *fakes.VersionParser
		workingDir    string
		detect        packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.WriteFile(filepath.Join(workingDir, "app.xsproj"), nil, 0600)).To(Succeed())

		versionParser = &fakes.VersionParser{}
		versionParser.ParseVersionCall.Returns.Version = "1.2.3"
		detect = dotnetpublish.Detect(versionParser)
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	it("returns a build plan", func() {
		result, err := detect(packit.DetectContext{
			WorkingDir: workingDir,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(Equal(packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "build"},
				},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "build",
						Metadata: dotnetpublish.BuildPlanMetadata{
							Build: true,
						},
					},
					{
						Name: "dotnet-sdk",
						Metadata: dotnetpublish.BuildPlanMetadata{
							Version: "1.2.0",
							Build:   true,
							Launch:  true,
						},
					},
					{
						Name: "dotnet-runtime",
						Metadata: dotnetpublish.BuildPlanMetadata{
							Version: "1.2.3",
							Build:   true,
							Launch:  true,
						},
					},
				},
			},
		}))
	})

	context("when aspnet is required", func() {
		it.Before(func() {
			versionParser.ASPNetIsRequiredCall.Returns.Bool = true
		})

		it("requires aspnet in the build plan", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(packit.DetectResult{
				Plan: packit.BuildPlan{
					Provides: []packit.BuildPlanProvision{
						{Name: "build"},
					},
					Requires: []packit.BuildPlanRequirement{
						{
							Name: "build",
							Metadata: dotnetpublish.BuildPlanMetadata{
								Build: true,
							},
						},
						{
							Name: "dotnet-sdk",
							Metadata: dotnetpublish.BuildPlanMetadata{
								Version: "1.2.0",
								Build:   true,
								Launch:  true,
							},
						},
						{
							Name: "dotnet-runtime",
							Metadata: dotnetpublish.BuildPlanMetadata{
								Version: "1.2.3",
								Build:   true,
								Launch:  true,
							},
						},
						{
							Name: "dotnet-aspnetcore",
							Metadata: dotnetpublish.BuildPlanMetadata{
								Version: "1.2.3",
								Build:   true,
								Launch:  true,
							},
						},
					},
				},
			}))
		})
	})

	context("failure cases", func() {
		context("when a project file cannot be found", func() {
			it.Before(func() {
				Expect(os.Remove(filepath.Join(workingDir, "app.xsproj"))).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(packit.Fail.WithMessage("no project file found")))
			})
		})

		context("when .?proj file cannot be parsed", func() {
			it.Before(func() {
				versionParser.ParseVersionCall.Returns.Err = errors.New("failed to parse project file version")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{WorkingDir: workingDir})
				Expect(err).To(MatchError("failed to parse project file version"))
			})
		})
	})
}
