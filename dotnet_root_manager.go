package dotnetpublish

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type DotnetRootManager struct{}

func NewDotnetRootManager() DotnetRootManager {
	return DotnetRootManager{}
}

func (m DotnetRootManager) Setup(existingRoot, sdkLocation string) (string, error) {
	fmt.Println("Existing $DOTNET_ROOT")
	err := filepath.Walk(existingRoot, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return err
	})
	if err != nil {
		panic(err)
	}

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

	err = os.Symlink(filepath.Join(existingRoot, "dotnet"), filepath.Join(root, "dotnet"))
	if err != nil {
		panic(err)
	}

	err = os.Symlink(sdkLocation, filepath.Join(root, "sdk"))
	if err != nil {
		panic(err)
	}

	fmt.Println("New $DOTNET_ROOT")
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		return err
	})
	if err != nil {
		panic(err)
	}

	return root, nil
}
