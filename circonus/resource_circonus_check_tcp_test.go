package circonus

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccCirconusCheckTCP_basic(t *testing.T) {
	checkName := fmt.Sprintf("Terraform test: TCP+TLS check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckTCPConfigFmt, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "collector.2388330941.id", "/broker/1"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "tcp.#", "1"),
					// resource.TestCheckResourceAttr("circonus_check.tls_cert", "tcp.453641246.banner_regexp", ""),
					// resource.TestCheckResourceAttr("circonus_check.tls_cert", "tcp.453641246.ca_chain", ""),
					// resource.TestCheckResourceAttr("circonus_check.tls_cert", "tcp.453641246.certificate_file", ""),
					// resource.TestCheckResourceAttr("circonus_check.tls_cert", "tcp.453641246.ciphers", ""),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "tcp.453641246.host", "127.0.0.1"),
					// resource.TestCheckResourceAttr("circonus_check.tls_cert", "tcp.453641246.key_file", ""),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "tcp.453641246.port", "443"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "notes", "Check to harvest cert expiration information"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.#", "9"),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.name", "cert_end"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.0.unit", "epoch"),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.name", "cert_end_in"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.1.unit", "seconds"),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.name", "cert_error"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.type", "text"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.2.unit", ""),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.name", "cert_issuer"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.type", "text"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.3.unit", ""),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.name", "cert_start"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.4.unit", "epoch"),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.name", "cert_subject"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.type", "text"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.5.unit", ""),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.name", "duration"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.6.unit", "milliseconds"),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.name", "tt_connect"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.7.unit", "milliseconds"),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.name", "tt_firstbyte"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "metric.8.unit", "milliseconds"),

					resource.TestCheckResourceAttr("circonus_check.tls_cert", "tags.#", "4"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "tags.30226350", "app:circonus"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "tags.213659730", "app:tls_cert"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "tags.862116066", "source:fastly"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "target", "127.0.0.1"),
					resource.TestCheckResourceAttr("circonus_check.tls_cert", "type", "tcp"),
				),
			},
		},
	})
}

const testAccCirconusCheckTCPConfigFmt = `
variable "tcp_check_tags" {
  type = "list"
  default = [ "app:circonus", "app:tls_cert", "lifecycle:unittest", "source:fastly" ]
}

resource "circonus_check" "tls_cert" {
  active = true
  name = "%s"
  notes = "Check to harvest cert expiration information"
  period = "60s"

  collector {
    id = "/broker/1"
  }

  tcp {
    host = "127.0.0.1"
    port = 443
  }

  metric {
    name = "cert_end"
    tags = "${var.tcp_check_tags}"
    type = "numeric"
    unit = "epoch"
  }

  metric {
    name = "cert_end_in"
    tags = "${var.tcp_check_tags}"
    type = "numeric"
    unit = "seconds"
  }

  metric {
    name = "cert_error"
    tags = "${var.tcp_check_tags}"
    type = "text"
  }

  metric {
    name = "cert_issuer"
    tags = "${var.tcp_check_tags}"
    type = "text"
  }

  metric {
    name = "cert_start"
    tags = "${var.tcp_check_tags}"
    type = "numeric"
    unit = "epoch"
  }

  metric {
    name = "cert_subject"
    tags = "${var.tcp_check_tags}"
    type = "text"
  }

  metric {
    name = "duration"
    tags = "${var.tcp_check_tags}"
    type = "numeric"
    unit = "milliseconds"
  }

  metric {
    name = "tt_connect"
    tags = "${var.tcp_check_tags}"
    type = "numeric"
    unit = "milliseconds"
  }

  metric {
    name = "tt_firstbyte"
    tags = "${var.tcp_check_tags}"
    type = "numeric"
    unit = "milliseconds"
  }

  tags = "${var.tcp_check_tags}"
}
`
