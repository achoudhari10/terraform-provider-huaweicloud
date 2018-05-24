package clusters

import "github.com/huaweicloud/golangsdk"

const (
	rootpath = "clusters"
)

func rootURL(client *golangsdk.ServiceClient) string {
	return client.ServiceURL(rootpath)
}
func resourceURL(c *golangsdk.ServiceClient, id string) string {
	return c.ServiceURL(rootpath, id)
}
