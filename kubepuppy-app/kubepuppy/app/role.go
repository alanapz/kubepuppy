package app

import (
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	"kubepuppy/utils"
)

type RoleKind = string

const ClusterRoleKind RoleKind = "ClusterRole"
const NamespacedRoleKind RoleKind = "Role"

type Role struct {
	RoleId       string
	Kind         RoleKind
	Name         string
	Namespace    string
	Rules        []v1.PolicyRule
	RoleBindings map[string]*RoleBinding
	Subjects     map[string]*Subject
}

func GenerateRoleIdFromRoleRef(roleRef v1.RoleRef, namespace string) (*string, error) {

	if roleRef.Kind == ClusterRoleKind {
		return utils.Ptr(GenerateClusterRoleId(roleRef.Name)), nil
	}

	if roleRef.Kind == NamespacedRoleKind {
		return utils.Ptr(GenerateNamespacedRoleId(roleRef.Name, namespace)), nil
	}

	return nil, fmt.Errorf("unexpected role kind: '%v'", roleRef.Kind)
}

func GenerateClusterRoleId(roleName string) string {
	return fmt.Sprintf("%v::%v", ClusterRoleKind, roleName)
}

func GenerateNamespacedRoleId(roleName string, namespace string) string {
	return fmt.Sprintf("%v::%v::%v", NamespacedRoleKind, namespace, roleName)
}
