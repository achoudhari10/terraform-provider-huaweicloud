package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
	"log"
)

func dataSourceCCEClusterV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCCEClusterV3Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"billing_mode": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"highway_subnet_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_network_type": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"container_network_cidr": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"internal_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"external_endpoint": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCCEClusterV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))

	listOpts := clusters.ListOpts{
		ID:    d.Get("id").(string),
		Name:  d.Get("name").(string),
		Type:  d.Get("cluster_type").(string),
		Phase: d.Get("status").(string),
		VpcID: d.Get("vpc_id").(string),
	}

	refinedClusters, err := clusters.List(cceClient).ExtractCluster(listOpts)
	log.Printf("[DEBUG] Value of allClusters: %#v", refinedClusters)
	if err != nil {
		return fmt.Errorf("Unable to retrieve clusters: %s", err)
	}

	if len(refinedClusters) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedClusters) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Cluster := refinedClusters[0]

	log.Printf("[DEBUG] Retrieved Clusters using given filter %s: %+v", Cluster.Metadata.Id, Cluster)
	d.SetId(Cluster.Metadata.Id)

	d.Set("id", Cluster.Metadata.Id)
	d.Set("name", Cluster.Metadata.Name)
	d.Set("flavor", Cluster.Spec.Flavor)
	d.Set("description", Cluster.Spec.Description)
	d.Set("cluster_version", Cluster.Spec.Version)
	d.Set("cluster_type", Cluster.Spec.Type)
	d.Set("billing_mode", Cluster.Spec.BillingMode)
	d.Set("vpc_id", Cluster.Spec.HostNetwork.VpcId)
	d.Set("subnet_id", Cluster.Spec.HostNetwork.SubnetId)
	d.Set("highway_subnet_id", Cluster.Spec.HostNetwork.HighwaySubnet)
	d.Set("container_network_cidr", Cluster.Spec.ContainerNetwork.Cidr)
	d.Set("container_network_type", Cluster.Spec.ContainerNetwork.Mode)
	d.Set("status", Cluster.Status.Phase)
	d.Set("internal_endpoint", Cluster.Status.Endpoints.Internal)
	d.Set("external_endpoint", Cluster.Status.Endpoints.External)
	d.Set("region", GetRegion(d, config))

	return nil
}