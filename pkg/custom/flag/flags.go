package flag

import "github.com/spf13/pflag"

// Flags - Options accepted by generator
type Flags struct {
	OutputDir        string
	InputDir         string
	ClientsetAPIPath string
	GroupVersions    *[]string
	InterfaceName    string
}

func (f *Flags) AddTo(flagset *pflag.FlagSet) {
	// TODO: FIgure out if its worth defaulting it to pkg/api/...
	flagset.StringVar(&f.InputDir, "input-dir", "", "Input directory where types are defined")
	flagset.StringVar(&f.OutputDir, "output-dir", "output", "Output directory where wrapped clients will be generatoed")
	flagset.StringVar(&f.ClientsetAPIPath, "clientset-api-path", "/apis", "package path where clients are generated")

	// TODO: Probably default this to be the package name
	flagset.StringVar(&f.InterfaceName, "interface", "", "name of the interface which needs to be wrapped")
	gv := flagset.StringSlice("group-versions", []string{}, "specify group versions for the clients")
	f.GroupVersions = gv
}
