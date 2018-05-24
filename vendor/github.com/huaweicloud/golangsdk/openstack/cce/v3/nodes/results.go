package nodes
import (
	"github.com/huaweicloud/golangsdk"
	//"github.com/huaweicloud/golangsdk/pagination"
)

type Nodes struct {
	Kind      	string 		`json:"kind"`
	Apiversion 	string 		`json:"apiversion"`
	Metadata  	Metadata   	`json:"metadata"`
	Spec   		Spec    	`json:"spec"`
	Status   	Status  	`json:"status"`
}

type Metadata struct {
	//Node unique id
	//Id string `json:"uid"`
	Id string
	//var Id = "cec124c2-58f1-11e8-ad73-0255ac101926" "cf4bc001-58f1-11e8-ad73-0255ac101926"
	//Node  name
	Name string `json:"name"`
	// Node tag, key/value pair format
	Labels map[string]string `json:"labels"`
	//Node annotation, key/value pair format
	Annotations map[string]string `json:"annotations"`
}

type Status struct {
	//The state of the Node
	Phase string `json:"phase"`
	//The ID of the Job that is operating asynchronously in the Node
	JobID string `json:"jobID"`
	//Reasons for the Node to become current
	Reason string `json:"reason"`
	//The status of each component in the Node
	Conditions Conditions `json:"conditions"`

}
type Conditions struct {
	//The type of component
	Type string `json:"type"`
	//The state of the component
	Status string `json:"status"`
	//The reason that the component becomes current
	Reason string `json:"reason"`
}

// Extract is a function that accepts a result and extracts a Node.
func (r commonResult) Extract() (*Nodes, error) {
	var s Nodes
	err := r.ExtractInto(&s)
	s.Metadata.Id="12345678899009767900865432246"
	return &s, err
}

type commonResult struct {
	golangsdk.Result
}

// CreateResult represents the result of a create operation. Call its Extract
// method to interpret it as a Node.
type CreateResult struct {
	commonResult
}
// GetResult represents the result of a get operation. Call its Extract
// method to interpret it as a Node.
type GetResult struct {
	commonResult
}
// UpdateResult represents the result of an update operation. Call its Extract
// method to interpret it as a Node.
type UpdateResult struct {
	commonResult
}
// DeleteResult represents the result of a delete operation. Call its ExtractErr
// method to determine if the request succeeded or failed.
type DeleteResult struct {
	golangsdk.ErrResult
}


