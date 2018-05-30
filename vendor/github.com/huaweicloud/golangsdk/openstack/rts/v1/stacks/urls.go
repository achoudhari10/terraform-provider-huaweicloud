package stacks

import (
	"github.com/huaweicloud/golangsdk"
)

func createURL(c *golangsdk.ServiceClient) string {
	return c.ServiceURL("stacks")
}
func listURL(c *golangsdk.ServiceClient) string {
	return createURL(c)
}

func getURL(c *golangsdk.ServiceClient, name, id string) string {
	return c.ServiceURL("stacks", name, id)
}
func updateURL(c *golangsdk.ServiceClient, name, id string) string {
	return getURL(c, name, id)
}

func deleteURL(c *golangsdk.ServiceClient, name, id string) string {
	return getURL(c, name, id)
}