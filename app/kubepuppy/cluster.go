package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"kubepuppy/utils"
	"log"
	"os"
	"sync"
)

type ClusterConfig struct {
	KubeconfigPath string
}

type PuppyCluster struct {
	lock            *sync.RWMutex
	kubeClient      *kubernetes.Clientset
	ServerAuthority []*x509.Certificate
	RoleBindings    map[string]*PuppyRoleBinding
	Roles           map[string]*PuppyRole
}

type PuppyRoleBinding struct {
	Kind            string
	RoleBindingName string
	RoleName        string
	Subjects        []v1.Subject
}

type PuppyRole struct {
	Kind     string
	RoleName string
	Rules    []v1.PolicyRule
}

type PuppySubject interface {
	Matches(v1.Subject) bool
}

type PuppyUserSubject struct {
	Name string
}

func (ps PuppyUserSubject) Matches(subject v1.Subject) bool {
	return subject.Kind == "User" && subject.Name == ps.Name
}

func NewUser(username string) PuppySubject {
	return PuppyUserSubject{Name: username}
}

func InitialiseCluster(ctx context.Context) (*PuppyCluster, error) {

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

	serverCaData, err := os.ReadFile(serverCaPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read KUBEPUPPY_SERVER_CA_FILE file: '%v': %v", serverCaPath, err)
	}

	serverCaBlock, _ := pem.Decode(serverCaData)
	if serverCaBlock == nil {
		return nil, fmt.Errorf("invalid KUBEPUPPY_SERVER_CA_FILE: '%v'", serverCaPath)
	}

	serverCa, err := x509.ParseCertificates(serverCaBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse KUBEPUPPY_SERVER_CA_FILE: '%v': %v", serverCaPath, err)
	}

	pc := PuppyCluster{}
	pc.lock = &sync.RWMutex{}
	pc.kubeClient = clientset
	pc.ServerAuthority = serverCa

	err = pc.RefreshRolesAndBindings(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to refresh role bindings: %v", err)
	}

	return &pc, nil
}

func (pc *PuppyCluster) RefreshRolesAndBindings(ctx context.Context) error {
	pc.lock.Lock()
	defer pc.lock.Unlock()

	roleBindings := make(map[string]*PuppyRoleBinding)
	roles := make(map[string]*PuppyRole)

	clusterRoleBindings, err := (*pc.kubeClient).RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {
		roleBindings[clusterRoleBinding.Name] = &PuppyRoleBinding{
			Kind:            clusterRoleBinding.Kind,
			RoleBindingName: clusterRoleBinding.Name,
			Subjects:        clusterRoleBinding.Subjects,
			RoleName:        clusterRoleBinding.RoleRef.Name}
	}

	clusterRoles, err := (*pc.kubeClient).RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, clusterRole := range clusterRoles.Items {
		roles[clusterRole.Name] = &PuppyRole{
			Kind:     clusterRole.Kind,
			RoleName: clusterRole.Name,
			Rules:    clusterRole.Rules}
	}

	namespacedRoleBindings, err := (*pc.kubeClient).RbacV1().RoleBindings("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, namespacedRoleBinding := range namespacedRoleBindings.Items {
		fullRoleBindingName := fmt.Sprintf("%v/%v", namespacedRoleBinding.Namespace, namespacedRoleBinding.Name)
		roleBindings[fullRoleBindingName] = &PuppyRoleBinding{
			Kind:            namespacedRoleBinding.Kind,
			RoleBindingName: fullRoleBindingName,
			Subjects:        namespacedRoleBinding.Subjects,
			RoleName:        fmt.Sprintf("%v/%v", namespacedRoleBinding.Namespace, namespacedRoleBinding.RoleRef.Name)}
	}

	namespacedRoles, err := (*pc.kubeClient).RbacV1().Roles("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, namespacedRole := range namespacedRoles.Items {
		fullRoleName := fmt.Sprintf("%v/%v", namespacedRole.Namespace, namespacedRole.Name)
		roles[fullRoleName] = &PuppyRole{
			Kind:     namespacedRole.Kind,
			RoleName: fullRoleName,
			Rules:    namespacedRole.Rules}
	}

	pc.RoleBindings = roleBindings
	pc.Roles = roles

	return nil
}

func (pc *PuppyCluster) FindRoleBindingsForSubject(ps PuppySubject) ([]*PuppyRoleBinding, []*PuppyRole) {

	roleBindings := utils.Filter(utils.MapValues(pc.RoleBindings), func(roleBinding *PuppyRoleBinding) bool {
		return utils.Any(roleBinding.Subjects, func(subject v1.Subject) bool { return ps.Matches(subject) })
	})

	roles := utils.Map(roleBindings, func(roleBinding *PuppyRoleBinding) *PuppyRole {
		role, ok := pc.Roles[roleBinding.RoleName]
		if !ok {
			log.Printf("ROLE NOT FOUND: %v FOR BINDING %v", roleBinding.RoleName, roleBinding.RoleBindingName)
		}
		return role
	})

	return roleBindings, roles
}
