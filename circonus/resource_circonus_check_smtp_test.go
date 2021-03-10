package circonus

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCirconusCheckSMTP_basic(t *testing.T) {
	checkName := fmt.Sprintf("SMTP check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckSMTPConfigFmt, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.smtp", "active", "true"),
					resource.TestMatchResourceAttr("circonus_check.smtp", "check_id", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckResourceAttr("circonus_check.smtp", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.smtp", "collector.0.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.smtp", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.smtp", "period", "300s"),
					resource.TestCheckResourceAttr("circonus_check.smtp", "metric.#", "3"),
					resource.TestCheckResourceAttr("circonus_check.smtp", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.smtp", "target", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.smtp", "type", "smtp"),
				),
			},
		},
	})
}

const testAccCirconusCheckSMTPConfigFmt = `
variable "test_tags" {
  type = "list"
  default = [ "author:terraform", "lifecycle:unittest" ]
}
resource "circonus_check" "smtp" {
  active = true
  name = "%s"
  period = "300s"

  collector {
    id = "/broker/1"
  }

  smtp {
    to = "test@example.com"
  }

  metric {
    name = "banner_time"
    type = "numeric"
  }

  metric {
    name = "ehlo_time"
    type = "numeric"
  }

  metric {
    name = "quit_time"
    type = "numeric"
  }

  tags = "${var.test_tags}"
  target = "127.0.0.1"
}
`
