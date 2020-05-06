package circonus

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCirconusCheckJMX_basic(t *testing.T) {
	statsdAccBrokerEnvVar := "TF_ACC_CIRC_ENT_BROKER_CID"
	statsdAccBrokerSkipMsg := "'%s' missing from env, unable to test w/o enterprise broker w/jmx enabled, skipping..."
	accEnterpriseBrokerCID := os.Getenv(statsdAccBrokerEnvVar)
	if accEnterpriseBrokerCID == "" {
		t.Skipf(statsdAccBrokerSkipMsg, statsdAccBrokerEnvVar)
	}

	checkName := fmt.Sprintf("Terraform test: JMX check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckJMXConfigFmt, checkName, accEnterpriseBrokerCID),
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
					resource.TestCheckResourceAttr("circonus_check.something", "mbean_domains.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.something", "mbean_properties.#", "1"),

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
variable "jmx_check_tags" {
  type = "list"
  default = [ "app:circonus", "app:something", "lifecycle:unittest", "source:fastly" ]
}

resource "circonus_check" "something" {
  active = true
  name = "%s"
  notes = "Check to harvest JMX info"
  period = "60s"
  target = "foo.foo.com"

  collector {
    id = "%s"
  }

  jmx {
    host = "127.0.0.1"
    port = 9999
    mbean_domains = ["foo"]
    mbean_properties {
       name = "Foo"
       type = "Thing"
       index = 1
    }
  }

  metric {
    name = "Foo"
    type = "numeric"
  }

  tags = var.jmx_check_tags
}
`
