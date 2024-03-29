package circonus

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCirconusCheckJSON_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckJSONConfig1, testAccBroker1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.usage", "active", "true"),
					resource.TestMatchResourceAttr("circonus_check.usage", "check_id", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckResourceAttr("circonus_check.usage", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.usage", "collector.0.id", testAccBroker1),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.headers.%", "3"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.headers.Accept", "application/json"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.headers.X-Circonus-App-Name", "TerraformCheck"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.headers.X-Circonus-Auth-Token", "<env 'CIRCONUS_API_TOKEN'>"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.version", "1.0"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.method", "GET"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.port", "443"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.read_limit", "1048576"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.url", "https://api.circonus.com/account/current"),
					resource.TestCheckResourceAttr("circonus_check.usage", "name", "Terraform test: api.circonus.com metric usage check"),
					resource.TestCheckResourceAttr("circonus_check.usage", "notes", ""),
					resource.TestCheckResourceAttr("circonus_check.usage", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.0.name", "_usage`0`_limit"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.0.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.1.name", "_usage`0`_used"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.1.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.usage", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.usage", "target", "api.circonus.com"),
					resource.TestCheckResourceAttr("circonus_check.usage", "type", "json"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCirconusCheckJSONConfig2, testAccBroker1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.usage", "active", "true"),
					resource.TestCheckResourceAttr("circonus_check.usage", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.usage", "collector.0.id", testAccBroker1),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.headers.%", "3"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.headers.Accept", "application/json"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.headers.X-Circonus-App-Name", "TerraformCheck"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.headers.X-Circonus-Auth-Token", "<env 'CIRCONUS_API_TOKEN'>"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.version", "1.1"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.method", "GET"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.port", "443"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.read_limit", "1048576"),
					resource.TestCheckResourceAttr("circonus_check.usage", "json.0.url", "https://api.circonus.com/account/current"),
					resource.TestCheckResourceAttr("circonus_check.usage", "name", "Terraform test: api.circonus.com metric usage check"),
					resource.TestCheckResourceAttr("circonus_check.usage", "notes", "notes!"),
					resource.TestCheckResourceAttr("circonus_check.usage", "period", "300s"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.0.name", "_usage`0`_limit"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.0.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.1.name", "_usage`0`_used"),
					resource.TestCheckResourceAttr("circonus_check.usage", "metric.1.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.usage", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.usage", "target", "api.circonus.com"),
					resource.TestCheckResourceAttr("circonus_check.usage", "type", "json"),
				),
			},
		},
	})
}

const testAccCirconusCheckJSONConfig1 = `

resource "circonus_metric" "limit" {
  name = "_usage` + "`0`" + `_limit"
  type = "numeric"
}

resource "circonus_metric" "used" {
  name = "_usage` + "`0`" + `_used"
  type = "numeric"
}

resource "circonus_check" "usage" {
  active = true
  name = "Terraform test: api.circonus.com metric usage check"
  period = "60s"

  collector {
    id = "%s"
  }

  json {
    url = "https://api.circonus.com/account/current"
    headers = {
      Accept                = "application/json",
      X-Circonus-App-Name   = "TerraformCheck",
      X-Circonus-Auth-Token = "<env 'CIRCONUS_API_TOKEN'>",
    }
    version = "1.0"
    method = "GET"
    port = 443
    read_limit = 1048576
  }

  metric {
    name = "${circonus_metric.limit.name}"
    type = "${circonus_metric.limit.type}"
  }

  metric {
    name = "${circonus_metric.used.name}"
    type = "${circonus_metric.used.type}"
  }

  tags = [ "source:circonus", "lifecycle:unittest" ]
}
`

const testAccCirconusCheckJSONConfig2 = `

resource "circonus_metric" "limit" {
  name = "_usage` + "`0`" + `_limit"
  type = "numeric"
}

resource "circonus_metric" "used" {
  name = "_usage` + "`0`" + `_used"
  type = "numeric"
}

resource "circonus_check" "usage" {
  active = true
  name = "Terraform test: api.circonus.com metric usage check"
  notes = "notes!"
  period = "300s"

  collector {
    id = "%s"
  }

  json {
    url = "https://api.circonus.com/account/current"
    headers = {
      Accept                = "application/json",
      X-Circonus-App-Name   = "TerraformCheck",
      X-Circonus-Auth-Token = "<env 'CIRCONUS_API_TOKEN'>",
    }
    version = "1.1"
    method = "GET"
    port = 443
    read_limit = 1048576
  }

  metric {
    name = "${circonus_metric.limit.name}"
    type = "${circonus_metric.limit.type}"
  }

  metric {
    name = "${circonus_metric.used.name}"
    type = "${circonus_metric.used.type}"
  }

  tags = [ "source:circonus", "lifecycle:unittest" ]
}
`
