package custom

import (
	"fmt"
	"go/types"
	"io"
	"strings"
	"text/template"

	gentype "k8s.io/code-generator/cmd/client-gen/types"

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

type wrapperInterface struct {
	InterfaceName string
	InputPath     string
	APIs          []Api
	codeWriter    *codeWriter
}

func groupVersionToAPIs(gvs []gentype.GroupVersions) []Api {
	result := make([]Api, 0)

	for _, gv := range gvs {
		if len(gv.Versions) <= 0 {
			continue
		}
		a := &Api{
			Name:    gv.Group.String(),
			Version: string(gv.Versions[0].Version),
			PkgName: gv.PackageName,
		}
		a.setCased()
		result = append(result, *a)
	}
	return result
}

func NewWrappedInterface(interfaceName string, inputPath string, gvs []gentype.GroupVersions, cw *codeWriter) (*wrapperInterface, error) {
	if len(gvs) == 0 {
		return nil, fmt.Errorf("no group version pair is defined")
	}

	apis := groupVersionToAPIs(gvs)
	return &wrapperInterface{
		InputPath:     inputPath,
		InterfaceName: interfaceName,
		APIs:          apis,
		codeWriter:    cw,
	}, nil
}

func (w *wrapperInterface) writeMethods() error {
	templ, err := template.New("wrapper").Parse(wrappedInterfacesTempl)
	if err != nil {
		return err
	}

	err = templ.Execute(w.codeWriter.out, w)
	if err != nil {
		return err
	}

	return nil
}

type Api struct {
	Name       string
	Version    string
	PkgName    string
	codeWriter *codeWriter

	PkgNameUpperFirst string
	VersionUpperFirst string
	NameLowerFirst    string
}

type api struct {
	Name       string
	Version    string
	PkgName    string
	codeWriter *codeWriter

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
	NameUpperFirst    string
	VersionUpperFirst string
	Version           string
	codeWriter        *codeWriter
}

func NewPackage(root *loader.Package, apiPath, clientPath, version, group string, cocodeWriter *codeWriter) (*packages, error) {
	p := &packages{
		Name:       group,
		APIPath:    apiPath,
		Version:    version,
		ClientPath: clientPath,
		codeWriter: cocodeWriter,
	}
	p.setCased()
	return p, nil
}

func (p *packages) setCased() {
	p.NameUpperFirst = upperFirst(p.Name)
	p.VersionUpperFirst = upperFirst(p.Version)
}

func NewAPI(root *loader.Package, info *markers.TypeInfo, version, group string, cocodeWriter *codeWriter) (*api, error) {
	typeInfo := root.TypesInfo.TypeOf(info.RawSpec.Name)
	if typeInfo == types.Typ[types.Invalid] {
		return nil, fmt.Errorf("unknown type: %s", info.Name)
	}

	api := &api{
		Name:       info.RawSpec.Name.Name,
		Version:    version,
		PkgName:    group,
		codeWriter: cocodeWriter,
	}

	api.setCased()
	return api, nil

}

func (a *api) writeMethods() error {
	templ, err := template.New("wrapper").Parse(wrapperTempl)
	if err != nil {
		return err
	}

	err = templ.Execute(a.codeWriter.out, a)
	if err != nil {
		return err
	}

	return nil
}

func (a *api) setCased() {
	a.PkgNameUpperFirst = upperFirst(a.PkgName)
	a.VersionUpperFirst = upperFirst(a.Version)
	a.NameLowerFirst = lowerFirst(a.Name)
}

func (a *Api) setCased() {
	a.PkgNameUpperFirst = upperFirst(a.PkgName)
	a.VersionUpperFirst = upperFirst(a.Version)
	a.NameLowerFirst = lowerFirst(a.Name)
}

func (p *packages) writeCommonContent() error {
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
