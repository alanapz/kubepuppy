/*
 * KubePuppy API
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package openapi

type SubjectDetailsRoleBinding struct {

	RoleBinding ApiObject `json:"roleBinding"`

	Role ApiObject `json:"role,omitempty"`

	Rules []interface{} `json:"rules,omitempty"`
}