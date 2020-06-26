package circonus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCirconusCheckHTTPTrap_basic(t *testing.T) {
	checkName := fmt.Sprintf("Terraform test: consul server httptrap check- %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckHTTPTrapConfigFmt, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.consul", "active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.consul", "collector.1263561585.id", "/broker/35"),
					resource.TestCheckResourceAttr("circonus_check.consul", "httptrap.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.consul", "httptrap.2067899660.async_metrics", "false"),
					resource.TestCheckResourceAttr("circonus_check.consul", "httptrap.2067899660.secret", "12345"),
					resource.TestCheckResourceAttr("circonus_check.consul", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.consul", "notes", "Check to receive consul server telemetry"),
					resource.TestCheckResourceAttr("circonus_check.consul", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.consul", "metric.#", "3"),

					resource.TestCheckResourceAttr("circonus_check.consul", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul", "metric.0.name", "consul`consul-server-10-151-2-8`consul`session_ttl`active"),
					resource.TestCheckResourceAttr("circonus_check.consul", "metric.0.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.consul", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul", "metric.1.name", "consul`consul-server-10-151-2-8`runtime`alloc_bytes"),
					resource.TestCheckResourceAttr("circonus_check.consul", "metric.1.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.consul", "metric.2.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul", "metric.2.name", "consul`consul`http`GET`v1`kv`_"),
					resource.TestCheckResourceAttr("circonus_check.consul", "metric.2.type", "histogram"),

					resource.TestCheckResourceAttr("circonus_check.consul", "tags.#", "3"),
					resource.TestCheckResourceAttr("circonus_check.consul", "tags.3728194417", "app:consul"),
					resource.TestCheckResourceAttr("circonus_check.consul", "tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.consul", "tags.2058715988", "source:consul"),
					resource.TestCheckResourceAttr("circonus_check.consul", "target", "consul-server-10-151-2-8"),
					resource.TestCheckResourceAttr("circonus_check.consul", "type", "httptrap"),
				),
			},
		},
	})
}

const testAccCirconusCheckHTTPTrapConfigFmt = `
variable "httptrap_check_tags" {
  type = "list"
  default = [ "app:consul", "lifecycle:unittest", "source:consul" ]
}

variable "consul_hostname" {
  type = "string"
  default = "consul-server-10-151-2-8"
}

resource "circonus_check" "consul" {
  active = true
  name = "%s"
  notes = "Check to receive consul server telemetry"
  period = "60s"

  collector {
    id = "/broker/35"
  }

  httptrap {
    async_metrics = "false"
    secret = "12345"
  }

  metric {
    name = "consul` + "`" + `${var.consul_hostname}` + "`" + `consul` + "`" + `session_ttl` + "`" + `active"
    type = "numeric"
  }

  metric {
    name = "consul` + "`" + `${var.consul_hostname}` + "`" + `runtime` + "`" + `alloc_bytes"
    type = "numeric"
  }

  metric {
    name = "consul` + "`" + `consul` + "`" + `http` + "`" + `GET` + "`" + `v1` + "`" + `kv` + "`" + `_"
    type = "histogram"
  }

  tags = "${var.httptrap_check_tags}"
  target = "${var.consul_hostname}"
}
`
