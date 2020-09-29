package dotnetpublish

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strings"
)

type ProjectFileParser struct{}

func NewProjectFileParser() ProjectFileParser {
	return ProjectFileParser{}
}

func (p ProjectFileParser) ParseVersion(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to read project file: %w", err)
	}
	defer file.Close()

	var project struct {
		PropertyGroups []struct {
			RuntimeFrameworkVersion string
			TargetFramework         string
		} `xml:"PropertyGroup"`
	}

	err = xml.NewDecoder(file).Decode(&project)
	if err != nil {
		return "", fmt.Errorf("failed to parse project file: %w", err)
	}

	for _, group := range project.PropertyGroups {
		if group.RuntimeFrameworkVersion != "" {
			return group.RuntimeFrameworkVersion, nil
		}
	}

	for _, group := range project.PropertyGroups {
		if strings.HasPrefix(group.TargetFramework, "netcoreapp") {
			return fmt.Sprintf("%s.0", strings.TrimPrefix(group.TargetFramework, "netcoreapp")), nil
		}
	}

	return "", errors.New("failed to find version in project file: missing TargetFramework property")
}

func (p ProjectFileParser) ASPNetIsRequired(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var project struct {
		SDK string `xml:"Sdk,attr"`
	}

	err = xml.NewDecoder(file).Decode(&project)
	if err != nil {
		panic(err)
	}

	return project.SDK == "Microsoft.NET.Sdk.Web"
}
