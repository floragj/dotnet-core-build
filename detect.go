package dotnetpublish

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit"
)

type BuildPlanMetadata struct {
	Version string `toml:"version,omitempty"`
	Build   bool   `toml:"build"`
	Launch  bool   `toml:"launch"`
}

//go:generate faux --interface ProjectParser --output fakes/project_parser.go
type ProjectParser interface {
	ParseVersion(path string) (version string, err error)

	ASPNetIsRequired(path string) bool
	NodeIsRequired(path string) bool
	NPMIsRequired(path string) bool
}

//go:generate faux --interface YMLParser --output fakes/yml_parser.go
type YMLParser interface {
	ParseProjectPath(path string) (projectFilePath string, err error)
}

func Detect(parser ProjectParser, ymlParser YMLParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		var projectPath string
		_, err := os.Stat(filepath.Join(context.WorkingDir, "buildpack.yml"))
		if err != nil {
			if !os.IsNotExist(err) {
				return packit.DetectResult{}, err
			}
		} else {
			projectPath, err = ymlParser.ParseProjectPath(filepath.Join(context.WorkingDir, "buildpack.yml"))
			if err != nil {
				return packit.DetectResult{}, fmt.Errorf("failed to parse buildpack.yml: %w", err)
			}
		}

		matches, err := filepath.Glob(filepath.Join(context.WorkingDir, projectPath, "*.?sproj"))
		if err != nil {
			return packit.DetectResult{}, err
		}

		if len(matches) == 0 {
			return packit.DetectResult{}, packit.Fail.WithMessage("no project file found")
		}

		projectFilePath := matches[0]
		runtimeVersion, err := parser.ParseVersion(projectFilePath)
		if err != nil {
			return packit.DetectResult{}, err
		}

		parts := strings.Split(runtimeVersion, ".")
		if len(parts) < 2 {
			panic("invalid version") // this replicates original buildpack behaviour
		}
		sdkVersion := strings.Join([]string{parts[0], parts[1], "0"}, ".")

		requirements := []packit.BuildPlanRequirement{
			{
				Name: "build",
				Metadata: BuildPlanMetadata{
					Build: true,
					// TODO: pass the project path here so we can grab it for build
				},
			},
			{
				Name: "dotnet-sdk",
				Metadata: BuildPlanMetadata{
					Version: sdkVersion,
					Build:   true,
					Launch:  true,
				},
			},
			{
				Name: "dotnet-runtime",
				Metadata: BuildPlanMetadata{
					Version: runtimeVersion,
					Build:   true,
					Launch:  true,
				},
			},
		}

		if parser.ASPNetIsRequired(projectFilePath) {
			requirements = append(requirements, packit.BuildPlanRequirement{
				Name: "dotnet-aspnetcore",
				Metadata: BuildPlanMetadata{
					Version: runtimeVersion,
					Build:   true,
					Launch:  true,
				},
			})
		}

		if parser.NodeIsRequired(projectFilePath) {
			requirements = append(requirements, packit.BuildPlanRequirement{
				Name: "node",
				Metadata: BuildPlanMetadata{
					Build:  true,
					Launch: true,
				},
			})
		}

		if parser.NPMIsRequired(projectFilePath) {
			requirements = append(requirements, packit.BuildPlanRequirement{
				Name: "npm",
				Metadata: BuildPlanMetadata{
					Build:  true,
					Launch: true,
				},
			})
		}

		return packit.DetectResult{
			Plan: packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{
					{Name: "build"},
				},
				Requires: requirements,
			},
		}, nil
	}
}
