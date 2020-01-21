package circonus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceCirconusAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCirconusAccountCurrentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCirconusAccountCheck("data.circonus_account.by_current", "/account/4536"),
				),
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCirconusAccountIDConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCirconusAccountCheck("data.circonus_account.by_id", "/account/4536"),
				),
			},
		},
	})
}

func testAccDataSourceCirconusAccountCheck(name, cid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		if attr[accountIDAttr] != cid {
			return fmt.Errorf("bad %s %s", accountIDAttr, attr[accountIDAttr])
		}

		return nil
	}
}

const testAccDataSourceCirconusAccountCurrentConfig = `
data "circonus_account" "by_current" {
  current = true
}
`

const testAccDataSourceCirconusAccountIDConfig = `
data "circonus_account" "by_id" {
  id = "/account/4536"
}
`
