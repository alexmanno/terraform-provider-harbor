package models

var PathUserGroups = "/usergroups"

//
type UserGroupBody struct {
  ID          int    `json:"id,omitempty"`
	GroupName   string `json:"group_name,omitempty"`
	LdapGroupDn string `json:"ldap_group_dn,omitempty"`
}
