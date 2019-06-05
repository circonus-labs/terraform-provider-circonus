package circonus

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCirconusWorksheet_basic(t *testing.T) {
	graphName := fmt.Sprintf("Test Graph - %s", acctest.RandString(5))
	checkName := fmt.Sprintf("ICMP Ping check - %s", acctest.RandString(5))
	worksheetName := fmt.Sprintf("Test worksheet - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusWorksheet,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusWorksheetConfigFmt, checkName, graphName, worksheetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("circonus_worksheet.test", "favourite"),
				),
			},
		},
	})
}

func testAccCheckDestroyCirconusWorksheet(s *terraform.State) error {
	ctxt := testAccProvider.Meta().(*providerContext)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circonus_worksheet" {
			continue
		}

		cid := rs.Primary.ID
		exists, err := checkWorksheetExists(ctxt, api.CIDType(&cid))
		switch {
		case !exists:
			// noop
		case exists:
			return fmt.Errorf("worksheet still exists after destroy")
		case err != nil:
			return fmt.Errorf("Error checking worksheet: %v", err)
		}
	}

	return nil
}

func checkWorksheetExists(c *providerContext, worksheetCID api.CIDType) (bool, error) {
	rs, err := c.client.FetchWorksheet(worksheetCID)
	if err != nil {
		if strings.Contains(err.Error(), defaultCirconus404ErrorString) {
			return false, nil
		}

		return false, err
	}

	if api.CIDType(&rs.CID) == worksheetCID {
		return true, nil
	}

	return false, nil
}

const testAccCirconusWorksheetConfigFmt = `
variable "test_tags" {
  type = "list"
  default = [ "author:terraform", "lifecycle:unittest" ]
}

resource "circonus_check" "api_latency" {
  active = true
  name = "%s"
  period = "60s"

  collector {
    id = "/broker/1"
  }

  icmp_ping {
    count = 5
  }

  metric {
    name = "maximum"
    tags = [ "${var.test_tags}" ]
    type = "numeric"
    unit = "seconds"
  }

  metric {
    name = "minimum"
    tags = [ "${var.test_tags}" ]
    type = "numeric"
    unit = "seconds"
  }

  tags = [ "${var.test_tags}" ]
  target = "api.circonus.com"
}

resource "circonus_graph" "mixed-points" {
  name = "%s"
  description = "Terraform Test: mixed graph"
  notes = "test notes"
  graph_style = "line"
  line_style = "stepped"

  metric {
    # caql = "" # conflicts with metric_name/check
    check = "${circonus_check.api_latency.checks[0]}"
    metric_name = "maximum"
    metric_type = "numeric"
    name = "Maximum Latency"
    axis = "left" # right
    color = "#657aa6"
    function = "gauge"
    active = true
  }

  metric {
    # caql = "" # conflicts with metric_name/check
    check = "${circonus_check.api_latency.checks[0]}"
    metric_name = "minimum"
    metric_type = "numeric"
    name = "Minimum Latency"
    axis = "right" # left
    color = "#657aa6"
    function = "gauge"
    active = true
  }

  left {
    max = 11
  }

  right {
    logarithmic = 10
    max = 20
    min = -1
  }

  tags = [ "${var.test_tags}" ]
}

resource "circonus_worksheet" "test" {
  title = "%s"
  graphs = [
    "${circonus_graph.mixed-points.id}",
  ]
}
`
