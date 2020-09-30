package dotnetpublish_test

import (
	"errors"
	"fmt"
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

		projectParser      *fakes.ProjectParser
		buildpackYMLParser *fakes.YMLParser
		workingDir         string
		detect             packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		Expect(ioutil.WriteFile(filepath.Join(workingDir, "app.xsproj"), nil, 0600)).To(Succeed())

		projectParser = &fakes.ProjectParser{}
		buildpackYMLParser = &fakes.YMLParser{}
		projectParser.ParseVersionCall.Returns.Version = "1.2.3"

		detect = dotnetpublish.Detect(projectParser, buildpackYMLParser)
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
			projectParser.ASPNetIsRequiredCall.Returns.Bool = true
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

	context("when node is required", func() {
		it.Before(func() {
			projectParser.NodeIsRequiredCall.Returns.Bool = true
		})

		it("requires node in the build plan", func() {
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
							Name: "node",
							Metadata: dotnetpublish.BuildPlanMetadata{
								Build:  true,
								Launch: true,
							},
						},
					},
				},
			}))
		})
	})

	context("when npm is required", func() {
		it.Before(func() {
			projectParser.NodeIsRequiredCall.Returns.Bool = true
			projectParser.NPMIsRequiredCall.Returns.Bool = true
		})

		it("requires node in the build plan", func() {
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
							Name: "node",
							Metadata: dotnetpublish.BuildPlanMetadata{
								Build:  true,
								Launch: true,
							},
						},
						{
							Name: "npm",
							Metadata: dotnetpublish.BuildPlanMetadata{
								Build:  true,
								Launch: true,
							},
						},
					},
				},
			}))
		})
	})

	context("when the .csproj file is not at the base of the directory and project_path is set in buildpack.yml", func() {

		it.Before(func() {
			buildpackYMLParser.ParseProjectPathCall.Returns.ProjectFilePath = "src/proj1"
			Expect(os.RemoveAll(workingDir)).To(Succeed())
			Expect(os.MkdirAll(filepath.Join(workingDir, "src", "proj1"), os.ModePerm)).To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(workingDir, "src", "proj1", "app.xsproj"), nil, 0600)).To(Succeed())
			Expect(ioutil.WriteFile(filepath.Join(workingDir, "buildpack.yml"), nil, 0600)).To(Succeed())
		})

		it.After(func() {
			Expect(os.RemoveAll(workingDir)).To(Succeed())
		})

		it("finds the projfile and passes detection", func() {

			_, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})
	context("failure cases", func() {
		context("when buildpack.yml cannot be stated", func() {
			it.Before(func() {
				Expect(os.Chmod(workingDir, 0000)).To(Succeed())
			})

			it.After(func() {
				Expect(os.Chmod(workingDir, os.ModePerm)).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("permission denied")))
			})
		})

		context("when buildpack.yml cannot be parsed", func() {
			it.Before(func() {
				buildpackYMLParser.ParseProjectPathCall.Returns.Err = fmt.Errorf("parsing error")
				Expect(ioutil.WriteFile(filepath.Join(workingDir, "buildpack.yml"), nil, 0600)).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to parse buildpack.yml:")))
			})
		})

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
				projectParser.ParseVersionCall.Returns.Err = errors.New("failed to parse project file version")
			})

			it("returns an error", func() {
				_, err := detect(packit.DetectContext{WorkingDir: workingDir})
				Expect(err).To(MatchError("failed to parse project file version"))
			})
		})
	})
}
