package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHuaweiCloudCCENodesV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCCENode(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{

			resource.TestStep{
				Config: testAccHuaweiCloudCCENodeV3DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCENodeV3DataSourceID("data.huaweicloud_cce_nodes_v3.nodes"),
					resource.TestCheckResourceAttr("data.huaweicloud_cce_nodes_v3.nodes", "name", "c2c-node"),
					resource.TestCheckResourceAttr("data.huaweicloud_cce_nodes_v3.nodes", "phase", "Abnormal"),
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
			return fmt.Errorf("node data source ID not set ")
		}

		return nil
	}
}

var testAccHuaweiCloudCCENodeV3DataSource_basic = fmt.Sprintf(`
	data "huaweicloud_cce_nodes_v3" "nodes" {
		cluster_id ="cec124c2-58f1-11e8-ad73-0255ac101926"
		name = "c2c-node"
}
`)
