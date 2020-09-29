package dotnetpublish

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/paketo-buildpacks/packit/fs"
)

type DotnetRootManager struct{}

func NewDotnetRootManager() DotnetRootManager {
	return DotnetRootManager{}
}

func (m DotnetRootManager) Setup(existingRoot, sdkLocation string) (string, error) {
	root, err := ioutil.TempDir("", "dotnet-root")
	if err != nil {
		panic(err)
	}

	paths, err := filepath.Glob(filepath.Join(existingRoot, "shared", "*"))
	if err != nil {
		panic(err)
	}

	for _, path := range paths {
		relPath, err := filepath.Rel(existingRoot, path)
		if err != nil {
			panic(err)
		}

		newPath := filepath.Join(root, relPath)
		err = os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
		if err != nil {
			panic(err)
		}

		err = os.Symlink(path, newPath)
		if err != nil {
			panic(err)
		}
	}

	err = os.Symlink(filepath.Join(existingRoot, "host"), filepath.Join(root, "host"))
	if err != nil {
		panic(err)
	}

	// NOTE: the dotnet CLI uses relative pathing that means we must copy it into
	// the final DOTNET_ROOT so that it can find SDKs.
	err = fs.Copy(filepath.Join(existingRoot, "dotnet"), filepath.Join(root, "dotnet"))
	if err != nil {
		panic(err)
	}

	err = os.Symlink(filepath.Join(sdkLocation, "sdk"), filepath.Join(root, "sdk"))
	if err != nil {
		panic(err)
	}

	return root, nil
}
