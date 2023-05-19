package app

import (
	"crypto/x509"
	"k8s.io/client-go/kubernetes"
	"sync"
)

type ClusterConfig struct {
	KubeConfigPath      string
	ClientAuthorityPath string
}

type Cluster struct {
	lock              sync.RWMutex
	KubeClient        *kubernetes.Clientset
	ServerAuthorities []*x509.Certificate
	ClientAuthorities []*x509.Certificate
	Roles             map[string]*Role
	RoleBindings      map[string]*RoleBinding
	Subjects          map[string]*Subject
}

type Principal interface {
	Matches(target Subject) bool
}
