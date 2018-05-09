package huaweicloud
import (
	"fmt"
	"log"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stackresources"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRtsStackResourcesV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRtsStackResourcesV1Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"stack_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"stack_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"resource_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"logical_resource_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"links": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"rel": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"required_by":&schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"resource_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_status_reason": &schema.Schema{
				Type:     schema.TypeString,
				Optional:true,
			},
			"physical_resource_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional:true,
			},
			"resource_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceRtsStackResourcesV1Read(d *schema.ResourceData,  meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))

	listOpts := stackresources.ListOpts{
		Name:       d.Get("resource_name").(string),
		LogicalID:  d.Get("logical_resource_id").(string),
		PhysicalID: d.Get("physical_resource_id").(string),
		Status:     d.Get("resource_status").(string),
		Type:       d.Get("resource_type").(string),
	}

	refinedResources ,err := stackresources.List(orchestrationClient,d.Get("stack_name").(string),d.Get("stack_id").(string),listOpts)
	log.Printf("[DEBUG] Value of allResources: %#v", refinedResources)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Resources: %s", err)
	}

	if  refinedResources == nil ||len(refinedResources) == 0 {
		return fmt.Errorf("No matching resource found. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedResources) > 1 {
		return fmt.Errorf("multiple resources matched; use additional constraints to reduce matches to a single resource")
	}

	resources := refinedResources[0]

	log.Printf("[DEBUG] Retrieved Resources using given filter %s: %+v", resources.Name,resources)
	d.SetId(resources.Name)

	var s []map[string]interface{}
	for _, stack := range resources.Links {
		mapping := map[string]interface{}{
			"href": stack.Href,
			"rel":  stack.Rel,
		}
		s = append(s, mapping)
	}

	d.Set("resource_name", resources.Name)
	d.Set("resource_status", resources.Status)
	d.Set("logical_resource_id", resources.LogicalID)
	d.Set("physical_resource_id", resources.PhysicalID)
	d.Set("required_by", resources.RequiredBy)
	d.Set("resource_status_reason", resources.StatusReason)
	d.Set("resource_type", resources.Type)
	d.Set("region", GetRegion(d, config))
	if err := d.Set("links", s); err != nil {
		return err
	}
	return nil
}
