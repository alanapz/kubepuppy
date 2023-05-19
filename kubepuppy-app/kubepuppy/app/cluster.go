package app

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func InitialiseCluster(ctx context.Context) (*Cluster, error) {

	kubeconfigPath, ok := os.LookupEnv("KUBEPUPPY_KUBECONFIG")
	if !ok {
		return nil, errors.New("KUBEPUPPY_KUBECONFIG not defined")
	}

	kubeconfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig from file: '%v': %v", kubeconfigPath, err)
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)

	if err != nil {
		return nil, fmt.Errorf("unable to initialise Kubernetes API client from file: '%v': %v", kubeconfigPath, err)
	}

	serverCaPath, ok := os.LookupEnv("KUBEPUPPY_SERVER_CA_FILE")

	if !ok {
		return nil, errors.New("KUBEPUPPY_SERVER_CA_FILE not defined")
	}

	serverAuthorities, err := loadCertificate(serverCaPath)

	if err != nil {
		return nil, err
	}

	clientCaPath, ok := os.LookupEnv("KUBEPUPPY_CLIENT_CA_FILE")

	if !ok {
		return nil, errors.New("KUBEPUPPY_CLIENT_CA_FILE not defined")
	}

	clientAuthorities, err := loadCertificate(clientCaPath)

	if err != nil {
		return nil, err
	}

	c := Cluster{}
	c.KubeClient = clientset
	c.ServerAuthorities = serverAuthorities
	c.ClientAuthorities = clientAuthorities

	c.lock.Lock()
	defer c.lock.Unlock()
	err = c.RefreshRolesAndBindings(ctx)

	if err != nil {
		return nil, fmt.Errorf("unable to refresh role bindings: %v", err)
	}

	return &c, nil
}

func loadCertificate(certPath string) ([]*x509.Certificate, error) {

	certData, err := os.ReadFile(certPath)

	if err != nil {
		return nil, err
	}

	certBlock, _ := pem.Decode(certData)

	if certBlock == nil {
		return nil, fmt.Errorf("invalid certificate file: '%v'", certPath)
	}

	certs, err := x509.ParseCertificates(certBlock.Bytes)

	if err != nil {
		return nil, fmt.Errorf("unable to parse certificates from: '%v': %v", certPath, err)
	}

	return certs, nil
}

func (c *Cluster) LockForRead(fn func()) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	fn()
}

func (c *Cluster) LockForWrite(fn func()) {
	c.lock.Lock()
	defer c.lock.Unlock()
	fn()
}

//func (c *Cluster) FindRoleBindingsForSubject(target Principal) ([]RoleBinding, []Role, error) {
//
//	roleBindings := make([]RoleBinding, 0)
//
//	for _, roleBinding := range c.RoleBindings {
//		if roleBinding.Matches(target) {
//			roleBindings = append(roleBindings, roleBinding)
//		}
//	}
//
//	roles := make([]Role, 0)
//
//	for _, roleBinding := range roleBindings {
//		role, ok := c.Roles[roleBinding.RoleId()]
//
//		if !ok {
//			return nil, nil, fmt.Errorf("role not found: '%v' for binding: '%v'", roleBinding.RoleId(), roleBinding.RoleBindingId())
//		}
//
//		roles = append(roles, role)
//	}
//
//	return roleBindings, roles, nil
//}
