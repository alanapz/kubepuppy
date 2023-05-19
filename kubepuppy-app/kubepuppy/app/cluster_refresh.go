package app

import (
	"context"
	"fmt"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubepuppy/utils"
)

func (c *Cluster) RefreshRolesAndBindings(ctx context.Context) error {

	if !utils.IsRWMutexWriteLocked(&c.lock) {
		panic("cannot update without holding write lock")
	}

	roles := make(map[string]*Role)
	roleBindings := make(map[string]*RoleBinding)
	subjects := make(map[string]*Subject)

	err := c.refreshClusterRoles(ctx, roles)

	if err != nil {
		return err
	}

	err = c.refreshNamespacedRoles(ctx, roles)

	if err != nil {
		return err
	}

	err = c.refreshClusterRoleBindings(ctx, roles, roleBindings, subjects)

	if err != nil {
		return err
	}

	err = c.refreshNamespacedRoleBindings(ctx, roles, roleBindings, subjects)

	if err != nil {
		return err
	}

	// We need to handle system:authenticated and system:unauthenticated virtual groupos

	c.Roles = roles
	c.RoleBindings = roleBindings
	c.Subjects = subjects

	return nil
}

func (c *Cluster) refreshClusterRoles(ctx context.Context, roles map[string]*Role) error {

	clusterRoles, err := (*c.KubeClient).RbacV1().ClusterRoles().List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, clusterRole := range clusterRoles.Items {

		if clusterRole.Namespace != "" {
			return fmt.Errorf("unexpected namespace: '%v' for cluster role: '%v'", clusterRole.Namespace, clusterRole.Name)
		}

		err = c.attachRole(
			roles,
			GenerateClusterRoleId(clusterRole.Name),
			ClusterRoleKind,
			clusterRole.Name,
			clusterRole.Namespace,
			clusterRole.Rules)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) refreshNamespacedRoles(ctx context.Context, roles map[string]*Role) error {

	namespacedRoles, err := (*c.KubeClient).RbacV1().Roles("").List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, namespacedRole := range namespacedRoles.Items {

		if namespacedRole.Namespace == "" {
			return fmt.Errorf("unexpected namespace: '%v' for role: '%v'", namespacedRole.Namespace, namespacedRole.Name)
		}

		err = c.attachRole(
			roles,
			GenerateNamespacedRoleId(namespacedRole.Name, namespacedRole.Namespace),
			NamespacedRoleKind,
			namespacedRole.Name,
			namespacedRole.Namespace,
			namespacedRole.Rules)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) attachRole(roles map[string]*Role, roleId string, roleKind string, roleName string, roleNamespace string, rules []v1.PolicyRule) error {

	role := &Role{
		RoleId:       roleId,
		Kind:         roleKind,
		Name:         roleName,
		Namespace:    roleNamespace,
		Rules:        rules,
		RoleBindings: make(map[string]*RoleBinding),
		Subjects:     make(map[string]*Subject),
	}

	roles[role.RoleId] = role

	return nil
}

func (c *Cluster) refreshClusterRoleBindings(ctx context.Context, roles map[string]*Role, roleBindings map[string]*RoleBinding, subjects map[string]*Subject) error {

	clusterRoleBindings, err := (*c.KubeClient).RbacV1().ClusterRoleBindings().List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, clusterRoleBinding := range clusterRoleBindings.Items {

		if clusterRoleBinding.Namespace != "" {
			return fmt.Errorf("unexpected namespace: '%v' for cluster role binding: '%v'", clusterRoleBinding.Namespace, clusterRoleBinding.Name)
		}

		err = c.attachRoleBinding(
			roles,
			roleBindings,
			subjects,
			GenerateClusterRoleBindingId(clusterRoleBinding.Name),
			ClusterRoleBindingKind,
			clusterRoleBinding.Name,
			clusterRoleBinding.Namespace,
			clusterRoleBinding.RoleRef,
			clusterRoleBinding.Subjects)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) refreshNamespacedRoleBindings(ctx context.Context, roles map[string]*Role, roleBindings map[string]*RoleBinding, subjects map[string]*Subject) error {

	namespacedRoleBindings, err := (*c.KubeClient).RbacV1().RoleBindings("").List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for _, namespacedRoleBinding := range namespacedRoleBindings.Items {

		if namespacedRoleBinding.Namespace == "" {
			return fmt.Errorf("unexpected namespace: '%v' for role binding: '%v'", namespacedRoleBinding.Namespace, namespacedRoleBinding.Name)
		}

		err = c.attachRoleBinding(
			roles,
			roleBindings,
			subjects,
			GenerateNamespacedRoleBindingId(namespacedRoleBinding.Name, namespacedRoleBinding.Namespace),
			NamespacedRoleBindingKind,
			namespacedRoleBinding.Name,
			namespacedRoleBinding.Namespace,
			namespacedRoleBinding.RoleRef,
			namespacedRoleBinding.Subjects)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Cluster) attachRoleBinding(roles map[string]*Role, roleBindings map[string]*RoleBinding, subjects map[string]*Subject, roleBindingId string, roleBindingKind string, roleBindingName string, roleBindingNamespace string, roleRef v1.RoleRef, roleSubjects []v1.Subject) error {

	roleId, err := GenerateRoleIdFromRoleRef(roleRef, roleBindingNamespace)

	if err != nil {
		return err
	}

	role := roles[*roleId]

	roleBinding := &RoleBinding{
		RoleBindingId: roleBindingId,
		Kind:          roleBindingKind,
		Name:          roleBindingName,
		Namespace:     roleBindingNamespace,
		RoleRef:       roleRef,
		Role:          role,
		Subjects:      make(map[string]*Subject),
	}

	roleBindings[roleBinding.RoleBindingId] = roleBinding

	if role != nil {
		role.RoleBindings[roleBinding.RoleBindingId] = roleBinding
	}

	for _, roleSubject := range roleSubjects {

		subjectId, err := GenerateSubjectId(roleSubject)

		if err != nil {
			return fmt.Errorf("couldn't build subject for role binding: '%v': %v", roleBindingId, err)
		}

		subject := subjects[*subjectId]

		if subject == nil {
			subject = &Subject{
				SubjectId:    *subjectId,
				Kind:         roleSubject.Kind,
				Name:         roleSubject.Name,
				Namespace:    roleSubject.Namespace,
				Roles:        make(map[string]*Role),
				RoleBindings: make(map[string]*RoleBinding),
			}
			subjects[*subjectId] = subject
		}

		roleBinding.Subjects[*subjectId] = subject
		subject.RoleBindings[roleBinding.RoleBindingId] = roleBinding

		if role != nil {
			role.Subjects[*subjectId] = subject
			subject.Roles[role.RoleId] = role
		}
	}

	return nil
}
