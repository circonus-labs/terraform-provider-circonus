package circonus

import (
	"fmt"
	"strings"
	"testing"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var rulesetGroupCheckName = fmt.Sprintf("ICMP Ping check - %s", acctest.RandString(5))

func TestAccCirconusRuleSetGroup_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDestroyCirconusRuleSetGroup,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCirconusRuleSetGroupConfigFmt, rulesetGroupCheckName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("circonus_rule_set.icmp_max_latency", "check"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_max_latency", "metric_name", "maximum"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_max_latency", "metric_type", "numeric"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_max_latency", "notes", "icmp max latency"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_max_latency", "if.#", "7"),

					resource.TestCheckResourceAttrSet("circonus_rule_set.icmp_min_latency", "check"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_min_latency", "metric_name", "minimum"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_min_latency", "metric_type", "numeric"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_min_latency", "notes", "icmp min latency"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_min_latency", "if.#", "1"),

					resource.TestCheckResourceAttrSet("circonus_rule_set.icmp_avg_latency", "check"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_avg_latency", "metric_name", "average"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_avg_latency", "metric_type", "numeric"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_avg_latency", "notes", "icmp avg latency"),
					resource.TestCheckResourceAttr("circonus_rule_set.icmp_avg_latency", "if.#", "3"),

					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "name", "icmp latency group 1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "notify.#", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "notify.0.sev1.#", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "notify.0.sev2.#", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "notify.0.sev3.#", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "formula.0.expression", "A and B"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "formula.0.raise_severity", "3"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "formula.0.wait", "0"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_1", "condition.#", "3"),

					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "name", "icmp latency group 2"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "notify.#", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "notify.0.sev1.#", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "notify.0.sev2.#", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "notify.0.sev3.#", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "formula.0.expression", "2"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "formula.0.raise_severity", "2"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "formula.0.wait", "1"),
					resource.TestCheckResourceAttr("circonus_rule_set_group.icmp_latency_2", "condition.#", "3"),
				),
			},
		},
	})
}

func testAccCheckDestroyCirconusRuleSetGroup(s *terraform.State) error {
	ctxt := testAccProvider.Meta().(*providerContext)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circonus_rule_set_group" {
			continue
		}

		cid := rs.Primary.ID
		exists, err := checkRuleSetGroupExists(ctxt, api.CIDType(&cid))
		switch {
		case !exists:
			// noop
		case exists:
			return fmt.Errorf("rule set still exists after destroy")
		case err != nil:
			return fmt.Errorf("Error checking rule set: %v", err)
		}
	}

	return nil
}

func checkRuleSetGroupExists(c *providerContext, ruleSetGroupCID api.CIDType) (bool, error) {
	rs, err := c.client.FetchRuleSetGroup(ruleSetGroupCID)
	if err != nil {
		if strings.Contains(err.Error(), defaultCirconus404ErrorString) {
			return false, nil
		}

		return false, err
	}

	if api.CIDType(&rs.CID) == ruleSetGroupCID {
		return true, nil
	}

	return false, nil
}

const testAccCirconusRuleSetGroupConfigFmt = `
variable "test_tags" {
  type = list(string)
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
    count = 1
  }

  metric {
    name = "average"
    type = "numeric"
  }

  metric {
    name = "maximum"
    type = "numeric"
  }

  metric {
    name = "minimum"
    type = "numeric"
  }

  tags = "${var.test_tags}"
  target = "api.circonus.com"
}

resource "circonus_rule_set" "icmp_max_latency" {
  check = "${circonus_check.api_latency.checks[0]}"
  metric_name = "maximum"
  notes = "icmp max latency"

  if {
    value {
      absent = "70"
    }

    then {
      notify = [
        "/contact_group/4680",
        "/contact_group/4679"
      ]
      severity = 1
    }
  }

  if {
    value {
      over {
        atleast = "30"
        last = "120"
        using = "average"
      }
      min_value = 2
    }

    then {
      notify = [ "/contact_group/4679" ]
      severity = 2
    }
  }

  if {
    value {
      over {
        atleast = "30"
        last = "180"
        using = "average"
      }

      max_value = 300
    }

    then {
      notify = [ "/contact_group/4679" ]
      severity = 3
    }
  }

  if {
    value {
      max_value = 400
    }

    then {
      notify = [ "/contact_group/4679" ]
      after = "2400"
      severity = 4
    }
  }

  if {
    value {
      max_value = 500
    }

    then {
      severity = 0
    }
  }

  if {
    value {
      eq_value = 600
    }
    then {
      severity = 0
    }
  }

  if {
    value {
      neq_value = 600
    }
    then {
      severity = 0
    }
  }
}

resource "circonus_rule_set" "icmp_min_latency" {
  check = "${circonus_check.api_latency.checks[0]}"
  metric_name = "minimum"
  notes = "icmp min latency"

  if {
    value {
      absent = "70"
    }

    then {
      notify = [
        "/contact_group/4680",
        "/contact_group/4679"
      ]
      severity = 3
    }
  }
}

resource "circonus_rule_set" "icmp_avg_latency" {
  check = "${circonus_check.api_latency.checks[0]}"
  metric_name = "average"
  notes = "icmp avg latency"

  if {
    value {
      absent = "300"
    }
    then {
      severity = 1
      notify = [
        "/contact_group/4680",
      ]
    }
  }
  if {
    value {
      absent = "70"
    }
    then {
      severity = 3
      notify = [
        "/contact_group/4680",
      ]
    }
  }
  if {
    value {
      max_value = "8000"
      over {
        atleast = "0"
        last    = "180"
        using   = "average"
      }
    }
    then {
      notify = [
        "/contact_group/4680",
      ]
      severity = 2
    }
  }
}

resource "circonus_rule_set_group" "icmp_latency_1" {
  name = "icmp latency group 1"

  notify {
    sev1 = [
      "/contact_group/4680"
    ]
    sev2 = [
      "/contact_group/4680"
    ]
    sev3 = [
      "/contact_group/4680"
    ]
  }

  formula {
    expression = "A and B"
    raise_severity = 3
    wait = 0
  }

  condition {
    index = 1
    rule_set = circonus_rule_set.icmp_max_latency.id
    matching_severities = ["3"]
  }

  condition {
    index = 2
    rule_set = circonus_rule_set.icmp_min_latency.id
    matching_severities = ["3"]
  }

  condition {
    index = 3
    rule_set = circonus_rule_set.icmp_avg_latency.id
    matching_severities = ["3"]
  }

}

resource "circonus_rule_set_group" "icmp_latency_2" {
  name = "icmp latency group 2"

  notify {
    sev1 = [
      "/contact_group/4680"
    ]
    sev2 = [
      "/contact_group/4680"
    ]
    sev3 = [
      "/contact_group/4680"
    ]
  }

  formula {
    expression = 2
    raise_severity = 2
    wait = 1
  }

  condition {
    index = 1
    rule_set = circonus_rule_set.icmp_max_latency.id
    matching_severities = ["3"]
  }

  condition {
    index = 2
    rule_set = circonus_rule_set.icmp_min_latency.id
    matching_severities = ["3"]
  }

  condition {
    index = 3
    rule_set = circonus_rule_set.icmp_avg_latency.id
    matching_severities = ["3"]
  }
}
`
