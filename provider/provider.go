package provider

import (
	"strings"

	"github.com/alexmanno/terraform-provider-harbor/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_URL", ""),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_USERNAME", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_PASSWORD", ""),
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				DefaultFunc: schema.EnvDefaultFunc("HARBOR_IGNORE_CERT", ""),
			},
			"api_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  2,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"harbor_user_group":             resourceUserGroup(),
		},
		DataSourcesMap: map[string]*schema.Resource{
		},

		ConfigureFunc: providerConfigure,
	}
}

func checkProjectid(id string) (projecid string) {
	path := "/projects/"
	if strings.Contains(id, path) == false {
		id = path + id
	}
	return id

}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	var apiPath string

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	insecure := d.Get("insecure").(bool)
	apiVersion := d.Get("api_version").(int)

	if strings.HasSuffix(url, "/") {
		url = strings.Trim(url, "/")
	}

	if apiVersion == 1 {
		apiPath = "/api"
	} else if apiVersion == 2 {
		apiPath = "/api/v2.0"
	}

	return client.NewClient(url+apiPath, username, password, insecure), nil
}
