package openapi

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"kubepuppy/app"
	"kubepuppy/utils"
	"net/http"
)

func GetRole(c *gin.Context) {
	roleId := c.Param("roleId")

	if len(roleId) == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("roleId required"))
		return
	}

	cluster := GetClusterFromRequest(c)

	cluster.LockForRead(func() {

		role, ok := cluster.Roles[roleId]

		if !ok {
			_ = c.AbortWithError(http.StatusBadRequest, fmt.Errorf("role not found: %v", roleId))
			return
		}

		c.JSON(http.StatusOK, RoleDetails{
			RoleId:    role.RoleId,
			Kind:      role.Kind,
			Name:      role.Name,
			Namespace: role.Namespace,
		})
	})
}

func ListRoles(c *gin.Context) {
	cluster := GetClusterFromRequest(c)

	cluster.LockForRead(func() {

		c.JSON(http.StatusOK, utils.Map(utils.MapValues(cluster.Roles), func(role *app.Role) Role {
			return Role{
				RoleId:    role.RoleId,
				Kind:      role.Kind,
				Name:      role.Name,
				Namespace: role.Namespace,
			}
		}))
	})
}

func DeleteRole(c *gin.Context) {
	roleId := c.Param("roleId")

	if len(roleId) == 0 {
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("roleId required"))
		return
	}

	cluster := GetClusterFromRequest(c)

	cluster.LockForWrite(func() {
		panic("!!! TODO !!!")
	})
}

func roleToApiObject(role *app.Role) ApiObject {
	if role == nil {
		return ApiObject{}
	}
	return ApiObject{
		Id:        role.RoleId,
		Kind:      role.Kind,
		Name:      role.Name,
		Namespace: role.Namespace,
	}
}
