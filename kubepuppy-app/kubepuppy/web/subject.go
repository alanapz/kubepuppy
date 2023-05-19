package openapi

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/rbac/v1"
	"kubepuppy/app"
	"kubepuppy/utils"
	"net/http"
)

func GetSubject(c *gin.Context) {
	subjectId := c.Param("subjectId")

	if len(subjectId) == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("subjectId required"))
		return
	}

	cluster := GetClusterFromRequest(c)

	cluster.LockForRead(func() {

		subject, ok := cluster.Subjects[subjectId]

		if !ok {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("subject not found: %v", subjectId))
			return
		}

		roleBindings := utils.Map(utils.MapValues(subject.RoleBindings), func(roleBinding *app.RoleBinding) SubjectDetailsRoleBinding {

			var rules []interface{}

			if roleBinding.Role != nil {
				rules = utils.Map(roleBinding.Role.Rules, func(rule v1.PolicyRule) interface{} {
					return rule
				})
			}

			return SubjectDetailsRoleBinding{
				RoleBinding: roleBindingToApiObject(roleBinding),
				Role:        roleToApiObject(roleBinding.Role),
				Rules:       rules,
			}
		})

		c.JSON(http.StatusOK, SubjectDetails{
			Id:           subject.SubjectId,
			Kind:         subject.Kind,
			Name:         subject.Name,
			Namespace:    subject.Namespace,
			RoleBindings: roleBindings,
		})
	})
}

func ListSubjects(c *gin.Context) {
	cluster := GetClusterFromRequest(c)

	cluster.LockForRead(func() {

		c.JSON(http.StatusOK, utils.Map(utils.MapValues(cluster.Subjects), func(subject *app.Subject) Subject {

			roleBindings := utils.Map(utils.MapValues(subject.RoleBindings), func(roleBinding *app.RoleBinding) SubjectRoleBinding {
				return SubjectRoleBinding{
					RoleBinding: roleBindingToApiObject(roleBinding),
					Role:        roleToApiObject(roleBinding.Role),
				}
			})

			return Subject{
				Id:           subject.SubjectId,
				Kind:         subject.Kind,
				Name:         subject.Name,
				Namespace:    subject.Namespace,
				RoleBindings: roleBindings,
			}
		}))
	})
}
