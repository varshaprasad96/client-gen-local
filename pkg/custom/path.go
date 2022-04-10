package custom

import "k8s.io/code-generator/cmd/client-gen/types"

type gvkConverter struct {
	importBasePath string
	dirPath        string
}

func (g *gvkConverter) getGroupVersion() []types.GroupVersions {
	return nil
}
