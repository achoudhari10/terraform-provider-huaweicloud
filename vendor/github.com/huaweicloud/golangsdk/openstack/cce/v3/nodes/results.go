package nodes

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

//Describes the Node Structure of cluster
type Node struct {
	Kind       string  `json:"kind"`
	Apiversion string  `json:"apiVersion"`
	Items      []Items `json:"items"`
}

//Individual nodes of the cluster
type Items struct {
	Kind       string   `json:"kind"`
	Apiversion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Spec       Spec     `json:"spec"`
	Status     Status   `json:"status"`
}

// Name, Status of the node
type Metadata struct {
	Name              string `json:"name"`
	Uid               string `json:"uid"`
	CreationTimestamp string `json:"creationTimestamp"`
	UpdateTimestamp   string `json:"updateTimestamp"`
	Labels            string `json:"labels"`
	Annotations       string `json:"annotations"`
}

// Describes Nodes specification
type Spec struct {
	Type        string   `json:"type"`
	Az          string   `json:"az"`
	Login       Login    `json:"login"`
	RootVolume  Volume   `json:"rootVolume"`
	DataVolumes []Volume `json:"dataVolumes"`
	PublicIP    PublicIP `json:"publicIP"`
	Count       int      `json:"count"`
	ExtendParam string   `json:"extendParam"`
	Flavor      string   `json:"flavor"`
	BillingMode int      `json:"billingMode"`
}

//Gives the current status of the node
type Status struct {
	Phase     string `json:"phase"`
	ServerID  string `json:"ServerID"`
	PublicIP  string `json:"PublicIP"`
	PrivateIP string `json:"privateIP"`
	JobID     string `json:"jobID"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
}
type Login struct {
	SshKey string `json:"sshKey"`
}

//Used by Root Volume and Data Volumes
type Volume struct {
	Size        int    `json:"size"`
	Volumetype  string `json:"volumetype"`
	ExtendParam string `json:"extendParam"`
}

type PublicIP struct {
	Ids   []string `json:"ids"`
	Count int      `json:"count"`
	Eip   Eip      `json:"eip"`
}

type Eip struct {
	Iptype    string    `json:"iptype"`
	Bandwidth Bandwidth `json:"bandwidth"`
}

type Bandwidth struct {
	Chargemode string `json:"chargemode"`
	Size       int    `json:"size"`
	Sharetype  string `json:"sharetype"`
}

type commonResult struct {
	golangsdk.Result
}

func (r commonResult) Extract() (*Node, error) {
	var n struct {
		Node *Node `json:"node"`
	}
	err := r.ExtractInto(&n)
	return n.Node, err
}

func (r commonResult) ExtractNode(opts ListOpts) ([]Items, error) {
	var s Node
	err := r.ExtractInto(&s)
	if err != nil {
		return nil, err
	}
	return FilterNodes(s.Items, opts)
}

type NodePage struct {
	pagination.LinkedPageBase
}

type ListResult struct {
	commonResult
}
