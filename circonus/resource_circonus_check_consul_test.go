package circonus

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	consulAccBrokerEnvVar  = "TF_ACC_CIRC_ENT_BROKER_CID"
	consulAccBrokerSkipMsg = "'%s' missing from env, unable to test w/o enterprise broker w/resmon:consul enabled, skipping..."
)

func TestAccCirconusCheckConsul_node(t *testing.T) {
	accEnterpriseBrokerCID := os.Getenv(consulAccBrokerEnvVar)
	if accEnterpriseBrokerCID == "" {
		t.Skipf(consulAccBrokerSkipMsg, consulAccBrokerEnvVar)
	}

	checkName := fmt.Sprintf("Terraform test: consul.service.consul mode=state check - %s", acctest.RandString(5))

	checkNode := fmt.Sprintf("my-node-name-or-node-id-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckConsulConfigV1HealthNodeFmt, checkName, accEnterpriseBrokerCID, checkNode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.consul_server", "active", "true"),
					resource.TestMatchResourceAttr("circonus_check.consul_server", "check_id", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "collector.2084916526.id", accEnterpriseBrokerCID),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.dc", "dc2"),
					resource.TestCheckNoResourceAttr("circonus_check.consul_server", "consul.0.headers"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.http_addr", "http://consul.service.consul:8501"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.node", checkNode),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.node_blacklist.#", "3"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.node_blacklist.0", "a"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.node_blacklist.1", "bad"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.node_blacklist.2", "node"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "notes", ""),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.name", "KnownLeader"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.type", "text"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.name", "LastContact"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.2058715988", "source:consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "target", "consul.service.consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "type", "consul"),
				),
			},
		},
	})
}

func TestAccCirconusCheckConsul_service(t *testing.T) {
	accEnterpriseBrokerCID := os.Getenv(consulAccBrokerEnvVar)
	if accEnterpriseBrokerCID == "" {
		t.Skipf(consulAccBrokerSkipMsg, consulAccBrokerEnvVar)
	}

	checkName := fmt.Sprintf("Terraform test: consul.service.consul mode=service check - %s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckConsulConfigV1HealthServiceFmt, accEnterpriseBrokerCID, checkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.consul_server", "active", "true"),
					resource.TestMatchResourceAttr("circonus_check.consul_server", "check_id", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "collector.2084916526.id", accEnterpriseBrokerCID),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.#", "1"),
					resource.TestCheckNoResourceAttr("circonus_check.consul_server", "consul.0.headers"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.http_addr", "http://consul.service.consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.service", "consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.service_blacklist.#", "3"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.service_blacklist.0", "bad"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.service_blacklist.1", "hombre"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.service_blacklist.2", "service"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "notes", ""),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.name", "KnownLeader"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.type", "text"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.name", "LastContact"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.2058715988", "source:consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "target", "consul.service.consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "type", "consul"),
				),
			},
		},
	})
}

func TestAccCirconusCheckConsul_state(t *testing.T) {
	accEnterpriseBrokerCID := os.Getenv(consulAccBrokerEnvVar)
	if accEnterpriseBrokerCID == "" {
		t.Skipf(consulAccBrokerSkipMsg, consulAccBrokerEnvVar)
	}

	checkName := fmt.Sprintf("Terraform test: consul.service.consul mode=state check - %s", acctest.RandString(5))

	checkState := "critical"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusCheckBundle,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusCheckConsulConfigV1HealthStateFmt, checkName, accEnterpriseBrokerCID, checkState),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("circonus_check.consul_server", "active", "true"),
					resource.TestMatchResourceAttr("circonus_check.consul_server", "check_id", regexp.MustCompile(config.CheckCIDRegex)),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "collector.#", "1"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "collector.2084916526.id", accEnterpriseBrokerCID),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.#", "1"),
					resource.TestCheckNoResourceAttr("circonus_check.consul_server", "consul.0.headers"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.http_addr", "http://consul.service.consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.state", checkState),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.check_blacklist.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.check_blacklist.0", "worthless"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "consul.0.check_blacklist.1", "check"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "name", checkName),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "notes", ""),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "period", "60s"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.name", "KnownLeader"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.0.type", "text"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.active", "true"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.name", "LastContact"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "metric.1.type", "numeric"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.#", "2"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.1401442048", "lifecycle:unittest"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "tags.2058715988", "source:consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "target", "consul.service.consul"),
					resource.TestCheckResourceAttr("circonus_check.consul_server", "type", "consul"),
				),
			},
		},
	})
}

const testAccCirconusCheckConsulConfigV1HealthNodeFmt = `
resource "circonus_check" "consul_server" {
  active = true
  name = "%s"
  period = "60s"

  collector {
    id = "%s"
  }

  consul {
    dc = "dc2"
    http_addr = "http://consul.service.consul:8501"
    node = "%s"
    node_blacklist = ["a","bad","node"]
  }

  metric {
    name = "KnownLeader"
    type = "text"
  }

  metric {
    name = "LastContact"
    type = "numeric"
  }

  tags = [ "source:consul", "lifecycle:unittest" ]

  target = "consul.service.consul"
}
`

const testAccCirconusCheckConsulConfigV1HealthServiceFmt = `
resource "circonus_check" "consul_server" {
  active = true
  name = "%s"
  period = "60s"

  collector {
    id = "%s"
  }

  consul {
    service = "consul"
    service_blacklist = ["bad","hombre","service"]
  }

  metric {
    name = "KnownLeader"
    type = "text"
  }

  metric {
    name = "LastContact"
    type = "numeric"
  }

  tags = [ "source:consul", "lifecycle:unittest" ]

  target = "consul.service.consul"
}
`

const testAccCirconusCheckConsulConfigV1HealthStateFmt = `
resource "circonus_check" "consul_server" {
  active = true
  name = "%s"
  period = "60s"

  collector {
    id = "%s"
  }

  consul {
    state = "%s"
    check_blacklist = ["worthless","check"]
  }

  metric {
    name = "KnownLeader"
    type = "text"
  }

  metric {
    name = "LastContact"
    type = "numeric"
  }

  tags = [ "source:consul", "lifecycle:unittest" ]

  target = "consul.service.consul"
}
`
