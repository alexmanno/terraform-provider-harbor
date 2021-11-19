package client

import (
	"github.com/alexmanno/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UserGroupBody return a json body
func UserGroupBody(d *schema.ResourceData) models.UserGroupBody {
	return models.UserGroupBody{
		GroupName:   d.Get("group_name").(string),
		LdapGroupDn: d.Get("ldap_group_dn").(string),
	}
}
