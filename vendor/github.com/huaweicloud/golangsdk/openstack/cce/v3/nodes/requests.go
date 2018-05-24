package nodes

import (
	"github.com/huaweicloud/golangsdk"
	//"github.com/huaweicloud/golangsdk/pagination"
	//"reflect"
)

var RequestOpts golangsdk.RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json"},
}
// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.

type CreateOpts struct {
	// API type, fixed value Node
	Kind string `json:"kind" required:"true"`
	//API version, fixed value v3
	ApiVersion string `json:"apiversion" required:"true"`
	//Medata required to create a Node
	Metadata CreateMetaData `json:"metadata"`
	//specifications to create a Node
	Spec Spec `json:"spec" required:"true"`
}
//Medata required to create a Node
type CreateMetaData struct {
	//Node name
	Name string `json:"name,omitempty"`
	// Node tag, key/value pair format
	Labels map[string]string `json:"labels,omitempty"`
	//Node annotation, key/value pair format
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Spec struct {
	// Node specifications
	Flavor 		  string 			`json:"flavor" required:"true"`
	Az 			  string 			`json:"az"  required:"true"`
	Login   	  LoginSpec			`json:"login" required:"true"`
	RootVolume 	  VolumeSpec 		`json:"rootVolume" required:"true"`
	DataVolumes []VolumeSpec 		`json:"dataVolumes" required:"true"`
	PublicIP	  PublicIPSpec 		`json:"publicIP,omitempty"`
	BillingMode   int 				`json:"billingMode ,omitempty"`
	Count         int				`json:"count" required:"true"`
	ExtendParam   string 			`json:"extendParam ,omitempty"`
}
type LoginSpec struct {
	SshKey   	string 		`json:"sshKey" required:"true"`
}
type VolumeSpec struct {
	Size 			int 	`json:"size" required:"true"`
	VolumeType 		string 	`json:"volumetype" required:"true"`
	ExtendParam   	string 	`json:"extendParam ,omitempty"`
}
type PublicIPSpec struct {
	Ids   []string  `json:"ids ,omitempty"`
	Count int       `json:"count ,omitempty"`
	Eip   EipSpec   `json:"eip, omitempty"`
}
type EipSpec struct {
	IpType 		string 			`json:"iptype" required:"true"`
	Bandwidth 	BandwidthOpts 	`json:"bandwidth" required:"true"`
}
type BandwidthOpts struct {
	ChargeMode 	string 		`json:"chargemode ,omitempty"`
	Size 		int         `json:"size" required:"true"`
	ShareType 	string 		`json:"sharetype" required:"true"`
}


// Create accepts a CreateOpts struct and uses the values to create a new
// logical Node. When it is created, the Node does not have an internal
// interface
//
type CreateOptsBuilder interface {
	ToNodeCreateMap() (map[string]interface{}, error)
}

func (opts CreateOpts) ToNodeCreateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

func Create(c *golangsdk.ServiceClient, clusterid string, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToNodeCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	reqOpt := &golangsdk.RequestOpts{OkCodes: []int{201}}
	_, r.Err = c.Post(rootURL(c, clusterid), b, &r.Body, reqOpt)
	return
}

// Get retrieves a particular vpc based on its unique ID.
func Get(c *golangsdk.ServiceClient, clusterid, nodeid string) (r GetResult) {
	//_, r.Err = c.Get(resourceURL(c, clusterid, nodeid), &r.Body, nil)
	_, r.Err = c.Get(resourceURL(c, clusterid, nodeid), &r.Body, &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to the
// Update request.
type UpdateOptsBuilder interface {
	ToNodeUpdateMap() (map[string]interface{}, error)
}
type UpdateOpts struct {
	Metadata UpdateMetadata `json:"metadata,omitempty"`
}
type UpdateMetadata struct {
	Name string `json:"name,omitempty"`
}

// ToNodeUpdateMap builds an update body based on UpdateOpts.
func (opts UpdateOpts) ToNodeUpdateMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Update allows vpcs to be updated. You can update the name, administrative
// state, and the external gateway. For more information about how to set the
// external gateway for a vpc, see Create. This operation does not enable
// the update of vpc interfaces. To do this, use the AddInterface and
// RemoveInterface functions.
func Update(c *golangsdk.ServiceClient, clusterid, nodeid string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToNodeUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = c.Put(resourceURL(c, clusterid, nodeid), b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

func Delete(c *golangsdk.ServiceClient, clusterid, nodeid string) (r DeleteResult) {
	_, r.Err = c.Delete(resourceURL(c, clusterid, nodeid), &golangsdk.RequestOpts{
		OkCodes:     []int{200},
		MoreHeaders: RequestOpts.MoreHeaders, JSONBody: nil,
	})
	return
}