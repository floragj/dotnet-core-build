package dotnetpublish

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type BuildpackYMLParser struct{}

func NewBuildpackYMLParser() BuildpackYMLParser {
	return BuildpackYMLParser{}
}

func (p BuildpackYMLParser) ParseProjectPath(path string) (string, error) {
	var buildpack struct {
		Config struct {
			ProjectPath string `yaml:"project-path"`
		} `yaml:"dotnet-build"`
	}

	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	err = yaml.NewDecoder(file).Decode(&buildpack)

	if err != nil {
		return "", fmt.Errorf("invalid buildpack.yml: %w", err)
	}

	return buildpack.Config.ProjectPath, nil
}
