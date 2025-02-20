/*
Nexodus API

This is the Nexodus API Server.

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package public

// ModelsUser struct for ModelsUser
type ModelsUser struct {
	CreatedAt string `json:"createdAt,omitempty"`
	// Since the ID comes from the IDP, we have no control over the format...
	Id              string `json:"id,omitempty"`
	SecurityGroupId string `json:"security_group_id,omitempty"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
	UserName        string `json:"userName,omitempty"`
}
