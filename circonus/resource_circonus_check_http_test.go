package circonus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCirconusCheckHTTP_basic(t *testing.T) {
	checkName := fmt.Sprintf("Terraform test: noit's jezebel availability check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckHTTPConfigFmt, checkName, 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.jezebel", "active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "collector.0.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.code", `^200$`),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.extract", `HTTP/1.1 200 OK`),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.headers.%", "1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.headers.Host", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.version", "1.1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.method", "GET"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.read_limit", "1048576"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.url", "http://127.0.0.1:8083/resmon"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.redirects", "2"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "notes", "Check to make sure jezebel is working as expected"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.#", "4"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.0.name", "code"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.0.type", "text"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.1.name", "duration"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.1.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.2.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.2.name", "tt_connect"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.2.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.3.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.3.name", "tt_firstbyte"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.3.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.0", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.1", "app:jezebel"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.2", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.3", "source:circonus"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "target", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "type", "http"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCirconusCheckHTTPConfigFmt, checkName, 0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.jezebel", "active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "collector.0.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.code", `^200$`),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.extract", `HTTP/1.1 200 OK`),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.headers.%", "1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.headers.Host", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.version", "1.1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.method", "GET"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.read_limit", "1048576"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.url", "http://127.0.0.1:8083/resmon"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "http.0.redirects", "0"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "notes", "Check to make sure jezebel is working as expected"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.#", "4"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.0.name", "code"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.0.type", "text"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.1.name", "duration"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.1.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.2.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.2.name", "tt_connect"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.2.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.3.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.3.name", "tt_firstbyte"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "metric.3.type", "numeric"),

					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.0", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.1", "app:jezebel"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.2", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "tags.3", "source:circonus"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "target", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.jezebel", "type", "http"),
				),
			},
		},
	})
}

const testAccCirconusCheckHTTPConfigFmt = `
variable "http_check_tags" {
  type = list(string)
  default = [ "app:circonus", "app:jezebel", "lifecycle:unittest", "source:circonus" ]
}

resource "circonus_metric" "status_code" {
  name = "code"
  type = "text"
}

resource "circonus_metric" "request_duration" {
  name = "duration"
  type = "numeric"
}

resource "circonus_metric" "request_ttconnect" {
  name = "tt_connect"
  type = "numeric"
}

resource "circonus_metric" "request_ttfb" {
  name = "tt_firstbyte"
  type = "numeric"
}

resource "circonus_check" "jezebel" {
  active = true
  name = "%s"
  notes = "Check to make sure jezebel is working as expected"
  period = "60s"

  collector {
    id = "/broker/1"
  }

  http {
    code = "^200$"
    extract     = "HTTP/1.1 200 OK"
    headers     = {
      Host = "127.0.0.1",
    }
    version     = "1.1"
    method      = "GET"
    read_limit  = 1048576
	url         = "http://127.0.0.1:8083/resmon"
	redirects   = "%d"
  }

  metric {
    name = "${circonus_metric.status_code.name}"
    type = "${circonus_metric.status_code.type}"
  }

  metric {
    name = "${circonus_metric.request_duration.name}"
    type = "${circonus_metric.request_duration.type}"
  }

  metric {
    name = "${circonus_metric.request_ttconnect.name}"
    type = "${circonus_metric.request_ttconnect.type}"
  }

  metric {
    name = "${circonus_metric.request_ttfb.name}"
    type = "${circonus_metric.request_ttfb.type}"
  }

  tags = "${var.http_check_tags}"
}
`
