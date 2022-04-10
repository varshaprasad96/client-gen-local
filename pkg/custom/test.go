package custom

// func (w *wrapped{{.Api.Name}}) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
// 	ctx, err := w.checkCluster(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return w.delegate.Delete(ctx, name, opts)
// }

// func (w *wrapped{{.Api.Name}}) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
// 	ctx, err := w.checkCluster(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return w.delegate.Delete(ctx, opts, listOpts)
// }

// func (w *wrapped{{.Api.Name}}) Get(ctx context.Context, name string, opts metav1.GetOptions) (*{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}, error) {
// 	ctx, err := w.checkCluster(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return w.delegate.Get(ctx, opts, listOpts)
// }

// func (w *wrapped{{.Api.Name}}) List(ctx context.Context, opts metav1.ListOptions) (*{{.Api.PkgName}}api{{.Api.Version}}.{{.Api.Name}}List, error) {
// 	ctx, err := w.checkCluster(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return w.delegate.List(ctx, opts)
// }
