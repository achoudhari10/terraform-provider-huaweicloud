package nodes
import "github.com/huaweicloud/golangsdk"

const rootPath = "clusters"
const resourcePath = "nodes"

func rootURL(c *golangsdk.ServiceClient, clusterid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath)
}
func resourceURL(c *golangsdk.ServiceClient, clusterid, nodeid string) string {
	return c.ServiceURL(rootPath, clusterid, resourcePath, nodeid)
}