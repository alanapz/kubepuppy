/*
 * KubePuppy API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type RoleBinding struct {

	RoleBindingId string `json:"roleBindingId"`

	Kind string `json:"kind"`

	Name string `json:"name"`

	// For a ServiceAccount, the account namespace. Null for Users and Groups
	Namespace string `json:"namespace,omitempty"`
}