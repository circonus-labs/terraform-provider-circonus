package circonus

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCirconusCheckNTP_basic(t *testing.T) {
	checkName := fmt.Sprintf("NTP check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckNTPConfigFmt, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.clock", "active", "true"),
					resource.TestCheckNoResourceAttr("circonus_check.clock", "check_id"),
					resource.TestCheckResourceAttr("circonus_check.clock", "checks.#", "2"),
					resource.TestMatchResourceAttr("circonus_check.clock", "checks.0", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestMatchResourceAttr("circonus_check.clock", "checks.1", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckNoResourceAttr("circonus_check.clock", "check_id"),
					resource.TestCheckResourceAttr("circonus_check.clock", "check_by_collector.%", "2"),
					resource.TestCheckResourceAttr("circonus_check.clock", "collector.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.clock", "collector.0.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.clock", "ntp.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.clock", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.clock", "period", "300s"),
					resource.TestCheckResourceAttr("circonus_check.clock", "metric.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.clock", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.clock", "target", "10.1.1.1"),
					resource.TestCheckResourceAttr("circonus_check.clock", "type", "ntp"),
				),
			},
		},
	})
}

const testAccCirconusCheckNTPConfigFmt = `
variable "test_tags" {
  type = "list"
  default = [ "author:terraform", "lifecycle:unittest" ]
}
resource "circonus_check" "clock" {
  active = true
  name = "%s"
  period = "300s"

  collector {
    id = "/broker/1"
  }

  collector {
    id = "/broker/275"
  }

  ntp {
    port = 123
    use_control = false
  }

  metric {
    name = "offset"
    type = "numeric"
  }

  metric {
    name = "rtdisp"
    type = "numeric"
  }

  tags = "${var.test_tags}"
  target = "10.1.1.1"
}
`
