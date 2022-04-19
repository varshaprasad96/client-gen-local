package generated

import (
	"github.com/kcp-dev/kcp-client-wrappers/kcp"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"

	rbacv1 "github.com/varshaprasad96/client-gen/testdata/pkg/apis/rbac/v1"

	appsv1 "github.com/varshaprasad96/client-gen/testdata/pkg/apis/apps/v1"
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

func (w *wrappedInterface) Discovery() discovery.DiscoveryInterface {
	return w.delegate.Discovery()
}

func (w *wrappedInterface) RbacV1() rbacv1.RbacV1Interface {
	return &wrappedRbacV1{
		cluster:  w.cluster,
		delegate: w.delegate.RbacV1(),
	}
}

func (w *wrappedInterface) AppsV1() appsv1.AppsV1Interface {
	return &wrappedAppsV1{
		cluster:  w.cluster,
		delegate: w.delegate.AppsV1(),
	}
}
