package circonus

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCirconusCheckMetricFilter_basic(t *testing.T) {
	checkName := fmt.Sprintf("Metric Filter check - %s", acctest.RandString(5))
	target := fmt.Sprintf("%s.circonus.com", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckMetricFilterConfigFmt, checkName, target),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "active", "true"),
					resource.TestCheckNoResourceAttr("circonus_check.metric_filter", "check_id"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "checks.#", "2"),
					resource.TestMatchResourceAttr("circonus_check.metric_filter", "checks.0", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestMatchResourceAttr("circonus_check.metric_filter", "checks.1", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckNoResourceAttr("circonus_check.metric_filter", "check_id"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "check_by_collector.%", "2"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "collector.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "collector.2388330941.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "icmp_ping.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "icmp_ping.979664239.availability", "100"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "icmp_ping.979664239.count", "5"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "icmp_ping.979664239.interval", "500ms"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "period", "300s"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "metric_filter.#", "3"),

					resource.TestCheckResourceAttr("circonus_check.metric_filter", "metric_filter.0.type", "allow"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "metric_filter.0.regex", "available"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "metric_filter.1.type", "allow"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "metric_filter.1.regex", "average"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "metric_filter.2.type", "deny"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "metric_filter.2.regex", ".*"),

					resource.TestCheckResourceAttr("circonus_check.metric_filter", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "tags.2087084518", "author:terraform"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "target", target),
					resource.TestCheckResourceAttr("circonus_check.metric_filter", "type", "ping_icmp"),
				),
			},
		},
	})
}

const testAccCirconusCheckMetricFilterConfigFmt = `
variable "test_tags" {
  type = "list"
  default = [ "author:terraform", "lifecycle:unittest" ]
}
resource "circonus_check" "metric_filter" {
  active = true
  name = "%s"
  period = "300s"

  collector {
    id = "/broker/1"
  }

  collector {
    id = "/broker/275"
  }

  icmp_ping {
    availability = "100.0"
    count = 5
    interval = "500ms"
  }

  metric_filter {
    type = "allow"
    regex = "available"
    comment = "Allow available percentage"
  }

  metric_filter {
    type = "allow"
    regex = "average"
    comment = "Allow average latency"
  }

  metric_filter {
    type = "deny"
    regex = ".*"
    comment = "Deny everything else"
  }

  tags = "${var.test_tags}"
  target = "%s"
}
`
