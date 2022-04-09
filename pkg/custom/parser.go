package custom

import (
	"fmt"
	"go/types"
	"io"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

var (
	// RuleDefinition is a marker for defining rules
	RuleDefinition = markers.Must(markers.MakeDefinition("genclient", markers.DescribesType, placeholder{}))
)

// Assigning marker's output to a placeholder struct, to verify to
// typecast the result and make sure if it exists for the type.
type placeholder struct{}

type codeWriter struct {
	out io.Writer
}

// Line writes a single line.
func (c *codeWriter) Line(line string) {
	fmt.Fprintln(c.out, line)
}

// Linef writes a single line with formatting (as per fmt.Sprintf).
func (c *codeWriter) Linef(line string, args ...interface{}) {
	fmt.Fprintf(c.out, line+"\n", args...)
}

type configMethodWriter struct {
	*importsList
	pkg loader.Package
	*codeWriter
}

func (c *configMethodWriter) GenerateConfigMethod(root *loader.Package, info *markers.TypeInfo) {
	typeInfo := root.TypesInfo.TypeOf(info.RawSpec.Name)
	if typeInfo == types.Typ[types.Invalid] {
		root.AddError(loader.ErrFromNode(fmt.Errorf("unknown type: %s", info.Name), info.RawSpec))
	}

	// Flaky condition. We can remove it because of isEnabledMethod(), but keeping this as a double check.
	if strings.HasSuffix(typeInfo.String(), "Status") || strings.HasSuffix(typeInfo.String(), "Spec") {
		return
	}

	c.Line("// DONOT EDIT!!")
	c.Linef(newClientsetForConfigTemplate, (&namingInfo{typeInfo: typeInfo}).Syntax(root, c.importsList))

	// Add the imports
	importsList := []string{"k8s.io/client-go/kubernetes", "github.com/kcp-dev/kcp-client-wrappers/kcp", "k8s.io/client-go/rest"}
	for _, imp := range importsList {
		c.NeedImport(imp)
	}

}

var newClientsetForConfigTemplate = `
// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable, 
// NewForConfig will generate a rate-limiter in configShallowCopy.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*ClusterClient, error) {
	client, err := rest.HTTPClientFor(config)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP client: %w", err)
	}

	clusterRoundTripper := kcp.NewClusterRoundTripper(client.Transport)
	client.Transport = clusterRoundTripper

	delegate, err := kubernetes.NewForConfigAndClient(config, client)
	if err != nil {
		return nil, fmt.Errorf("error creating delegate clientset: %w", err)
	}

	return &ClusterClient{
		delegate: delegate,
	}, nil

}
`
