package huaweicloud


import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/huaweicloud/golangsdk/openstack/rts/v1/stacks"
)

// PASS
func TestAccOTCRtsStackV1_basic(t *testing.T) {
	var stacks stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRtsStackV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsStackV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsStackV1Exists("opentelekomcloud_rts_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "name", "terraform_provider_stack"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "status", "CREATE_COMPLETE"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "description", "A HOT template that create a single server and boot from volume."),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "status_reason", "Stack CREATE completed successfully"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "disable_rollback", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "timeout_mins", "60"),


				),

			},
		},
	})
}

func TestAccOTCRtsStackV1_update(t *testing.T) {
	var stacks stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRtsStackV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsStackV1_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsStackV1Exists("opentelekomcloud_rts_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "name", "terraform_provider_stack"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "status", "CREATE_COMPLETE"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "description", "A HOT template that create a single server and boot from volume."),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "status_reason", "Stack CREATE completed successfully"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "disable_rollback", "true"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "timeout_mins", "60"),
				),
			},
			resource.TestStep{
				Config: testAccRtsStackV1_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsStackV1Exists("opentelekomcloud_rts_stack_v1.stack_1", &stacks),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "disable_rollback", "false"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "timeout_mins", "50"),
					resource.TestCheckResourceAttr(
						"opentelekomcloud_rts_stack_v1.stack_1", "status", "UPDATE_COMPLETE"),

				),
			},
		},
	})
}

// PASS
func TestAccOTCRtsStackV1_timeout(t *testing.T) {
	var stacks stacks.RetrievedStack

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOTCRtsStackV1Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRtsStackV1_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOTCRtsStackV1Exists("opentelekomcloud_rts_stack_v1.stack_1", &stacks),
				),
			},
		},
	})
}

func testAccCheckOTCRtsStackV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenTelekomCloud orchestration client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "opentelekomcloud_rts_stack_v1" {
			continue
		}

		stack , err := stacks.Get(orchestrationClient,"terraform_provider_stack" ,rs.Primary.ID).Extract()

			if stack.Status != "DELETE_COMPLETE" {
				return fmt.Errorf("Stack still exists %s", err)
			}

	}

	return nil
}

func testAccCheckOTCRtsStackV1Exists(n string, stack *stacks.RetrievedStack) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		orchestrationClient, err := config.orchestrationV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenTelekomCloud orchestration Client : %s", err)
		}

		found, err := stacks.Get(orchestrationClient, "terraform_provider_stack",rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("stack not found")
		}

		*stack = *found

		return nil
	}
}

const testAccRtsStackV1_basic = `
resource "opentelekomcloud_rts_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= true
  timeout_mins=60
  template = <<JSON
          {
    "outputs": {
      "str1": {
        "description": "The description of the nat server.",
        "value": {
          "get_resource": "random"
        }
      }
    },
    "heat_template_version": "2013-05-23",
    "description": "A HOT template that create a single server and boot from volume.",
    "parameters": {
      "key_name": {
        "type": "string",
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

}
`

const testAccRtsStackV1_update = `
resource "opentelekomcloud_rts_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= false
  timeout_mins=50
  template = <<JSON
           {
    "outputs": {
      "str1": {
        "description": "The description of the nat server.",
        "value": {
          "get_resource": "random"
        }
      }
    },
    "heat_template_version": "2013-05-23",
    "description": "A HOT template that create a single server and boot from volume.",
    "parameters": {
      "key_name": {
        "type": "string",
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

}
`
const testAccRtsStackV1_timeout = `
resource "opentelekomcloud_rts_stack_v1" "stack_1" {
  name = "terraform_provider_stack"
  disable_rollback= true
  timeout_mins=60

  template = <<JSON
          {
    "outputs": {
      "str1": {
        "description": "The description of the nat server.",
        "value": {
          "get_resource": "random"
        }
      }
    },
    "heat_template_version": "2013-05-23",
    "description": "A HOT template that create a single server and boot from volume.",
    "parameters": {
      "key_name": {
        "type": "string",
  		"default": "keysclick",
        "description": "Name of existing key pair for the instance to be created."
      }
    },
    "resources": {
      "random": {
        "type": "OS::Heat::RandomString",
        "properties": {
          "length": 6
        }
      }
    }
  }
JSON

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
