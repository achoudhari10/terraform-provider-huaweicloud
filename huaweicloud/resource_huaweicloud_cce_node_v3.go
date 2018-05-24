package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/cce/v3/nodes"
)

func resourceCCENodeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCCENodeV3Create,
		Read:   resourceCCENodeV3Read,
		Update: resourceCCENodeV3Update,
		Delete: resourceCCENodeV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kind": &schema.Schema{
				Type:     schema.TypeString,
				Default: "Node",
				Optional: true,
			},
			"api_version": &schema.Schema{
				Type:     schema.TypeString,
				Default: "v3",
				Optional: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"annotations": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"flavor": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"az": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sshkey": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"root_volume": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				//ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"volumetype": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"extend_param": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					}},
			},
			"data_volumes": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				//ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							//ForceNew: true,
						},
						"volumetype": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							//ForceNew: true,
						},
						"extend_param": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							//ForceNew: true,
						},
					}},
			},
			"eip_ids": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"eip_count": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"iptype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"chargemode": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"sharetype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"billing_mode": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"node_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"extend_param": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCCENodeLabelsV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("labels").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceCCENodeAnnotationsV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("annotations").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceCCEDataVolume(d *schema.ResourceData) []nodes.VolumeSpec {
	log.Printf("[DEBUG] dataVolumes: %+v", d.Get("data_volumes"))
	volumeRaw := d.Get("data_volumes").(*schema.Set).List()
	log.Printf("[DEBUG] dataVolumessss: %+v", volumeRaw)
	volumes := make([]nodes.VolumeSpec , len(volumeRaw))
	for i, raw := range volumeRaw {
		rawMap := raw.(map[string]interface{})
		volumes[i] = nodes.VolumeSpec{
			Size:  rawMap["size"].(int),
			VolumeType: rawMap["volumetype"].(string),
			ExtendParam: rawMap["extend_param"].(string),
		}
	}
	return volumes
}
func resourceCCERootVolume(d *schema.ResourceData) nodes.VolumeSpec {
	var nics nodes.VolumeSpec
	nicsRaw := d.Get("root_volume").([]interface{})
	if len(nicsRaw) == 1 {
		nics.Size = nicsRaw[0].(map[string]interface{})["size"].(int)
		nics.VolumeType = nicsRaw[0].(map[string]interface{})["volumetype"].(string)
		nics.ExtendParam = nicsRaw[0].(map[string]interface{})["extend_param"].(string)
	}
	log.Printf("[DEBUG] nics: %+v", nics)
	return nics
}
func resourceCCENodeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.CCEv3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}
	idsraw := d.Get("eip_ids").([]interface{})
	eip_ids := make([]string, len(idsraw))
	for i, publicidraw := range idsraw {
		eip_ids[i] = publicidraw.(string)
	}

	createOpts := nodes.CreateOpts{
			Kind:       d.Get("kind").(string),
			ApiVersion:            d.Get("api_version").(string),
			Metadata:   nodes.CreateMetaData{
				Name: d.Get("name").(string),
				Labels: resourceCCENodeLabelsV2(d),
				Annotations: resourceCCENodeAnnotationsV2(d),
				},
			Spec: nodes.Spec{
				Flavor:     d.Get("flavor").(string),
				Az:             d.Get("az").(string),
				Login:          nodes.LoginSpec{SshKey: d.Get("sshkey").(string)},
				RootVolume: resourceCCERootVolume(d),
				DataVolumes: resourceCCEDataVolume(d),
				PublicIP:     nodes.PublicIPSpec{
					Ids: eip_ids,
					Count: d.Get("eip_count").(int),
					Eip: nodes.EipSpec{
						IpType: d.Get("iptype").(string),
						Bandwidth: nodes.BandwidthOpts{
							ChargeMode: d.Get("chargemode").(string),
							Size: d.Get("size").(int),
							ShareType: d.Get("sharetype").(string),
						},
					},
				},
				BillingMode:     d.Get("billing_mode").(int),
				Count:     d.Get("node_count").(int),
				ExtendParam:     d.Get("extend_param").(string),
			},
	}

	log.Printf("[DEBUG] Value of CreateOpts: %+v", createOpts)

	clusterid := d.Get("cluster_id").(string)
	log.Printf("[DEBUG] Value of clusterid:  %+v", clusterid)
	s, err := nodes.Create(nodeClient, clusterid, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE Node: %s", err)
	}

	log.Printf("[DEBUG] Waiting for CCE Node (%s) to become available", s.Metadata.Id)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"Available","Unavailable","Empty"},
		Refresh:    waitForCceNodeActive(nodeClient, clusterid,s.Metadata.Id),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(s.Metadata.Id)

	log.Printf("[DEBUG] Created Node %s: %#v", s.Metadata.Id, s)
	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.CCEv3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	s, err := nodes.Get(nodeClient, clusterid, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving HuaweiCloud Node: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Node %s: %#v", d.Id(), s)

	d.Set("kind", s.Kind)
	d.Set("api_version", s.Apiversion)
	d.Set("name", s.Metadata.Name)
	d.Set("id", s.Metadata.Id)
	d.Set("labels", s.Metadata.Labels)
	d.Set("annotations", s.Metadata.Annotations)
	d.Set("flavor", s.Spec.Flavor)
	d.Set("az", s.Spec.Az)
	d.Set("billing_mode", s.Spec.BillingMode)
	d.Set("node_count", s.Spec.Count)
	d.Set("extend_param", s.Spec.ExtendParam)
	d.Set("sshkey", s.Spec.Login.SshKey)
	d.Set("size", s.Spec.RootVolume.Size)
	d.Set("volumetype", s.Spec.RootVolume.VolumeType)
	d.Set("extend_param", s.Spec.RootVolume.ExtendParam)
	var volumes []map[string]interface{}
	for _, pairObject := range s.Spec.DataVolumes {
		volume := make(map[string]interface{})
		volume["size"] = pairObject.Size
		volume["volumetype"] = pairObject.VolumeType
		volume["extend_param"] = pairObject.ExtendParam
		volumes = append(volumes, volume)
	}
	if err := d.Set("data_volumes", volumes); err != nil {
		return fmt.Errorf("[DEBUG] Error saving dataVolumes to state for HuaweiCloud Node (%s): %s", d.Id(), err)
	}
	d.Set("extend_param", s.Spec.ExtendParam)
	d.Set("eip_ids", s.Spec.PublicIP.Ids)
	d.Set("eip_count", s.Spec.PublicIP.Count)
	d.Set("iptype", s.Spec.PublicIP.Eip.IpType)
	d.Set("chargemode", s.Spec.PublicIP.Eip.Bandwidth.ChargeMode)
	d.Set("size", s.Spec.PublicIP.Eip.Bandwidth.Size)
	d.Set("sharetype", s.Spec.PublicIP.Eip.Bandwidth.ShareType)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCCENodeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	nodeClient, err := config.CCEv3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	var updateOpts nodes.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Metadata.Name = d.Get("name").(string)
	}

	log.Printf("[DEBUG] Updating Node %s with options: %+v", d.Id(), updateOpts)

	clusterid := d.Get("cluster_id").(string)
	log.Printf("[DEBUG] clusterid: %+v",  clusterid)
	_, err = nodes.Update(nodeClient, clusterid, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud  Node: %s", err)
	}

	return resourceCCENodeV3Read(d, meta)
}

