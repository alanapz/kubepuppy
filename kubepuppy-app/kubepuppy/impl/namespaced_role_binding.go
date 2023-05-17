package impl

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	"kubepuppy/utils"
)

type NamedspacedRoleBinding struct {
	QualifiedName_     string
	QualifiedRoleName_ string
	Manifest           v1.RoleBinding
	subjects           []Subject
}

var _ RoleBinding = (*NamedspacedRoleBinding)(nil)

func NewNamespacedRoleBindingFromManifest(manifest v1.RoleBinding) (*NamedspacedRoleBinding, error) {

	subjects := make([]Subject, len(manifest.Subjects))

	for k, v := range manifest.Subjects {
		subject, err := NewSubjectFromManifest(v)
		if err != nil {
			return nil, err
		}
		subjects[k] = subject
	}

	return &NamedspacedRoleBinding{
		QualifiedName_:     fmt.Sprintf("RoleBinding[%v/%v]", manifest.ObjectMeta.Namespace, manifest.Name),
		QualifiedRoleName_: fmt.Sprintf("%v[%v/%v]", manifest.RoleRef.Kind, manifest.ObjectMeta.Namespace, manifest.RoleRef.Name),
		Manifest:           manifest,
		subjects:           subjects}, nil
}

func (rb NamedspacedRoleBinding) QualifiedName() string {
	return rb.QualifiedName_
}

func (rb NamedspacedRoleBinding) QualifiedRoleName() string {
	return rb.QualifiedRoleName_
}

func (rb NamedspacedRoleBinding) Subjects() []Subject {
	return rb.subjects
}

func (rb NamedspacedRoleBinding) Matches(principal Principal) bool {
	return utils.Any(rb.Subjects(), func(subject Subject) bool {
		return principal.Matches(subject)
	})
}
