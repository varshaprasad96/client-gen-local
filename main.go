package main

import (
	"fmt"
	"log"
	"os"

	"github.com/varshaprasad96/client-gen/pkg/custom"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func main() {

	err := os.Chdir("./testdata")
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	pkgs, err := loader.LoadRoots(".")
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	// TODO: call custom.RegisterInto instead
	reg := &markers.Registry{}
	err = reg.Register(custom.RuleDefinition)
	if err != nil {
		log.Fatalf(err.Error())
		os.Exit(1)
	}

	ctx := &genall.GenerationContext{
		Collector: &markers.Collector{Registry: reg},
		Roots:     pkgs,
	}

	g := custom.Generator{"test"}
	err = g.Generate(ctx)
	if err != nil {
		fmt.Println(err)
	}

}
