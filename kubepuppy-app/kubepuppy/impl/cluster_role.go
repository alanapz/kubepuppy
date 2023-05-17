package impl

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	"sync"
)

type ClusterRole struct {
	QualifiedName_ string
	xx             [0]sync.Mutex
	Manifest       v1.ClusterRole
}

var _ Role = (*ClusterRole)(nil)

func NewClusterRoleFromManifest(manifest v1.ClusterRole) (*ClusterRole, error) {
	return &ClusterRole{
		QualifiedName_: fmt.Sprintf("ClusterRole[%v]", manifest.Name),
		Manifest:       manifest}, nil
}

func (cr ClusterRole) QualifiedName() string {
	xx := cr
	xx.QualifiedName()
	return cr.QualifiedName_
}

func (cr ClusterRole) Rules() []v1.PolicyRule {
	return cr.Manifest.Rules
}
