package circonus

import (
	"fmt"
	"strings"
	"testing"
	"time"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCirconusMaintenance_basic(t *testing.T) {
	checkName := fmt.Sprintf("ICMP Ping check - %s", acctest.RandString(5))
	startTime, _ := time.Parse(time.RFC3339, "2020-01-26T19:00:00-05:00")
	stopTime, _ := time.Parse(time.RFC3339, "2020-01-30T19:00:00-05:00")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusMaintenance,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusMaintenanceConfigFmt, checkName, startTime.Format(time.RFC3339), stopTime.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("circonus_maintenance.check-maintenance", "check"),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "start", startTime.Format(time.RFC3339)),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "stop", stopTime.Format(time.RFC3339)),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "notes", "foo notes"),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "severities.#", "5"),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "severities.0", "1"),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "severities.1", "2"),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "severities.2", "3"),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "severities.3", "4"),
					resource.TestCheckResourceAttr("circonus_maintenance.check-maintenance", "severities.4", "5"),
				),
			},
		},
	})
}

func testAccCheckDestroyCirconusMaintenance(s *terraform.State) error {
	ctxt := testAccProvider.Meta().(*providerContext)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circonus_maintenance" {
			continue
		}

		cid := rs.Primary.ID
		exists, err := checkMaintenanceExists(ctxt, api.CIDType(&cid))
		switch {
		case !exists:
			// noop
		case exists:
			return fmt.Errorf("maintenance still exists after destroy")
		case err != nil:
			return fmt.Errorf("Error checking maintenance: %v", err)
		}
	}

	return nil
}

func checkMaintenanceExists(c *providerContext, maintenanceCID api.CIDType) (bool, error) {
	m, err := c.client.FetchMaintenanceWindow(maintenanceCID)
	if err != nil {
		if strings.Contains(err.Error(), defaultCirconus404ErrorString) {
			return false, nil
		}

		return false, err
	}

	if api.CIDType(&m.CID) == maintenanceCID {
		return true, nil
	}

	return false, nil
}

const testAccCirconusMaintenanceConfigFmt = `
resource "circonus_check" "api_latency" {
  active = true
  name = "%s"
  period = "60s"

  collector {
    id = "/broker/1"
  }

  icmp_ping {
    count = 1
  }

  metric {
    name = "maximum"
    type = "numeric"
    unit = "seconds"
  }

  target = "api.circonus.com"
}

resource "circonus_maintenance" "check-maintenance" {
  check = circonus_check.api_latency.check_id
  start = "%s"
  stop = "%s"
  notes = "foo notes"
  severities = ["1", "2", "3", "4", "5"]
}

`
