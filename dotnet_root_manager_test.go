package dotnetpublish_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	dotnetpublish "github.com/paketo-buildpacks/dotnet-publish"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDotnetRootManager(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		existingRootDir string
		sdkLayerDir     string
		manager         dotnetpublish.DotnetRootManager
	)

	it.Before(func() {
		var err error
		existingRootDir, err = ioutil.TempDir("", "existing-root-dir")
		Expect(err).NotTo(HaveOccurred())

		sdkLayerDir, err = ioutil.TempDir("", "existing-root-dir")
		Expect(err).NotTo(HaveOccurred())

		Expect(os.MkdirAll(filepath.Join(existingRootDir, "host"), os.ModePerm)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(existingRootDir, "shared", "some-dir"), os.ModePerm)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(existingRootDir, "shared", "some-file"), nil, 0600)).To(Succeed())
		Expect(ioutil.WriteFile(filepath.Join(existingRootDir, "dotnet"), nil, 0700)).To(Succeed())
		Expect(os.MkdirAll(filepath.Join(sdkLayerDir, "sdk"), os.ModePerm)).To(Succeed())

		manager = dotnetpublish.NewDotnetRootManager()
	})

	it.After(func() {
		Expect(os.RemoveAll(existingRootDir)).To(Succeed())
		Expect(os.RemoveAll(sdkLayerDir)).To(Succeed())
	})

	context("Setup", func() {
		it("sets up the DOTNET_ROOT directory", func() {
			root, err := manager.Setup(existingRootDir, sdkLayerDir)
			Expect(err).NotTo(HaveOccurred())

			var files []string
			err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				files = append(files, path)

				return nil
			})
			Expect(err).NotTo(HaveOccurred())

			Expect(files).To(ConsistOf([]string{
				filepath.Join(root),
				filepath.Join(root, "host"),
				filepath.Join(root, "dotnet"),
				filepath.Join(root, "shared"),
				filepath.Join(root, "shared", "some-dir"),
				filepath.Join(root, "shared", "some-file"),
				filepath.Join(root, "sdk"),
			}))

			Expect(os.Getenv("PATH")).To(ContainSubstring(root))
		})
	})
}
