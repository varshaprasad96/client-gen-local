package flag

import "github.com/spf13/pflag"

// Flags - Options accepted by generator
type Flags struct {
	OutputDir        string
	InputDir         string
	ClientsetAPIPath string
	Version          string
}

func (f *Flags) AddTo(flagset *pflag.FlagSet) {
	// TODO: FIgure out if its worth defaulting it to pkg/api/...
	flagset.StringVar(&f.InputDir, "input-dir", "", "Input directory where types are defined")
	flagset.StringVar(&f.OutputDir, "output-dir", "/output", "Output directory where wrapped clients will be generatoed")
	flagset.StringVar(&f.ClientsetAPIPath, "clientset-api-path", "/apis", "package path where clients are generated")
	flagset.StringVar(&f.Version, "version", "v1", "API version")
}
