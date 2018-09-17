package mongodbatlas

import (
	"fmt"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceContainer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceContainerRead,

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"atlas_cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"provider_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioned": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func getContainerByIDOrCidr(client *ma.Client, gid string, cidr string, id string) (*ma.Container, error) {
	containers, _, err := client.Containers.List(gid)
	if err != nil {
		return nil, fmt.Errorf("Couldn't import container %s in group %s, error: %s", cidr, gid, err.Error())
	}
	for i := range containers {
		if containers[i].AtlasCidrBlock == cidr {
			return &containers[i], nil
		}
	}
	for i := range containers {
		if containers[i].ID == id {
			return &containers[i], nil
		}
	}
	return nil, fmt.Errorf("Couldn't find container with cidr %s or id %s in group %s", cidr, id, gid)

}

func dataSourceContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	id := d.Get("identifier").(string)
	group := d.Get("group").(string)
	cidr := d.Get("atlas_cidr_block").(string)

	c, err := getContainerByIDOrCidr(client, group, cidr, id)
	if err != nil {
		return err
	}

	d.SetId(c.ID)
	d.Set("atlas_cidr_block", c.AtlasCidrBlock)
	d.Set("provider_name", c.ProviderName)
	d.Set("region", c.RegionName)
	d.Set("vpc_id", c.VpcID)
	d.Set("identifier", c.ID)
	d.Set("provisioned", c.Provisioned)

	return nil
}
