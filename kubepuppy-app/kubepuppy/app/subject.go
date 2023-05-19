package app

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	"kubepuppy/utils"
)

type SubjectKind = string

const UserSubjectKind SubjectKind = "User"
const GroupSubjectKind SubjectKind = "Group"
const ServiceAccountSubjectKind SubjectKind = "ServiceAccount"

type Subject struct {
	SubjectId    string
	Kind         SubjectKind
	Name         string
	Namespace    string
	Roles        map[string]*Role
	RoleBindings map[string]*RoleBinding
}

func GenerateSubjectId(subject v1.Subject) (*string, error) {

	if subject.Kind == UserSubjectKind {
		return utils.Ptr(fmt.Sprintf("%v::%v", UserSubjectKind, subject.Name)), nil
	}

	if subject.Kind == GroupSubjectKind {
		return utils.Ptr(fmt.Sprintf("%v::%v", GroupSubjectKind, subject.Name)), nil
	}

	if subject.Kind == ServiceAccountSubjectKind {

		if subject.Namespace == "" {
			return nil, fmt.Errorf("unexpected namespace: '%v' for service account subject: '%v'", subject.Namespace, subject.Name)
		}

		return utils.Ptr(fmt.Sprintf("%v::%v::%v", ServiceAccountSubjectKind, subject.Namespace, subject.Name)), nil
	}

	return nil, fmt.Errorf("unexpected subject kind: %v", subject.Kind)
}
