package custom

import (
	"fmt"
	"go/types"
	"io"
	"strings"
	"text/template"

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

type api struct {
	Name    string
	Version string
	PkgName string

	PkgNameUpperFirst string
	VersionUpperFirst string
	NameLowerFirst    string
}

// TODO: move this to internal, so the exported fields are not accessible to users.
// Need to export them for executing as template
type packages struct {
	Name              string
	APIPath           string
	ClientPath        string
	Api               *api
	NameUpperFirst    string
	VersionUpperFirst string
	Version           string
	codeWriter        *codeWriter
}

func NewPackage(root *loader.Package, apiPath, clientPath, version, group string, cocodeWriter *codeWriter) error {
	p := &packages{
		Name:       group,
		APIPath:    apiPath,
		Version:    version,
		ClientPath: clientPath,
		codeWriter: cocodeWriter,
	}
}

func NewAPI(root *loader.Package, info *markers.TypeInfo, apiPath, clientPath, version, group string, cocodeWriter *codeWriter) (*packages, error) {
	typeInfo := root.TypesInfo.TypeOf(info.RawSpec.Name)
	if typeInfo == types.Typ[types.Invalid] {
		return nil, fmt.Errorf("unknown type: %s", info.Name)
	}

	api := &api{
		Name:    info.RawSpec.Name.Name,
		Version: version,
		PkgName: group,
	}

	p := &packages{
		Name:       group,
		APIPath:    apiPath,
		Version:    version,
		ClientPath: clientPath,
		Api:        api,
		codeWriter: cocodeWriter,
	}

	p.setCased()
	return p, nil

}

func (p *packages) writeMethods() error {
	templ, err := template.New("wrapper").Parse(wrapperTempl)
	if err != nil {
		return err
	}

	err = templ.Execute(p.codeWriter.out, p)
	if err != nil {
		return err
	}

	return nil
}

func (p *packages) setCased() {
	p.NameUpperFirst = upperFirst(p.Name)
	p.VersionUpperFirst = upperFirst(p.Version)
	p.Api.setCased()
}

func (a *api) setCased() {
	a.PkgNameUpperFirst = upperFirst(a.PkgName)
	a.VersionUpperFirst = upperFirst(a.Version)
	a.NameLowerFirst = lowerFirst(a.Name)
}

func (p *packages) writeCommonContent(out io.Writer) error {
	templ, err := template.New("wrapper").Parse(commonTempl)
	if err != nil {
		return err
	}

	err = templ.Execute(p.codeWriter.out, p)
	if err != nil {
		return err
	}

	return nil
}

func lowerFirst(s string) string {
	return strings.ToLower(string(s[0])) + s[1:]
}

func upperFirst(s string) string {
	return strings.ToUpper(string(s[0])) + s[1:]
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

var clusterClientDef = `
type ClusterClient struct {
	delegate kubernetes.Interface
}
`

var clusterMethod = `
func (c *ClusterClient) Cluster(cluster string) kubernetes.Interface {
	return &wrappedInterface{
		cluster:  cluster,
		delegate: c.delegate,
	}
}
`

var wrappedInterface = `
type wrappedInterface struct {
	cluster  string
	delegate kubernetes.Interface
}
`
