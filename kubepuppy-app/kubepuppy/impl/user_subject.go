package impl

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
)

type UserSubject struct {
	Name      string
	Manifest_ *v1.Subject
}

var _ Subject = (*UserSubject)(nil)

func NewUserSubject(username string) UserSubject {
	return UserSubject{Name: username}
}

func NewUserSubjectFromManifest(subject v1.Subject) (*UserSubject, error) {
	if subject.Kind != "User" {
		return nil, fmt.Errorf("unexpected subject kind: %v", subject.Kind)
	}
	return &UserSubject{Name: subject.Name, Manifest_: &subject}, nil
}

func (us UserSubject) Matches(target Subject) bool {
	targetUser, isUser := target.(UserSubject)
	return isUser && us.Name == targetUser.Name
}

func (us UserSubject) Manifest() *v1.Subject {
	return us.Manifest_
}
