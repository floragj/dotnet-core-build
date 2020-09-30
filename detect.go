package dotnetpublish

import (
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

func Detect(parser ProjectParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		matches, err := filepath.Glob(filepath.Join(context.WorkingDir, "*.?sproj"))
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
			panic("invalid version")
		}
		sdkVersion := strings.Join([]string{parts[0], parts[1], "0"}, ".")

		requirements := []packit.BuildPlanRequirement{
			{
				Name: "build",
				Metadata: BuildPlanMetadata{
					Build: true,
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
