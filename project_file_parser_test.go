package dotnetpublish_test

import (
	"io/ioutil"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	dotnetpublish "github.com/paketo-buildpacks/dotnet-publish"
	"github.com/sclevine/spec"
)

func testProjectFileParser(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect
		parser dotnetpublish.ProjectFileParser
	)

	it.Before(func() {
		parser = dotnetpublish.NewProjectFileParser()
	})

	context("ParseVersion", func() {
		var path string

		it.Before(func() {
			file, err := ioutil.TempFile("", "app.csproj")
			Expect(err).NotTo(HaveOccurred())
			defer file.Close()

			_, err = file.WriteString(`
				<Project>
					<PropertyGroup>
						<Obfuscate>true</Obfuscate>
					</PropertyGroup>
					<PropertyGroup>
						<RuntimeFrameworkVersion>1.2.3</RuntimeFrameworkVersion>
					</PropertyGroup>
				</Project>
			`)
			Expect(err).NotTo(HaveOccurred())

			path = file.Name()
		})

		it.After(func() {
			Expect(os.Remove(path)).To(Succeed())
		})

		it("parses the dotnet runtime version from the ?sproj file", func() {
			version, err := parser.ParseVersion(path)
			Expect(err).NotTo(HaveOccurred())
			Expect(version).To(Equal("1.2.3"))
		})

		context("when the RuntimeFrameworkVersion is not specified", func() {
			it.Before(func() {
				err := ioutil.WriteFile(path, []byte(`
					<Project>
						<PropertyGroup>
							<TargetFramework>netcoreapp1.2</TargetFramework>
						</PropertyGroup>
					</Project>
				`), 0600)
				Expect(err).NotTo(HaveOccurred())
			})

			it("falls back to using the TargetFramework version", func() {
				version, err := parser.ParseVersion(path)
				Expect(err).NotTo(HaveOccurred())
				Expect(version).To(Equal("1.2.0"))
			})
		})

		context("failure cases", func() {
			context("when the project file does not exist", func() {
				it("returns an error", func() {
					_, err := parser.ParseVersion("no-such-file")
					Expect(err).To(MatchError(MatchRegexp(`failed to read project file: .* no such file or directory`)))
				})
			})

			context("when the project file is malformed", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`<<< %%%`), 0600)
					Expect(err).NotTo(HaveOccurred())
				})

				it("returns an error", func() {
					_, err := parser.ParseVersion(path)
					Expect(err).To(MatchError(MatchRegexp(`failed to parse project file: XML syntax error .*`)))
				})
			})

			context("when the project file does not contain a version", func() {
				it.Before(func() {
					err := ioutil.WriteFile(path, []byte(`<Project></Project>`), 0600)
					Expect(err).NotTo(HaveOccurred())
				})

				it("returns an error", func() {
					_, err := parser.ParseVersion(path)
					Expect(err).To(MatchError("failed to find version in project file: missing TargetFramework property"))
				})
			})
		})
	})

	context("ASPNetIsRequired", func() {
		var path string

		it.Before(func() {
			file, err := ioutil.TempFile("", "app.csproj")
			Expect(err).NotTo(HaveOccurred())
			defer file.Close()

			_, err = file.WriteString(`<Project Sdk="Microsoft.NET.Sdk.Web"></Project>`)
			Expect(err).NotTo(HaveOccurred())

			path = file.Name()
		})

		it.After(func() {
			Expect(os.Remove(path)).To(Succeed())
		})

		context("when project SDK is Microsoft.NET.Sdk.Web", func() {
			it("reports if the application requires ASP.Net", func() {
				Expect(parser.ASPNetIsRequired(path)).To(BeTrue())
			})
		})
	})
}
