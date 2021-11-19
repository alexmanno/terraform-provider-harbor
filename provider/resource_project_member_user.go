package provider

import (
	"encoding/json"
	"fmt"

	"github.com/alexmanno/terraform-provider-harbor/client"
	"github.com/alexmanno/terraform-provider-harbor/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMembersUser() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"member_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if v != "projectadmin" && v != "developer" && v != "guest" && v != "master" && v != "limitedguest" {
						errs = append(errs, fmt.Errorf("%q must be either projectadmin, developer, guest, limitedguest or master, got: %s", key, v))
					}
					return
				},
			},
		},
		Create: resourceMembersUserCreate,
		Read:   resourceMembersUserRead,
		Update: resourceMembersUserUpdate,
		Delete: resourceMembersUserDelete,
	}
}

func resourceMembersUserCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)
	path := d.Get("project_id").(string) + "/members"

	body := client.ProjectMembersUserBody(d)

	_, headers, err := apiClient.SendRequest("POST", path, body, 201)
	if err != nil {
		return err
	}

	id, err := client.GetID(headers)
	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceMembersUserRead(d, m)
}

func resourceMembersUserRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	resp, _, err := apiClient.SendRequest("GET", d.Id(), nil, 200)
	if err != nil {
		d.SetId("")
		return nil
	}

	var jsonData models.ProjectMembersBody
	err = json.Unmarshal([]byte(resp), &jsonData)
	if err != nil {
		return fmt.Errorf("Resource not found %s", d.Id())
	}

	d.Set("role", client.RoleTypeNumber(jsonData.RoleID))
	return nil
}

func resourceMembersUserUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	body := client.ProjectMembersUserBody(d)
	_, _, err := apiClient.SendRequest("PUT", d.Id(), body, 200)
	if err != nil {
		fmt.Println(err)
	}

	return resourceMembersUserRead(d, m)
}

func resourceMembersUserDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	_, _, err := apiClient.SendRequest("DELETE", d.Id(), nil, 200)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
