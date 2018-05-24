package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/cce/v3/clusters"
)

func TestAccCCEClusterV1_basic(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCCE(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCCEClusterV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV1Exists("huaweicloud_cce_cluster_v3.cluster_1", &cluster),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "name", "huaweicloud-cce-cluster"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "status", "Available"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "kind", "Cluster"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "api_version", "v3"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "cluster_type", "VirtualMachine"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "flavor", "cce.s1.small"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "cluster_version", "v1.7.3-r10"),
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "container_network_type", "overlay_l2"),
				),
			},
			resource.TestStep{
				Config: testAccCCEClusterV1_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"huaweicloud_cce_cluster_v3.cluster_1", "description", "new description"),
				),
			},
		},
	})
}

// PASS
func TestAccCCEClusterV1_timeout(t *testing.T) {
	var cluster clusters.Clusters

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCCE(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCCEClusterV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCCEClusterV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCCEClusterV1Exists("huaweicloud_cce_cluster_v3.cluster_1", &cluster),
				),
			},
		},
	})
}

func testAccCheckCCEClusterV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	cceClient, err := config.cceV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "huaweicloud_cce_cluster_v3" {
			continue
		}

		_, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Cluster still exists")
		}
	}

	return nil
}

func testAccCheckCCEClusterV1Exists(n string, cluster *clusters.Clusters) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		cceClient, err := config.cceV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating HuaweiCloud CCE client: %s", err)
		}

		found, err := clusters.Get(cceClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Metadata.Id != rs.Primary.ID {
			return fmt.Errorf("Cluster not found")
		}

		*cluster = *found

		return nil
	}
}

var testAccCCEClusterV1_basic = fmt.Sprintf(`
resource "huaweicloud_cce_cluster_v3" "cluster_1" {
  kind="Cluster"
  api_version="v3"
  name = "huaweicloud-cce-cluster"
  cluster_type="VirtualMachine"
  flavor="cce.s1.small"
  cluster_version = "v1.7.3-r10"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}`, OS_VPC_ID, OS_SUBNET_ID)

var testAccCCEClusterV1_update = fmt.Sprintf(`
resource "huaweicloud_cce_cluster_v3" "cluster_1" {
  kind="Cluster"
  api_version="v3"
  name = "huaweicloud-cce-cluster"
  cluster_type="VirtualMachine"
  flavor="cce.s1.small"
  cluster_version = "v1.7.3-r10"
  vpc_id="%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
  description="new description"
}`, OS_VPC_ID, OS_SUBNET_ID)

var testAccCCEClusterV1_timeout = fmt.Sprintf(`
resource "huaweicloud_cce_cluster_v3" "cluster_1" {
  kind="Cluster"
  api_version="v3"
  name = "huaweicloud-cce-cluster"
  cluster_type="VirtualMachine"
  flavor="cce.s1.small"
  cluster_version = "v1.7.3-r10"
  vpc_id=""%s"
  subnet_id="%s"
  container_network_type="overlay_l2"
}
`, OS_VPC_ID, OS_SUBNET_ID)
