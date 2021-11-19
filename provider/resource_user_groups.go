package provider

import (
	"encoding/json"
	"fmt"

	"github.com/alexmanno/terraform-provider-harbor/client"
	"github.com/alexmanno/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUserGroup() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ldap_group_dn": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
		},
		Create: resourceUserGroupCreate,
		Read:   resourceUserGroupRead,
		Update: resourceUserGroupUpdate,
		Delete: resourceUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceUserGroupCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.UserBody(d)

	_, header, err := apiClient.SendRequest("POST", models.PathUserGroups, &body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(header)
	if err != nil {
		return nil
	}

	d.SetId(id)
	return resourceUserGroupRead(d, m)
}

func resourceUserGroupRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	resp, _, err := apiClient.SendRequest("GET", models.PathUserGroups+"/"+d.Id(), nil, 200)
	if err != nil {
		return err
	}
	var jsonData models.UserGroupBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	d.Set("group_name", jsonData.GroupName)
	d.Set("ldap_group_dn", jsonData.LdapGroupDn)

	return nil
}

func resourceUserGroupUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.UserBody(d)
	_, _, err := apiClient.SendRequest("PUT", models.PathUserGroups+"/"+d.Id(), body, 200)
	if err != nil {
		return err
	}

	return resourceUserGroupRead(d, m)
}

func resourceUserGroupDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, err := apiClient.SendRequest("DELETE", models.PathUserGroups+"/"+d.Id(), nil, 200)
	if err != nil {
		return err
	}
	return nil
}
