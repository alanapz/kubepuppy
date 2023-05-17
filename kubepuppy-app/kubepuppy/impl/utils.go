package impl

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
)

func NewSubjectFromManifest(subject v1.Subject) (Subject, error) {

	if subject.Kind == "User" {
		userSubject, err := NewUserSubjectFromManifest(subject)
		if err != nil {
			return nil, err
		}
		return userSubject, nil
	}

	if subject.Kind == "Group" {
		groupSubject, err := NewGroupSubjectFromManifest(subject)
		if err != nil {
			return nil, err
		}
		return groupSubject, nil
	}

	if subject.Kind == "ServiceAccount" {
		serviceAccountSubject, err := NewServiceAccountSubjectFromManifest(subject)
		if err != nil {
			return nil, err
		}
		return serviceAccountSubject, nil
	}

	return nil, fmt.Errorf("unexpected subject kind: %v", subject.Kind)
}
