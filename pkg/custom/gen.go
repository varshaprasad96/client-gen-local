package custom

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

type Generator struct {
	// Name - placeholder for now
	Name string
}

type result struct {
	Group string
	Type  string
}

func (g Generator) Generate(ctx *genall.GenerationContext) error {
	res, err := GenerateHelper(ctx)
	if err != nil {
		return err
	}

	var output string
	for _, r := range res {
		output = output + fmt.Sprintf("Extracted values types: %s Groups: %s \n", r.Type, r.Group)
	}

	for _, root := range ctx.Roots {
		err := writeOut(ctx, root, []byte(output))
		if err != nil {
			return err
		}
	}
	return nil
}

func GenerateHelper(ctx *genall.GenerationContext) ([]result, error) {
	var out []result
	for _, root := range ctx.Roots {
		markerSet, err := markers.PackageMarkers(ctx.Collector, root)
		if err != nil {
			return nil, err
		}

		for _, markerValues := range markerSet[RuleDefinition.Name] {
			rule := markerValues.(Rule)
			r := result{
				Group: rule.Groups[0],
				Type:  rule.Types[0],
			}
			out = append(out, r)
		}
	}

	return out, nil
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

	path := filepath.Join(wd, "zz_generated_test.txt")

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
