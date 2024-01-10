package netbox

import (
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSourceNetboxFHRPGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSourceNetboxFHRPGroupCreate,
		Read:   resourceSourceNetboxFHRPGroupRead,
		Update: resourceSourceNetboxFHRPGroupUpdate,
		Delete: resourceSourceNetboxFHRPGroupDelete,

		Description: `:meta:subcategory:IP Address Management (IPAM):From the [official documentation](https://docs.netbox.dev/en/stable/models/ipam/fhrpgroup/):

> The FHRPGroup model represents an FHRP group. first-hop redundancy protocol (FHRP) enables multiple physical interfaces to present a virtual IP address (VIP) in a redundant manner.`,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSourceNetboxFHRPGroupCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	data := models.FHRPGroup{}

	name := d.Get("name").(string)
	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to generated slug if not given
	if !slugOk {
		slug = getSlug(name)
	} else {
		slug = slugValue.(string)
	}

	data.Name = name
	data.Slug = slug
	data.Description = getOptionalStr(d, "description", true)

	_, err := api.IPAM.CreateFHRPGroup(&data)
	if err != nil {
		return err
	}

	d.SetId(strconv.FormatInt(data.ID, 10))

	return resourceSourceNetboxFHRPGroupRead(d, m)
}

func resourceSourceNetboxFHRPGroupRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	res, err := api.IPAM.GetFHRPGroupByID(id)
	if err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(res.ID, 10))
	d.Set("name", res.Name)
	d.Set("slug", res.Slug)
	d.Set("description", res.Description)
	return nil
}

func resourceSourceNetboxFHRPGroupUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.FHRPGroup{}
	data.ID = id
	data.Name = d.Get("name").(string)
	data.Slug = d.Get("slug").(string)
	data.Description = getOptionalStr(d, "description", true)
	_, err := api.IPAM.UpdateFHRPGroupByID(&data)
	if err != nil {
		return err
	}
	return resourceSourceNetboxFHRPGroupRead(d, m)
}

func resourceSourceNetboxFHRPGroupDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	_, err := api.IPAM.DeleteFHRPGroupByID(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}
