package impl

import (
	"crypto/x509"
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"sync"
)

type ClusterConfig struct {
	KubeConfigPath      string
	ClientAuthorityPath string
}

type Cluster struct {
	KubeClient        *kubernetes.Clientset
	ServerAuthorities []*x509.Certificate
	ClientAuthorities []*x509.Certificate
	Roles             map[string]Role
	RoleBindings      map[string]RoleBinding
	lock              sync.RWMutex
}

type RoleBinding interface {
	QualifiedName() string
	QualifiedRoleName() string
	Subjects() []Subject
	Matches(target Principal) bool
}

type Role interface {
	QualifiedName() string
	Rules() []v1.PolicyRule
}

type Principal interface {
	Matches(target Subject) bool
}

type Subject interface {
	Principal
	Manifest() *v1.Subject
}
