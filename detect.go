package dotnetpublish

import (
	"path/filepath"

	"github.com/paketo-buildpacks/packit"
)

type BuildPlanMetadata struct {
	Version string `toml:"version"`
	Build   bool   `toml:"build"`
	Launch  bool   `toml:"launch"`
}

//go:generate faux --interface VersionParser --output fakes/version_parser.go
type VersionParser interface {
	ParseVersion(path string) (version string, err error)
	ASPNetIsRequired(path string) bool
}

func Detect(parser VersionParser) packit.DetectFunc {
	return func(context packit.DetectContext) (packit.DetectResult, error) {
		matches, err := filepath.Glob(filepath.Join(context.WorkingDir, "*.?sproj"))
		if err != nil {
			return packit.DetectResult{}, err
		}

		if len(matches) == 0 {
			return packit.DetectResult{}, packit.Fail.WithMessage("no project file found")
		}

		projectFilePath := matches[0]

		dotnetRuntimeVersion, err := parser.ParseVersion(projectFilePath)
		if err != nil {
			return packit.DetectResult{}, err
		}

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
					Build:  true,
					Launch: true,
				},
			},
			{
				Name: "dotnet-runtime",
				Metadata: BuildPlanMetadata{
					Version: dotnetRuntimeVersion,
					Build:   true,
					Launch:  true,
				},
			},
		}

		if parser.ASPNetIsRequired(projectFilePath) {
			requirements = append(requirements, packit.BuildPlanRequirement{
				Name: "dotnet-aspnetcore",
				Metadata: BuildPlanMetadata{
					Version: dotnetRuntimeVersion,
					Build:   true,
					Launch:  true,
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
