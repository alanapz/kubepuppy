package impl

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
)

type NamespacedRole struct {
	QualifiedName_ string
	Manifest       v1.Role
}

var _ Role = (*ClusterRole)(nil)

func NewNamespacedRoleFromManifest(manifest v1.Role) (*NamespacedRole, error) {
	return &NamespacedRole{
		QualifiedName_: fmt.Sprintf("Role[%v/%v]", manifest.ObjectMeta.Namespace, manifest.Name),
		Manifest:       manifest,
	}, nil
}

func (nr NamespacedRole) QualifiedName() string {
	return nr.QualifiedName_
}

func (nr NamespacedRole) Rules() []v1.PolicyRule {
	return nr.Manifest.Rules
}
