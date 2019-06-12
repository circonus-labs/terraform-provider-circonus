package circonus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCirconusCheckJMX_basic(t *testing.T) {
	checkName := fmt.Sprintf("Terraform test: JMX check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckJMXConfigFmt, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.something", "active", "true"),
					resource.TestCheckResourceAttr("circonus_check.something", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.something", "collector.1893401625.id", "/broker/1286"),
					resource.TestCheckResourceAttr("circonus_check.something", "jmx.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.something", "jmx.453641246.host", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.something", "jmx.453641246.port", "9999"),
					resource.TestCheckResourceAttr("circonus_check.something", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.something", "notes", "Check to harvest JMX info"),
					resource.TestCheckResourceAttr("circonus_check.something", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.something", "mbean_domains.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.something", "mbean_properties.#", "2"),

					resource.TestCheckResourceAttr("circonus_check.something", "tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.something", "tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.something", "tags.213659730", "app:something"),
					resource.TestCheckResourceAttr("circonus_check.something", "tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.something", "tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.something", "target", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.something", "type", "tcp"),
				),
			},
		},
	})
}

const testAccCirconusCheckJMXConfigFmt = `
variable "tcp_check_tags" {
  type = "list"
  default = [ "app:circonus", "app:something", "lifecycle:unittest", "source:fastly" ]
}

resource "circonus_check" "something" {
  active = true
  name = "%s"
  notes = "Check to harvest JMX info"
  period = "60s"

  collector {
    id = "/broker/1286"
  }

  jmx {
    host = "127.0.0.1"
    port = 9999
    mbean_domains = ["foo", "bar"]
    mbean_properties {
       name = "Foo"
       type = "Thing"
    }
    mbean_properties {
       name = "Baz"
       type = "Quux"
    }
  }
  tags = [ "${var.tcp_check_tags}" ]
}
`
