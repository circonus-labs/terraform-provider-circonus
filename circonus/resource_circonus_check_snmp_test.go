package circonus

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCirconusCheckSNMP_basic(t *testing.T) {
	checkName := fmt.Sprintf("SNMP check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckSNMPConfigFmt, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.snmp", "active", "true"),
					resource.TestMatchResourceAttr("circonus_check.snmp", "check_id", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckResourceAttr("circonus_check.snmp", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "collector.2388330941.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.snmp", "period", "300s"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "metric.#", "3"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "target", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "type", "snmp"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "snmp.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "snmp.0.version", "3"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "snmp.0.port", "161"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "snmp.0.security_name", "admin"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "snmp.0.oid.#", "3"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "snmp.0.oid.0.name", "upsBatCapacity"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "snmp.0.oid.1.name", "upsBatTimeRemaining"),
					resource.TestCheckResourceAttr("circonus_check.snmp", "snmp.0.oid.2.name", "upsBatVoltage"),
				),
			},
		},
	})
}

const testAccCirconusCheckSNMPConfigFmt = `
variable "test_tags" {
  type = "list"
  default = [ "author:terraform", "lifecycle:unittest" ]
}
resource "circonus_check" "snmp" {
  active = true
  name = "%s"
  period = "300s"

  collector {
    id = "/broker/1"
  }

  snmp {
    version = "3"
    community = "public"
    security_name = "admin"
    auth_passphrase = "foo"
    auth_protocol = "SHA"
    port = 161
    separate_queries = false
    security_level = "authNoPriv"

    oid {
      name = "upsBatCapacity"
      path = ".1.3.6.1.4.1.318.1.1.1.2.2.1.0"
    }
    oid {
      name = "upsBatTimeRemaining"
      path = ".1.3.6.1.4.1.318.1.1.1.2.2.3.0"
    }
    oid {
      name = "upsBatVoltage"
      path = ".1.3.6.1.4.1.318.1.1.1.2.2.8.0"
    }
  }

  metric {
    name = "upsBatCapacity"
    type = "numeric"
  }

  metric {
    name = "upsBatTimeRemaining"
    type = "numeric"
  }

  metric {
    name = "upsBatVoltage"
    type = "numeric"
  }

  tags = "${var.test_tags}"
  target = "127.0.0.1"
}
`
