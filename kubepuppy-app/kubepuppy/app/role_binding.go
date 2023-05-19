package app

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
)

type RoleBindingKind = string

const ClusterRoleBindingKind RoleBindingKind = "ClusterRoleBinding"
const NamespacedRoleBindingKind RoleBindingKind = "RoleBinding"

type RoleBinding struct {
	RoleBindingId string
	Kind          RoleBindingKind
	Name          string
	Namespace     string
	RoleRef       v1.RoleRef
	Role          *Role
	Subjects      map[string]*Subject
}

func GenerateClusterRoleBindingId(roleBindingName string) string {
	return fmt.Sprintf("%v::%v", ClusterRoleBindingKind, roleBindingName)
}

func GenerateNamespacedRoleBindingId(roleBindingName string, namespace string) string {
	return fmt.Sprintf("%v::%v::%v", NamespacedRoleBindingKind, roleBindingName, namespace)
}
