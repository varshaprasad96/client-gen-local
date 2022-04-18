package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"
	"github.com/varshaprasad96/client-gen/pkg/custom"
	flag "github.com/varshaprasad96/client-gen/pkg/custom/flag"
	"k8s.io/code-generator/cmd/client-gen/args"
	"k8s.io/code-generator/cmd/client-gen/types"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func main() {

	f := flag.Flags{}
	f.AddTo(pflag.CommandLine)
	pflag.Parse()

	reg := &markers.Registry{}
	err := reg.Register(custom.RuleDefinition)
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	ctx := &genall.GenerationContext{
		Collector: &markers.Collector{Registry: reg},
	}

	g := custom.Generator{}
	err = g.Generate(ctx, f)
	if err != nil {
		fmt.Println(err)
	}

}

func main_test(input string, f flag.Flags) error {

	i := *f.GroupVersions
	arr := strings.Split(i[0], ":")
	fmt.Println("***", arr[0], arr[1])
	input = filepath.Join(input, arr[0], arr[1])
	fmt.Println(input)
	groups := []types.GroupVersions{}

	builder := args.NewGroupVersionsBuilder(&groups)
	_ = args.NewGVPackagesValue(builder, []string{input})

	fmt.Println(len(groups))
	fmt.Println(groups[0].Group, groups[0].Versions, groups[0].PackageName)
	return nil
}
