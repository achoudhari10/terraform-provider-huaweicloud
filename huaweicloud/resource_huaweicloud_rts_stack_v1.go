package huaweicloud

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stacks"
	"log"
	"time"
)

func resourceRTSStackV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceRTSStackV1Create,
		Read:   resourceRTSStackV1Read,
		Update: resourceRTSStackV1Update,
		Delete: resourceRTSStackV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_body": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStackTemplate,
				StateFunc: func(v interface{}) string {
					template, _ := normalizeStackTemplate(v)
					return template
				},
			},
			"template_url": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"files": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"environment": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateJsonString,
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"timeout_mins": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"disable_rollback": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status_reason": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"outputs": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
			"capabilities": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"notification_topics": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceTemplateOptsV1(d *schema.ResourceData) *stacks.Template {
	var template = new(stacks.Template)
	if _, ok := d.GetOk("template_body"); ok {
		rawTemplate := d.Get("template_body").(string)
		template.Bin = []byte(rawTemplate)
	}
	if _, ok := d.GetOk("template_url"); ok {
		rawTemplateUrl := d.Get("template_url").(string)
		template.URL = rawTemplateUrl

	}
	if _, ok := d.GetOk("files"); ok {
		rawFiles := make(map[string]string)
		for key, val := range d.Get("files").(map[string]interface{}) {
			rawFiles[key] = val.(string)
		}
		template.Files = rawFiles

	}
	return template
}

func resourceEnvironmentV1(d *schema.ResourceData) *stacks.Environment {
	rawTemplate := d.Get("environment").(string)
	environment := new(stacks.Environment)
	environment.Bin = []byte(rawTemplate)
	return environment
}
func resourceParametersV1(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("parameters").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
func resourceRTSStackV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))

	log.Printf("[DEBUG] Value of orchestration client: %#v", orchestrationClient)

	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud orchestration client: %s", err)
	}

	rollback := d.Get("disable_rollback").(bool)
	createOpts := stacks.CreateOpts{
		Name:            d.Get("name").(string),
		TemplateOpts:    resourceTemplateOptsV1(d),
		DisableRollback: &rollback,
		EnvironmentOpts: resourceEnvironmentV1(d),
		Parameters:      resourceParametersV1(d),
		Timeout:         d.Get("timeout_mins").(int),
	}

	n, err := stacks.Create(orchestrationClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud stack: %s", err)
	}
	d.SetId(n.ID)

	log.Printf("[INFO] stack ID: %s", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATE_IN_PROGRESS"},
		Target:     []string{"CREATE_COMPLETE"},
		Refresh:    waitForRTSStackActive(orchestrationClient, d.Get("name").(string), n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()

	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for Stack (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	return resourceRTSStackV1Read(d, meta)

}

func resourceRTSStackV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud orchestration Client: %s", err)
	}

	n, err := stacks.Get(orchestrationClient, d.Get("name").(string), d.Id()).Extract()
	if err != nil {
		d.SetId("")
		return fmt.Errorf("Error retrieving HuaweiCloud Stacks: %s", err)

	}

	log.Printf("[DEBUG] Retrieved Stack %s: %+v", d.Id(), n)

	d.Set("disable_rollback", n.DisableRollback)
	d.Set("description", n.Description)
	d.Set("parameters", n.Parameters)
	d.Set("status_reason", n.StatusReason)
	d.Set("name", n.Name)
	d.Set("outputs", flattenStackOutputs(n.Outputs))
	d.Set("capabilities", n.Capabilities)
	d.Set("notification_topics", n.NotificationTopics)
	d.Set("timeout_mins", n.Timeout)
	d.Set("status", n.Status)
	d.Set("id", n.ID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceRTSStackV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud orchestration Client: %s", err)
	}

	var updateOpts stacks.UpdateOpts

	updateOpts.TemplateOpts = resourceTemplateOptsV1(d)

	if d.HasChange("environment") {

		updateOpts.EnvironmentOpts = resourceEnvironmentV1(d)
	}
	if d.HasChange("parameters") {

		updateOpts.Parameters = resourceParametersV1(d)
	}
	if d.HasChange("timeout_mins") {

		updateOpts.Timeout = d.Get("timeout_mins").(int)

	}
	if d.HasChange("disable_rollback") {

		rollback := d.Get("disable_rollback").(bool)
		updateOpts.DisableRollback = &rollback

	}

	err = stacks.Update(orchestrationClient, d.Get("name").(string), d.Id(), updateOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error updating HuaweiCloud Stack: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending: []string{"UPDATE_IN_PROGRESS",
			"CREATE_COMPLETE",
			"ROLLBACK_IN_PROGRESS"},
		Target:     []string{"UPDATE_COMPLETE"},
		Refresh:    waitForRTSStackUpdate(orchestrationClient, d.Get("name").(string), d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()

	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for updating stack: %s", stateErr)
	}

	return resourceRTSStackV1Read(d, meta)
}

func resourceRTSStackV1Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy Stack: %s", d.Id())

	config := meta.(*Config)
	orchestrationClient, err := config.orchestrationV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud orchestration Client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"DELETE_IN_PROGRESS",
			"CREATE_COMPLETE",
			"CREATE_FAILED",
			"UPDATE_COMPLETE",
			"UPDATE_FAILED",
			"CREATE_FAILED",
			"ROLLBACK_COMPLETE",
			"ROLLBACK_IN_PROGRESS"},
		Target: []string{"DELETE_COMPLETE",
			"DELETE_FAILED"},
		Refresh:    waitForRTSStackDelete(orchestrationClient, d.Get("name").(string), d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	out, err := stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud Stack: %s", err)
	}

	stack := out.(*stacks.RetrievedStack)

	if stack.Status == "DELETE_FAILED" {
		return fmt.Errorf("%s: %q", stack.Status, stack.StatusReason)

	}

	d.SetId("")
	return nil
}

