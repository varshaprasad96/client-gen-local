package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/varshaprasad96/client-gen/pkg/custom"
	flag "github.com/varshaprasad96/client-gen/pkg/custom/flag"
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
