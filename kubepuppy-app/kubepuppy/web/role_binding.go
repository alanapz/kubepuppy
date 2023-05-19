package openapi

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"kubepuppy/app"
	"kubepuppy/utils"
	"net/http"
)

func GetRoleBinding(c *gin.Context) {
	roleBindingId := c.Param("roleBindingId")

	if len(roleBindingId) == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("roleBindingId required"))
		return
	}

	cluster := GetClusterFromRequest(c)

	cluster.LockForRead(func() {

		roleBinding, ok := cluster.RoleBindings[roleBindingId]

		if !ok {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("roleBinding not found: %v", roleBindingId))
			return
		}

		c.JSON(http.StatusOK, RoleBindingDetails{
			RoleBindingId: roleBinding.RoleBindingId,
			Kind:          roleBinding.Kind,
			Name:          roleBinding.Name,
			Namespace:     roleBinding.Namespace,
		})
	})
}

func ListRoleBindings(c *gin.Context) {
	cluster := GetClusterFromRequest(c)

	cluster.LockForRead(func() {

		c.JSON(http.StatusOK, utils.Map(utils.MapValues(cluster.RoleBindings), func(roleBinding *app.RoleBinding) RoleBinding {
			return RoleBinding{
				RoleBindingId: roleBinding.RoleBindingId,
				Kind:          roleBinding.Kind,
				Name:          roleBinding.Name,
				Namespace:     roleBinding.Namespace,
			}
		}))
	})
}

func DeleteRoleBinding(c *gin.Context) {
	roleBindingId := c.Param("roleBindingId")

	if len(roleBindingId) == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("roleBindingId required"))
		return
	}

	cluster := GetClusterFromRequest(c)

	cluster.LockForWrite(func() {
		panic("!!! TODO !!!")
	})
}

func roleBindingToApiObject(roleBinding *app.RoleBinding) ApiObject {
	if roleBinding == nil {
		return ApiObject{}
	}
	return ApiObject{
		Id:        roleBinding.RoleBindingId,
		Kind:      roleBinding.Kind,
		Name:      roleBinding.Name,
		Namespace: roleBinding.Namespace,
	}
}
