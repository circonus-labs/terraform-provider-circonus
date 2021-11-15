package circonus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceCirconusCollector(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceCirconusCollectorConfig,
					testAccBroker1,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceCirconusCollectorCheck("data.circonus_collector.by_id", testAccBroker1),
				),
			},
		},
	})
}

func testAccDataSourceCirconusCollectorCheck(name, cid string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("root module has no resource called %s", name)
		}

		attr := rs.Primary.Attributes

		if attr[collectorIDAttr] != cid {
			return fmt.Errorf("bad id %s", attr[collectorIDAttr])
		}

		return nil
	}
}

const testAccDataSourceCirconusCollectorConfig = `
data "circonus_collector" "by_id" {
  id = "%s"
}
`
