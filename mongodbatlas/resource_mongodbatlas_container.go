package mongodbatlas

import (
	"errors"
	"fmt"
	"log"
	"strings"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceContainerCreate,
		Read:   resourceContainerRead,
		Update: resourceContainerUpdate,
		Delete: resourceContainerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceContainerImportState,
		},

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"atlas_cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"provider_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"provisioned": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := ma.Container{
		AtlasCidrBlock: d.Get("atlas_cidr_block").(string),
		ProviderName:   d.Get("provider_name").(string),
		RegionName:     d.Get("region").(string),
	}

	container, _, err := client.Containers.Create(d.Get("group").(string), &params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB Container: %s", err)
	}
	d.SetId(container.ID)
	log.Printf("[INFO] MongoDB Container ID: %s", d.Id())

	return resourceContainerRead(d, meta)
}

func resourceContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	c, _, err := client.Containers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Container %s: %s", d.Id(), err)
	}

	d.Set("atlas_cidr_block", c.AtlasCidrBlock)
	d.Set("provider_name", c.ProviderName)
	d.Set("region", c.RegionName)
	d.Set("vpc_id", c.VpcID)
	d.Set("identifier", c.ID)
	d.Set("provisioned", c.Provisioned)

	return nil
}

func resourceContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	requestUpdate := false

	c, _, err := client.Containers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Container %s: %s", d.Id(), err)
	}

	if d.HasChange("atlas_cidr_block") {
		c.AtlasCidrBlock = d.Get("atlas_cidr_block").(string)
		requestUpdate = true
	}
	if d.HasChange("provider_name") {
		c.ProviderName = d.Get("provider_name").(string)
		requestUpdate = true
	}
	if d.HasChange("region") {
		c.RegionName = d.Get("region").(string)
		requestUpdate = true
	}

	if requestUpdate {
		_, _, err := client.Containers.Update(d.Get("group").(string), d.Id(), c)
		if err != nil {
			return fmt.Errorf("Error reading MongoDB Container %s: %s", d.Id(), err)
		}
	}

	return resourceContainerRead(d, meta)
}

func getContainer(client *ma.Client, gid string, cidr string) (*ma.Container, error) {
	containers, _, err := client.Containers.List(gid)
	if err != nil {
		return nil, fmt.Errorf("Couldn't import container %s in group %s, error: %s", cidr, gid, err.Error())
	}
	for i := range containers {
		if containers[i].AtlasCidrBlock == cidr {
			return &containers[i], nil
		}
	}
	return nil, fmt.Errorf("Couldn't find container with cidr %s in group %s", cidr, gid)

}

func resourceContainerImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*ma.Client)
	parts := strings.SplitN(d.Id(), "-", 2)
	if len(parts) != 2 {
		return nil, errors.New("To import a Container, use the format {group id}-{container cidr}")
	}
	gid := parts[0]
	cidr := parts[0]
	container, err := getContainer(client, gid, cidr)
	if err != nil {
		return nil, err
	}
	d.SetId(container.ID)
	d.Set("group", gid)
	return []*schema.ResourceData{d}, nil
}

func resourceContainerDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
