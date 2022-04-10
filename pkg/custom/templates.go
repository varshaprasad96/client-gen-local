package custom

const commonTempl = `
import (
	{{.Name}}api{{.Version}} "{{.APIPath}}"
	{{.Name}}{{.Version}} "{{.ClientPath}}"
)

func NewForConfig(config *rest.Config) (*ClusterClient, error) {
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

type ClusterClient struct {
	delegate kubernetes.Interface
}

func (c *ClusterClient) Cluster(cluster string) kubernetes.Interface {
	return &wrappedInterface{
		cluster:  cluster,
		delegate: c.delegate,
	}
}

type wrappedInterface struct {
	cluster  string
	delegate kubernetes.Interface
}

func (w *wrappedInterface) {{.NameUpperFirst}}{{.VersionUpperFirst}}() {{.Name}}{{.Version}}.{{.NameUpperFirst}}{{.VersionUpperFirst}}Interface {
	return &wrapped{{.NameUpperFirst}}{{.VersionUpperFirst}}{
		cluster:  w.cluster,
		delegate: w.delegate.{{.NameUpperFirst}}{{.VersionUpperFirst}},
	}
}

type wrapped{{.NameUpperFirst}}{{.VersionUpperFirst}} struct {
	cluster  string
	delegate {{.Name}}{{.Version}}.{{.NameUpperFirst}}{{.VersionUpperFirst}}Interface
}

func (w *wrapped{{.NameUpperFirst}}{{.VersionUpperFirst}}) RESTClient() rest.Interface {
	//TODO
	panic("no")
}

`

const wrapperTempl = `
func (w *wrapped{{.Api.PkgNameUpperFirst}}{{.Api.VersionUpperFirst}}) {{.Api.Name}}s() {{.Api.PkgName}}{{.Api.Version}}.{{.Api.Name}}Interface {
	return &wrapped{{.Api.Name}}{
		cluster:  w.cluster,
		delegate: w.delegate.{{.Api.Name}}s(),
	}
}

type wrapped{{.Api.Name}} struct {
	cluster  string
	delegate {{.Api.PkgName}}.{{.Api.Name}}Interface
}

func (w *wrappedInterface) {{.Api.PkgNameUpperFirst}}{{.Api.VersionUpperFirst}}() {{.Api.PkgName}}{{.Api.Version}}.{{.Api.PkgNameUpperFirst}}{{.Api.VersionUpperFirst}}Interface {
	return &wrapped{{.Api.PkgNameUpperFirst}}{{.Api.VersionUpperFirst}}{
		cluster:  w.cluster,
		delegate: w.delegate.{{.Api.PkgNameUpperFirst}}{{.Api.VersionUpperFirst}}(),
	}
}

func (w *wrapped{{.Api.Name}}) checkCluster(ctx context.Context) (context.Context, error) {
	ctxCluster, ok := kcp.ClusterFromContext(ctx)
	if !ok {
		return kcp.WithCluster(ctx, w.cluster), nil
	} else if ctxCluster != w.cluster {
		return ctx, fmt.Errorf("cluster mismatch: context=%q, client=%q", ctxCluster, w.cluster)
	}
	return ctx, nil
}

func (w *wrapped{{.Api.Name}}) Create(ctx context.Context, {{.Api.NameLowerFirst}} *{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}, opts metav1.CreateOptions) (*{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Create(ctx, {{.Api.NameLowerFirst}}}, opts)
}

func (w *wrapped{{.Api.Name}}) Update(ctx context.Context, {{.Api.NameLowerFirst}}} *{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}, opts metav1.UpdateOptions) (*{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Update(ctx, {{.Api.NameLowerFirst}}}, opts)
}

func (w *wrapped{{.Api.Name}}) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Delete(ctx, name, opts)
}

func (w *wrapped{{.Api.Name}}) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Delete(ctx, opts, listOpts)
}

func (w *wrapped{{.Api.Name}}) Get(ctx context.Context, name string, opts metav1.GetOptions) (*{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Get(ctx, opts, listOpts)
}

func (w *wrapped{{.Api.Name}}) List(ctx context.Context, opts metav1.ListOptions) (*{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}List, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.List(ctx, opts)
}

func (w *wrapped{{.Api.Name}}) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Watch(ctx, opts)
}

func (w *wrapped{{.Api.Name}}) Patch(ctx context.Context, name string, pt apiTypes.PatchapiType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}, err error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Patch(ctx, name, pt, data, opts, subresources)
}
`
