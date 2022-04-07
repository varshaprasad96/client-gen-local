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
	RuleDefinition = markers.Must(markers.MakeDefinition("genclient", markers.DescribesType, Rule{}))
)

type Rule struct {
	// Types defines a test type
	Types []string `marker:",optional"`
	// Groups specifies the API groups that this rule encompasses.
	Groups []string `marker:",optional"`
	// Resources specifies the API resources that this rule encompasses.
	Resources []string `marker:",optional"`
	// ResourceNames specifies the names of the API resources that this rule encompasses.
	//
	// Create requests cannot be restricted by resourcename, as the object's name
	// is not known at authorization time.
	ResourceNames []string `marker:",optional"`
	// Verbs specifies the (lowercase) kubernetes API verbs that this rule encompasses.
	Verbs []string `marker:",optional"`
	// URL specifies the non-resource URLs that this rule encompasses.
	URLs []string `marker:"urls,optional"`
	// Namespace specifies the scope of the Rule.
	// If not set, the Rule belongs to the generated ClusterRole.
	// If set, the Rule belongs to a Role, whose namespace is specified by this field.
	Namespace string `marker:",optional"`
}

// ruleKey represents the resources and non-resources a Rule applies.
type ruleKey struct {
	Types         string
	Groups        string
	Resources     string
	ResourceNames string
	URLs          string
}

func (key ruleKey) String() string {
	return fmt.Sprintf("%s + %s + %s + %s + %s", key.Types, key.Groups, key.Resources, key.ResourceNames, key.URLs)
}

type codeWriter struct {
	out io.Writer
}

// Line writes a single line.
func (c *codeWriter) Line(line string) {
	fmt.Fprintln(c.out, line)
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

	if strings.HasSuffix(typeInfo.String(), "Status") || strings.HasSuffix(typeInfo.String(), "Spec") {
		return
	}

	c.Line("// DONOT EDIT!!")
	c.Line(newClientsetForConfigTemplate)

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
