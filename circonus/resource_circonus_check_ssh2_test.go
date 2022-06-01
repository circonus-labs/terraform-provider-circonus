package circonus

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCirconusCheckSSH2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckSSH2Config, testAccBroker1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.ssh2", "active", "true"),
					resource.TestMatchResourceAttr("circonus_check.ssh2", "check_id", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "collector.0.id", testAccBroker1),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "ssh2.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "ssh2.0.port", "22"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "ssh2.0.method_kex", "diffie-hellman-group14-sha1"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "ssh2.0.method_hostkey", "ssh-rsa"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "ssh2.0.method_comp_cs", "none"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "ssh2.0.method_comp_sc", "none"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "name", "Terraform test: api.circonus.com ssh2 check"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "notes", ""),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "metric.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "metric.0.name", "bits"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "metric.0.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "metric.1.name", "duration"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "metric.1.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "target", "10.1.1.1"),
					resource.TestCheckResourceAttr("circonus_check.ssh2", "type", "ssh2"),
				),
			},
		},
	})
}

const testAccCirconusCheckSSH2Config = `

resource "circonus_metric" "bits" {
  name = "bits"
  type = "numeric"
}

resource "circonus_metric" "duration" {
  name = "duration"
  type = "numeric"
}

resource "circonus_check" "ssh2" {
  active = true
  name = "Terraform test: api.circonus.com ssh2 check"
  period = "60s"

  collector {
    id = "%s"
  }

  ssh2 {
    port = 22
  }

  metric {
    name = "${circonus_metric.bits.name}"
    type = "${circonus_metric.bits.type}"
  }

  metric {
    name = "${circonus_metric.duration.name}"
    type = "${circonus_metric.duration.type}"
  }

  tags = [ "source:circonus", "lifecycle:unittest" ]
  target = "10.1.1.1"
}
`
