package impl

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	"kubepuppy/utils"
)

type ClusterRoleBinding struct {
	QualifiedName_     string
	QualifiedRoleName_ string
	Manifest           v1.ClusterRoleBinding
	subjects           []Subject
}

var _ RoleBinding = (*ClusterRoleBinding)(nil)

func NewClusterRoleBindingFromManifest(manifest v1.ClusterRoleBinding) (*ClusterRoleBinding, error) {

	subjects := make([]Subject, len(manifest.Subjects))

	for k, v := range manifest.Subjects {
		subject, err := NewSubjectFromManifest(v)
		if err != nil {
			return nil, err
		}
		subjects[k] = subject
	}

	return &ClusterRoleBinding{
		QualifiedName_:     fmt.Sprintf("ClusterRoleBinding[%v]", manifest.Name),
		QualifiedRoleName_: fmt.Sprintf("%v[%v]", manifest.RoleRef.Kind, manifest.RoleRef.Name),
		Manifest:           manifest,
		subjects:           subjects}, nil
}

func (crb ClusterRoleBinding) QualifiedName() string {
	return crb.QualifiedName_
}

func (crb ClusterRoleBinding) QualifiedRoleName() string {
	return crb.QualifiedRoleName_
}

func (crb ClusterRoleBinding) Subjects() []Subject {
	return crb.subjects
}

func (crb ClusterRoleBinding) Matches(principal Principal) bool {
	return utils.Any(crb.Subjects(), func(subject Subject) bool {
		return principal.Matches(subject)
	})
}
