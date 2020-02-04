package circonus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCirconusCheckRedis_basic(t *testing.T) {
	checkName := fmt.Sprintf("Terraform test: Redis check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckRedisConfigFmt, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.redis", "active", "true"),
					resource.TestCheckResourceAttr("circonus_check.redis", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.redis", "collector.2388330941.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.redis", "redis.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.redis", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.redis", "notes", "Check to grab redis metrics"),
					resource.TestCheckResourceAttr("circonus_check.redis", "period", "60s"),

					resource.TestCheckResourceAttr("circonus_check.redis", "tags.#", "3"),
					resource.TestCheckResourceAttr("circonus_check.redis", "target", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.redis", "type", "redis"),
				),
			},
		},
	})
}

const testAccCirconusCheckRedisConfigFmt = `
variable "tcp_check_tags" {
  type = "list"
  default = [ "app:redis", "lifecycle:unittest", "source:fastly" ]
}

resource "circonus_check" "redis" {
  active = true
  name = "%s"
  notes = "Check to grab redis metrics"
  period = "60s"
  target = "127.0.0.1"

  collector {
    id = "/broker/1"
  }

  redis {
  }

  metric_filter {
    type    = "allow"
    regex   = ".*"
    comment = "Allow all metrics"
  }
  tags = "${var.tcp_check_tags}"
}
`
