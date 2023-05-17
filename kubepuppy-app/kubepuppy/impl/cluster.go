package impl

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (c *Cluster) RefreshRolesAndBindings(ctx context.Context) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	roles := make(map[string]Role)
	roleBindings := make(map[string]RoleBinding)

	clusterRoles, err := (*c.KubeClient).RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, clusterRole := range clusterRoles.Items {
		cr, err := NewClusterRoleFromManifest(clusterRole)

		if err != nil {
			return err
		}

		roles[cr.QualifiedName()] = cr
	}

	namespacedRoles, err := (*c.KubeClient).RbacV1().Roles("").List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, namespacedRole := range namespacedRoles.Items {
		nr, err := NewNamespacedRoleFromManifest(namespacedRole)

		if err != nil {
			return err
		}

		roles[nr.QualifiedName()] = nr
	}

	clusterRoleBindings, err := (*c.KubeClient).RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		crb, err := NewClusterRoleBindingFromManifest(clusterRoleBinding)

		if err != nil {
			return err
		}

		roleBindings[crb.QualifiedName()] = crb
	}

	namespacedRoleBindings, err := (*c.KubeClient).RbacV1().RoleBindings("").List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, namespacedRoleBinding := range namespacedRoleBindings.Items {
		crb, err := NewNamespacedRoleBindingFromManifest(namespacedRoleBinding)

		if err != nil {
			return err
		}

		roleBindings[crb.QualifiedName()] = crb
	}

	c.Roles = roles
	c.RoleBindings = roleBindings

	return nil
}

func (c *Cluster) FindRoleBindingsForSubject(target Principal) ([]RoleBinding, []Role, error) {

	roleBindings := make([]RoleBinding, 0)

	for _, roleBinding := range c.RoleBindings {
		if roleBinding.Matches(target) {
			roleBindings = append(roleBindings, roleBinding)
		}
	}

	roles := make([]Role, 0)

	for _, roleBinding := range roleBindings {
		role, ok := c.Roles[roleBinding.QualifiedRoleName()]

		if !ok {
			return nil, nil, fmt.Errorf("role not found: '%v' for binding: '%v'", roleBinding.QualifiedRoleName(), roleBinding.QualifiedName())
		}

		roles = append(roles, role)
	}

	return roleBindings, roles, nil
}
