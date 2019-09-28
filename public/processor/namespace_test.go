package processor

import "testing"

func TestNamespaceBeforeResource(t *testing.T) {
	TestFromDir(t, "namespace")
}


func TestNamespaceBeforeResourceForce(t *testing.T) {
	TestFromDir(t, "namespace-force")
}

func TestNamespaceBeforeResourceClusterRole(t *testing.T) {
	TestFromDir(t, "namespace-clusterrole")
}