/*
 * KubePuppy API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteRoleBinding - Deletes the specified role binding.
func DeleteRoleBinding(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// GetRoleBinding - Shows details regarding the specified role binding (either cluster or namespaced).
func GetRoleBinding(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// ListRoleBindings - Lists all role bindings (both cluster and namespaced).
func ListRoleBindings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}