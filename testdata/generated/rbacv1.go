// +build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package generated


import (
	"ctx"
	rbacapiv1 "github.com/varshaprasad96/client-gen/testdata/pkg/apis/rbac/v1"
	rbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"

	"github.com/kcp-dev/kcp-client-wrappers/kcp"
)

type wrappedRbacV1 struct {
	cluster  string
	delegate rbacv1.RbacV1Interface
}

func (w *wrappedRbacV1) RESTClient() rest.Interface {
	return w.delegate.RESTClient()
}

func (w *wrappedrbac) checkCluster(ctx context.Context) (context.Context, error) {
	ctxCluster, ok := kcp.ClusterFromContext(ctx)
	if !ok {
		return kcp.WithCluster(ctx, w.cluster), nil
	} else if ctxCluster != w.cluster {
		return ctx, fmt.Errorf("cluster mismatch: context=%q, client=%q", ctxCluster, w.cluster)
	}
	return ctx, nil
}

func (w *wrappedRbacV1) ClusterRoles() rbacv1.ClusterRoleInterface {
	return &wrappedClusterRole{
		cluster:  w.cluster,
		delegate: w.delegate.ClusterRoles(),
	}
}

type wrappedClusterRole struct {
	cluster  string
	delegate rbac.ClusterRoleInterface
}

func (w *wrappedClusterRole) Create(ctx context.Context, clusterRole *rbacapiv1.ClusterRole, opts metav1.CreateOptions) (*rbacapiv1.ClusterRole, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Create(ctx, clusterRole}, opts)
}

func (w *wrappedClusterRole) Update(ctx context.Context, clusterRole} *rbacapiv1.ClusterRole, opts metav1.UpdateOptions) (*rbacapiv1.ClusterRole, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Update(ctx, clusterRole}, opts)
}

func (w *wrappedClusterRole) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Delete(ctx, name, opts)
}

func (w *wrappedClusterRole) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Delete(ctx, opts, listOpts)
}

func (w *wrappedClusterRole) Get(ctx context.Context, name string, opts metav1.GetOptions) (*rbacapiv1.ClusterRole, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Get(ctx, opts, listOpts)
}

func (w *wrappedClusterRole) List(ctx context.Context, opts metav1.ListOptions) (*rbacapiv1.ClusterRoleList, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.List(ctx, opts)
}

func (w *wrappedClusterRole) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Watch(ctx, opts)
}

func (w *wrappedClusterRole) Patch(ctx context.Context, name string, pt apiTypes.PatchapiType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *rbacapiv1.ClusterRole, err error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Patch(ctx, name, pt, data, opts, subresources)
}

func (w *wrappedRbacV1) ClusterRoleBindings() rbacv1.ClusterRoleBindingInterface {
	return &wrappedClusterRoleBinding{
		cluster:  w.cluster,
		delegate: w.delegate.ClusterRoleBindings(),
	}
}

type wrappedClusterRoleBinding struct {
	cluster  string
	delegate rbac.ClusterRoleBindingInterface
}

func (w *wrappedClusterRoleBinding) Create(ctx context.Context, clusterRoleBinding *rbacapiv1.ClusterRoleBinding, opts metav1.CreateOptions) (*rbacapiv1.ClusterRoleBinding, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Create(ctx, clusterRoleBinding}, opts)
}

func (w *wrappedClusterRoleBinding) Update(ctx context.Context, clusterRoleBinding} *rbacapiv1.ClusterRoleBinding, opts metav1.UpdateOptions) (*rbacapiv1.ClusterRoleBinding, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Update(ctx, clusterRoleBinding}, opts)
}

func (w *wrappedClusterRoleBinding) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Delete(ctx, name, opts)
}

func (w *wrappedClusterRoleBinding) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Delete(ctx, opts, listOpts)
}

func (w *wrappedClusterRoleBinding) Get(ctx context.Context, name string, opts metav1.GetOptions) (*rbacapiv1.ClusterRoleBinding, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Get(ctx, opts, listOpts)
}

func (w *wrappedClusterRoleBinding) List(ctx context.Context, opts metav1.ListOptions) (*rbacapiv1.ClusterRoleBindingList, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.List(ctx, opts)
}

func (w *wrappedClusterRoleBinding) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Watch(ctx, opts)
}

func (w *wrappedClusterRoleBinding) Patch(ctx context.Context, name string, pt apiTypes.PatchapiType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *rbacapiv1.ClusterRoleBinding, err error) {
	ctx, err := w.checkCluster(ctx)
	if err != nil {
		return nil, err
	}
	return w.delegate.Patch(ctx, name, pt, data, opts, subresources)
}