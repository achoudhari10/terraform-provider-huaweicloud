package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
	"log"
)

//Creates schema for data source
func dataSourceCceNodesV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCceNodesV3Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"apiversion": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"node_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"az": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"charge_mode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"share_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"disk_size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"volume_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"extend_param": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_volumes": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disk_size": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
						},
						"volume_type": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"billing_mode": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"phase": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"server_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"public_ip_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"spec_extend_param": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"spec_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"job_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"reason": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"message": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"publicip_id_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

//filters the nodes and sets the redifined searched nodes
func dataSourceCceNodesV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	cceClient, err := config.cceV3Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Unable to get Nodes: %s", err)
	}

	listOpts := nodes.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("node_id"); ok {
		listOpts.Uid = v.(string)
	}

	if v, ok := d.GetOk("phase"); ok {
		listOpts.Phase = v.(string)
	}

	refinedNodes, err := nodes.List(cceClient, d.Get("cluster_id").(string)).ExtractNode(listOpts)
	log.Printf("[DEBUG] Value of all Nodes: %#v", refinedNodes)

	if err != nil {
		return fmt.Errorf("Unable to retrieve Nodes: %s", err)
	}

	if len(refinedNodes) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNodes) > 1 {
		return fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	Node := refinedNodes[0]

	var v []map[string]interface{}
	for _, volume := range Node.Spec.DataVolumes {

		mapping := map[string]interface{}{
			"disk_size":   volume.Size,
			"volume_type": volume.Volumetype,
		}
		v = append(v, mapping)
	}



	pids := Node.Spec.PublicIP.Ids
	PublicIDs := make([]string, len(pids))
	for i, val := range pids {
		PublicIDs[i] = val
	}
	log.Printf("[DEBUG] Retrieved Clusters using given filter %s: %+v", Node.Metadata.Uid, Node)
	d.SetId(Node.Metadata.Uid)
	d.Set("kind", Node.Kind)
	d.Set("node_id", Node.Metadata.Uid)
	d.Set("apiversion", Node.Apiversion)
	d.Set("name", Node.Metadata.Name)
	d.Set("flavor", Node.Spec.Flavor)
	d.Set("az", Node.Spec.Az)
	d.Set("billing_mode", Node.Spec.BillingMode)
	d.Set("phase", Node.Status.Phase)
	d.Set("data_volumes", v)
	d.Set("disk_size", Node.Spec.RootVolume.Size)
	d.Set("volume_type", Node.Spec.RootVolume.Volumetype)
	d.Set("extend_param", Node.Spec.RootVolume.ExtendParam)
	d.Set("sshkey", Node.Spec.Login.SshKey)
	d.Set("charge_mode", Node.Spec.PublicIP.Eip.Bandwidth.Chargemode)
	d.Set("size", Node.Spec.PublicIP.Eip.Bandwidth.Size)
	d.Set("share_type", Node.Spec.PublicIP.Eip.Bandwidth.Sharetype)
	d.Set("ip_type", Node.Spec.PublicIP.Eip.Iptype)
	d.Set("server_id", Node.Status.ServerID)
	d.Set("public_ip", Node.Status.PublicIP)
	d.Set("private_ip", Node.Status.PrivateIP)
	d.Set("spec_extend_param", Node.Spec.ExtendParam)
	d.Set("spec_count", Node.Spec.Count)
	d.Set("job_id", Node.Status.JobID)
	d.Set("reason", Node.Status.Reason)
	d.Set("message", Node.Status.Message)
	d.Set("publicip_id_count", Node.Spec.PublicIP.Count)
	d.Set("public_ip_ids", PublicIDs)

	return nil
}
