package circonus

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCirconusCheckDNS_basic(t *testing.T) {
	checkName := fmt.Sprintf("DNS check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckDNSConfigFmt,
					checkName,
					testAccBroker1,
					testAccBroker2,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.google", "active", "true"),
					resource.TestCheckNoResourceAttr("circonus_check.google", "check_id"),
					resource.TestCheckResourceAttr("circonus_check.google", "checks.#", "2"),
					resource.TestMatchResourceAttr("circonus_check.google", "checks.0", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestMatchResourceAttr("circonus_check.google", "checks.1", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckNoResourceAttr("circonus_check.google", "check_id"),
					resource.TestCheckResourceAttr("circonus_check.google", "check_by_collector.%", "2"),
					resource.TestCheckResourceAttr("circonus_check.google", "collector.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.google", "collector.0.id", testAccBroker1),
					resource.TestCheckResourceAttr("circonus_check.google", "dns.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.google", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.google", "period", "300s"),
					resource.TestCheckResourceAttr("circonus_check.google", "metric.#", "3"),
					resource.TestCheckResourceAttr("circonus_check.google", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.google", "target", "api.circonus.com"),
					resource.TestCheckResourceAttr("circonus_check.google", "type", "dns"),
				),
			},
		},
	})
}

const testAccCirconusCheckDNSConfigFmt = `
variable "test_tags" {
  type = list(string)
  default = [ "author:terraform", "lifecycle:unittest" ]
}
resource "circonus_check" "google" {
  active = true
  name = "%s"
  period = "300s"

  collector {
    id = "%s"
  }

  collector {
    id = "%s"
  }

  dns {
    query = "google.com"
    rtype = "A"
  }

  metric {
    name = "answer"
    type = "text"
  }

  metric {
    name = "rtt"
    type = "numeric"
  }

  metric {
    name = "ttl"
    type = "numeric"
  }

  tags = "${var.test_tags}"
  target = "api.circonus.com"
}
`
