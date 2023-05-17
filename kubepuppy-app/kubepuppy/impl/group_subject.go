package impl

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
)

type GroupSubject struct {
	Name      string
	Manifest_ *v1.Subject
}

var _ Subject = (*GroupSubject)(nil)

func NewGroupSubject(groupName string) GroupSubject {
	return GroupSubject{Name: groupName}
}

func NewGroupSubjectFromManifest(subject v1.Subject) (*GroupSubject, error) {
	if subject.Kind != "Group" {
		return nil, fmt.Errorf("unexpected subject kind: %v", subject.Kind)
	}
	return &GroupSubject{Name: subject.Name, Manifest_: &subject}, nil
}

func (gs GroupSubject) Matches(target Subject) bool {
	targetGroup, isGroup := target.(*GroupSubject)
	return isGroup && gs.Name == targetGroup.Name
}

func (gs GroupSubject) Manifest() *v1.Subject {
	return gs.Manifest_
}