func resourceCCENodeV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	log.Printf("[DEBUG] config Value: %+v",  config)
	nodeClient, err := config.CCEv3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}
	clusterid := d.Get("cluster_id").(string)
	log.Printf("[DEBUG] clusterid: %+v",  clusterid)
	err = nodes.Delete(nodeClient, clusterid, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE Cluster: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Deleting"},
		Target:     []string{"DELETED"},
		Refresh:    waitForCceNodeDelete(nodeClient, clusterid,d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud CCE Node: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForCceNodeActive(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := nodes.Get(cceClient, clusterId, nodeId).Extract()
		if err != nil {
			return nil, "", err
		}
		if n.Status.Phase != "Creating"  {
			return n,"Creating", nil
		}

		return n, n.Status.Phase, nil
	}
}

func waitForCceNodeDelete(cceClient *golangsdk.ServiceClient, clusterId, nodeId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud CCE Node %s.\n", clusterId)

		r, err := nodes.Get(cceClient, clusterId, nodeId).Extract()

		log.Printf("[DEBUG] Value after extract: %#v", r)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud CCE Node %s", nodeId)
				return r, "Deleted", nil
			}
			return r, "Deleting", err
		}

		log.Printf("[DEBUG] HuaweiCloud CCE Node %s still available.\n", nodeId)
		return r, r.Status.Phase, nil
	}
}