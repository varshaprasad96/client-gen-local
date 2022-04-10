package custom

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/varshaprasad96/client-gen/pkg/custom/flag"
	"k8s.io/code-generator/cmd/client-gen/args"
	"k8s.io/code-generator/cmd/client-gen/types"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type Generator struct {
	// BaseImportPath refers to the base path of the package
	inputDir string
	// Output Dir
	outputDir string
	// ClienSetAPI path
	clientSetAPIPath string
}

func (g *Generator) setdefualtsFromFlags(f flag.Flags) error {
	if f.InputDir == "" {
		return fmt.Errorf("currently generator does not run without input path to API definition")
	}
	g.inputDir = f.InputDir

	if f.OutputDir != "" {
		g.outputDir = f.OutputDir
	}

	if f.ClientsetAPIPath == "" {
		return fmt.Errorf("specifying client API path is required currently")
	}

	g.clientSetAPIPath = f.ClientsetAPIPath
	return nil
}

func (g Generator) Generate(ctx *genall.GenerationContext, f flag.Flags) error {
	if err := g.setdefualtsFromFlags(f); err != nil {
		return err
	}

	err := g.generateHelper(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (g Generator) generateHelper(ctx *genall.GenerationContext) error {

	for _, root := range ctx.Roots {
		root.NeedTypesInfo()

		byType := make(map[string][]byte)

		group, err := getGroups(root, g.inputDir)
		if err != nil {
			fmt.Println(err.Error())
		}

		outCommonContent := new(bytes.Buffer)
		pkgmg, err := NewPackage(root, root.Package.PkgPath, g.clientSetAPIPath, string(group.Versions[0].Version), group.PackageName, &codeWriter{out: outCommonContent})
		if err != nil {
			return err
		}

		err = pkgmg.writeCommonContent()
		if err != nil {
			return err
		}

		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			outContent := new(bytes.Buffer)

			// if not enabled for this type, skip
			if !isEnabledForMethod(info) {
				return
			}

			p, err := NewAPI(root, info, string(group.Versions[0].Version), group.PackageName, &codeWriter{out: outContent})
			if err != nil {
				fmt.Println(err.Error())
			}

			err = p.writeMethods()
			if err != nil {
				fmt.Println(err.Error())
			}

			outBytes := outContent.Bytes()
			if len(outBytes) > 0 {
				byType[info.Name] = outBytes
			}
		}); err != nil {
			return err
		}

		if len(byType) == 0 {
			return nil
		}

		outContent := new(bytes.Buffer)
		writeHeader(root, outContent, root.Name)
		outContent.Write(outCommonContent.Bytes())
		writeMethods(root, outContent, byType)

		outBytes := outContent.Bytes()
		formattedBytes, err := format.Source(outBytes)
		if err != nil {
			root.AddError(err)
			// we still write the invalid source to disk to figure out what went wrong
		} else {
			outBytes = formattedBytes
		}

		err = writeOut(ctx, root, outBytes)
		if err != nil {
			return err
		}
	}
	return nil
}

func getGroups(pkg *loader.Package, basePath string) (*types.GroupVersions, error) {
	groups := []types.GroupVersions{}

	// Using the builder from code-gen
	builder := args.NewGroupVersionsBuilder(&groups)
	_ = args.NewGVPackagesValue(builder, []string{pkg.PkgPath})

	if len(groups) == 0 {
		return nil, fmt.Errorf("error finding the group version from import path %q", basePath)
	}

	if len(groups) > 1 {
		return nil, fmt.Errorf("specifying multiple groups in the same package %q is not supported", basePath)
	}

	return &groups[0], nil
}

func writeMethods(pkg *loader.Package, out io.Writer, byType map[string][]byte) {
	soretedNames := make([]string, 0, len(byType))
	for name := range byType {
		soretedNames = append(soretedNames, name)
	}
	sort.Strings(soretedNames)

	for _, name := range soretedNames {
		_, err := out.Write(byType[name])
		if err != nil {
			// expose this error
			pkg.AddError(err)
		}
	}
}

func (g Generator) RegisterMarkers(into *markers.Registry) error {
	if err := into.Register(RuleDefinition); err != nil {
		return err
	}
	// Skipping adding Help for this marker for now
	return nil
}

// Wire in output rules instead of creating a file in here. Use pkg/genall/output.go
func writeOut(ctx *genall.GenerationContext, root *loader.Package, outbytes []byte) error {

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	path := filepath.Join(wd, "zz_generated_test.go")

	outputFile, err := os.Create(path)
	if err != nil {
		return err
	}

	defer outputFile.Close()

	n, err := outputFile.Write(outbytes)
	if err != nil {
		return err
	}

	if n < len(outbytes) {
		return err
	}

	return nil
}

func writeHeader(pkg *loader.Package, out io.Writer, packageName string) {
	_, err := fmt.Fprintf(out, `// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package %[1]s

`, packageName)
	if err != nil {
		pkg.AddError(err)
	}
}

// isEnabledForMethod verifies if the genclient marker is enabled for
// this type or not
func isEnabledForMethod(info *markers.TypeInfo) bool {
	enabled := info.Markers.Get(RuleDefinition.Name)
	return enabled != nil
}
