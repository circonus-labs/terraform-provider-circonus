package circonus

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCirconusCheckICMPPing_basic(t *testing.T) {
	checkName := fmt.Sprintf("ICMP Ping check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckICMPPingConfigFmt,
					checkName,
					testAccBroker1,
					testAccBroker2,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "active", "true"),
					resource.TestCheckNoResourceAttr("circonus_check.loopback_latency", "check_id"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "checks.#", "2"),
					resource.TestMatchResourceAttr("circonus_check.loopback_latency", "checks.0", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestMatchResourceAttr("circonus_check.loopback_latency", "checks.1", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckNoResourceAttr("circonus_check.loopback_latency", "check_id"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "check_by_collector.%", "2"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "collector.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "collector.0.id", testAccBroker2),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "icmp_ping.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "icmp_ping.0.availability", "100"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "icmp_ping.0.count", "5"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "icmp_ping.0.interval", "500ms"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "period", "300s"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.#", "5"),

					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.0.name", "available"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.0.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.1.name", "average"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.1.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.2.name", "count"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.2.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.3.name", "maximum"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.3.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.4.name", "minimum"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "metric.4.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "tags.0", "author:terraform"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "tags.1", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "target", "api.circonus.com"),
					resource.TestCheckResourceAttr("circonus_check.loopback_latency", "type", "ping_icmp"),
				),
			},
		},
	})
}

const testAccCirconusCheckICMPPingConfigFmt = `
variable "test_tags" {
  type = list(string)
  default = [ "author:terraform", "lifecycle:unittest" ]
}
resource "circonus_check" "loopback_latency" {
  active = true
  name = "%s"
  period = "300s"

  collector {
    id = "%s"
  }

  collector {
    id = "%s"
  }

  icmp_ping {
    availability = "100.0"
    count = 5
    interval = "500ms"
  }

  metric {
    name = "available"
    type = "numeric"
  }

  metric {
    name = "average"
    type = "numeric"
  }

  metric {
    name = "count"
    type = "numeric"
  }

  metric {
    name = "maximum"
    type = "numeric"
  }

  metric {
    name = "minimum"
    type = "numeric"
  }

  tags = "${var.test_tags}"
  target = "api.circonus.com"
}
`
