package circonus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccCirconusMetric_basic(t *testing.T) {
	metricAvgName := fmt.Sprintf("Average Ping Time - %s", acctest.RandString(5))
	metricMaxName := fmt.Sprintf("Maximum Ping Time - %s", acctest.RandString(5))
	metricMinName := fmt.Sprintf("Minimum Ping Time - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyMetric,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusMetricConfigFmt, metricAvgName, metricMaxName, metricMinName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_average", "name", metricAvgName),
					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_average", "active", "false"),
					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_average", "type", "numeric"),

					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_maximum", "name", metricMaxName),
					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_maximum", "active", "true"),
					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_maximum", "type", "numeric"),

					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_minimum", "name", metricMinName),
					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_minimum", "active", "true"),
					resource.TestCheckResourceAttr("circonus_metric.icmp_ping_minimum", "type", "numeric"),
				),
			},
		},
	})
}

func TestAccCirconusMetric_tagsets(t *testing.T) {
	metricName := fmt.Sprintf("foo - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyMetric,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusMetricTagsFmt0, metricName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_metric.t", "name", metricName),
					resource.TestCheckResourceAttr("circonus_metric.t", "type", "numeric"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCirconusMetricTagsFmt1, metricName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_metric.t", "name", metricName),
					resource.TestCheckResourceAttr("circonus_metric.t", "type", "numeric"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCirconusMetricTagsFmt2, metricName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_metric.t", "name", metricName),
					resource.TestCheckResourceAttr("circonus_metric.t", "type", "numeric"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCirconusMetricTagsFmt3, metricName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_metric.t", "name", metricName),
					resource.TestCheckResourceAttr("circonus_metric.t", "type", "numeric"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCirconusMetricTagsFmt4, metricName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_metric.t", "name", metricName),
					resource.TestCheckResourceAttr("circonus_metric.t", "type", "numeric"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCirconusMetricTagsFmt5, metricName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_metric.t", "name", metricName),
					resource.TestCheckResourceAttr("circonus_metric.t", "type", "numeric"),
				),
			},
		},
	})
}

func testAccCheckDestroyMetric(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circonus_metric" {
			continue
		}

		id := rs.Primary.ID
		exists := id == ""
		switch {
		case !exists:
			// noop
		case exists:
			return fmt.Errorf("metric still exists after destroy")
		}
	}

	return nil
}

const testAccCirconusMetricConfigFmt = `
resource "circonus_metric" "icmp_ping_average" {
  name = "%s"
  active = false
  type = "numeric"
}

resource "circonus_metric" "icmp_ping_maximum" {
  name = "%s"
  active = true
  type = "numeric"
}

resource "circonus_metric" "icmp_ping_minimum" {
  name = "%s"
  type = "numeric"
}
`

const testAccCirconusMetricTagsFmt0 = `
resource "circonus_metric" "t" {
  name = "%s"
  type = "numeric"
}
`

const testAccCirconusMetricTagsFmt1 = `
resource "circonus_metric" "t" {
  name = "%s"
  type = "numeric"
}
`

const testAccCirconusMetricTagsFmt2 = `
resource "circonus_metric" "t" {
  name = "%s"
  type = "numeric"
}
`

const testAccCirconusMetricTagsFmt3 = `
resource "circonus_metric" "t" {
  name = "%s"
  type = "numeric"
}
`

const testAccCirconusMetricTagsFmt4 = `
resource "circonus_metric" "t" {
  name = "%s"
  type = "numeric"
}
`

const testAccCirconusMetricTagsFmt5 = `
resource "circonus_metric" "t" {
  name = "%s"
  type = "numeric"
}
`