func waitForRTSStackActive(orchestrationClient *golangsdk.ServiceClient, stackName string, stackId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := stacks.Get(orchestrationClient, stackName, stackId).Extract()
		if err != nil {
			return nil, "", err
		}
		log.Printf("[INFO] HuaweiCloud stack: %+v", n)
		if n.Status == "CREATE_IN_PROGRESS" {
			return n, n.Status, nil
		}

		if n.Status == "CREATE_FAILED" {
			return nil, "", fmt.Errorf("%s: %s", n.Status, n.StatusReason)
		}
		return n, n.Status, nil
	}
}

func waitForRTSStackDelete(orchestrationClient *golangsdk.ServiceClient, stackName string, stackId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud Stack %s.\n", stackId)
		r, err := stacks.Get(orchestrationClient, stackName, stackId).Extract()

		if r.Status == "DELETE_COMPLETE" {
			log.Printf("[INFO] Successfully deleted HuaweiCloud stack %s", r.ID)
			return r, "DELETE_COMPLETE", nil
		}

		err = stacks.Delete(orchestrationClient, stackName, stackId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud Stack %s", stackId)
				return r, r.Status, nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, r.Status, nil
				}
			}

			return r, r.Status, err
		}

		log.Printf("[DEBUG] HuaweiCloud Stack %s still active.\n", stackId)
		return r, r.Status, nil
	}
}

func waitForRTSStackUpdate(orchestrationClient *golangsdk.ServiceClient, stackName string, stackId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := stacks.Get(orchestrationClient, stackName, stackId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud stack: %+v", n)
		if n.Status == "UPDATE_IN_PROGRESS" {
			return n, "UPDATE_IN_PROGRESS", nil
		}
		if n.Status == "ROLLBACK_COMPLETE" || n.Status == "ROLLBACK_FAILED" || n.Status == "UPDATE_FAILED" {

			return nil, "", fmt.Errorf("%s: %s", n.Status, n.StatusReason)
		}

		return n, n.Status, nil
	}
}
