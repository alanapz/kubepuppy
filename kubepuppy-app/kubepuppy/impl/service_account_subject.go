package impl

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
)

type ServiceAccountSubject struct {
	AccountName string
	Namespace   string
	Manifest_   *v1.Subject
}

var _ Subject = (*ServiceAccountSubject)(nil)

func NewServiceAccountSubject(accountName string, namespace string) ServiceAccountSubject {
	return ServiceAccountSubject{AccountName: accountName, Namespace: namespace}
}

func NewServiceAccountSubjectFromManifest(subject v1.Subject) (*ServiceAccountSubject, error) {
	if subject.Kind != "ServiceAccount" {
		return nil, fmt.Errorf("unexpected subject kind: %v", subject.Kind)
	}
	return &ServiceAccountSubject{AccountName: subject.Name, Namespace: subject.Namespace, Manifest_: &subject}, nil
}

func (sas ServiceAccountSubject) Matches(target Subject) bool {
	targetServiceAccount, isServiceAccount := target.(ServiceAccountSubject)
	return isServiceAccount && sas.AccountName == targetServiceAccount.AccountName && sas.Namespace == targetServiceAccount.Namespace
}

func (sas ServiceAccountSubject) Manifest() *v1.Subject {
	return sas.Manifest_
}
