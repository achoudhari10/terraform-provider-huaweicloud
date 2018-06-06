package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var node_name = fmt.Sprintf("%s",OS_NODE_NAME)
func TestAccCCENodesV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCCENode(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccCCENodeV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3DataSourceID("data.huaweicloud_cce_nodes_v3.nodes"),
					resource.TestCheckResourceAttr("data.huaweicloud_cce_nodes_v3.nodes", "name", node_name),
				),
			},
		},
	})
}
func testAccCheckCCENodeV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find nodes data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Node data source ID not set ")
		}

		return nil
	}
}


var testAccCCENodeV3DataSource_basic = fmt.Sprintf(`
data "huaweicloud_cce_nodes_v3" "nodes" {
		cluster_id ="%s"
		name = "%s"
}
`, OS_CLUSTER_ID,OS_NODE_NAME)

