---
layout: "circonus"
page_title: "Circonus: circonus_rule_set_group"
sidebar_current: "docs-circonus-resource-circonus_rule_set_group"
description: |-
  Manages a Circonus rule set group.
---

# circonus\_rule\_set_group

The ``circonus_rule_set_group`` resource creates and manages a
[Circonus Rule Set Group](https://login.circonus.com/resources/api/calls/rule_set_group).

## Usage

```hcl
variable "myapp-tags" {
  type    = "list"
  default = [ "app:myapp", "owner:myteam" ]
}

resource "circonus_rule_set_group" "myapp-ruleset-group" {
  name = "Test Rule Set Group 1"
  notify {
    sev3 = [
      "${circonus_contact_group.myapp-owners.id}"
    ]
  }

  formula {
    expression = "A and B"
    raise_severity = 4
    wait = 0
  }

  condition {
    index = 1
    rule_set = circonus_rule_set.myapp-cert-ttl-alert.id
    matching_severities = ["3"]
  }

  condition {
    index = 2
    rule_set = circonus_rule_set.myapp-healthy-alert.id
    matching_severities = ["3"]
  }
}

resource "circonus_rule_set" "myapp-cert-ttl-alert" {
  check       = "${circonus_check.myapp-https.checks[0]}"
  metric_name = "cert_end_in"
  link        = "https://wiki.example.org/playbook/how-to-renew-cert"

  if {
    value {
      min_value = "${2 * 24 * 3600}"
    }

    then {
      notify = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 1
    }
  }

  if {
    value {
      min_value = "${7 * 24 * 3600}"
    }

    then {
      notify = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 2
    }
  }

  if {
    value {
      min_value = "${21 * 24 * 3600}"
    }

    then {
      notify = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 3
    }
  }

  if {
    value {
      absent = "24h"
    }

    then {
      notify = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 1
    }
  }

  tags = [ "${var.myapp-tags}" ]
}

resource "circonus_rule_set" "myapp-healthy-alert" {
  check = "${circonus_check.myapp-https.checks[0]}"
  metric_name = "duration"
  link = "https://wiki.example.org/playbook/debug-down-app"

  if {
    value {
      # SEV1 if it takes more than 9.5s for us to complete an HTTP request
      max_value = "${9.5 * 1000}"
    }

    then {
      notify = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 1
    }
  }

  if {
    value {
      # SEV2 if it takes more than 5s for us to complete an HTTP request
      max_value = "${5 * 1000}"
    }

    then {
      notify = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 2
    }
  }

  if {
    value {
      # SEV3 if the average response time is more than 500ms using a moving
      # average over the last 10min.  Any transient problems should have
      # resolved themselves by now.  Something's wrong, need to page someone.
      over {
        last  = "10m"
        using = "average"
      }
      max_value = "500"
    }

    then {
      notify = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 3
    }
  }

  if {
    value {
      # SEV4 if it takes more than 500ms for us to complete an HTTP request.  We
      # want to record that things were slow, but not wake anyone up if it
      # momentarily pops above 500ms.
      min_value = "500"
    }

    then {
      notify   = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 3
    }
  }

  if {
    value {
      # If for whatever reason we're not recording any values for the last
      # 24hrs, fire off a SEV1.
      absent = "24h"
    }

    then {
      notify = [ "${circonus_contact_group.myapp-owners.id}" ]
      severity = 1
    }
  }

  tags = [ "${var.myapp-tags}" ]
}

resource "circonus_contact_group" "myapp-owners" {
  name = "My App Owners"
  tags = [ "${var.myapp-tags}" ]
}

resource "circonus_check" "myapp-https" {
  name = "My App's HTTPS Check"

  notes = <<-EOF
A check to create metric streams for Time to First Byte, HTTP transaction
duration, and the TTL of a TLS cert.
EOF

  collector {
    id = "/broker/1"
  }

  http {
    code = "^200$"
    headers = {
      X-Request-Type = "health-check",
    }
    url = "https://www.example.com/myapp/healthz"
  }

  metric {
    name = "${circonus_metric.myapp-cert-ttl.name}"
    tags = "${circonus_metric.myapp-cert-ttl.tags}"
    type = "${circonus_metric.myapp-cert-ttl.type}"
    unit = "${circonus_metric.myapp-cert-ttl.unit}"
  }

  metric {
    name = "${circonus_metric.myapp-duration.name}"
    tags = "${circonus_metric.myapp-duration.tags}"
    type = "${circonus_metric.myapp-duration.type}"
    unit = "${circonus_metric.myapp-duration.unit}"
  }

  period       = 60
  tags         = ["source:circonus", "author:terraform"]
  timeout      = 10
}

resource "circonus_metric" "myapp-cert-ttl" {
  name = "cert_end_in"
  type = "numeric"
  unit = "seconds"
  tags = [ "${var.myapp-tags}", "resource:tls" ]
}

resource "circonus_metric" "myapp-duration" {
  name = "duration"
  type = "numeric"
  unit = "miliseconds"
  tags = [ "${var.myapp-tags}" ]
}
```

## Argument Reference

* `name` - (Required) The name of the rule set group, must be unique across
  all rule set groups within the account.

* `notify` - (Required) The list of contact groups to notify should
  the expression evaluate to true.  See below for details on the
  structure of a `notify` configuration clause.

* `formula` - (Required) Instructions for how to compare the member rule sets
  for trigger notification.  See below for details on the
  structure of a `formula` configuration clause.

* `condition` - (Required) The rule set reference and condition levels to watch.  
  See below for details on the structure of a `condition` configuration clause.

## `notify` Configuration

The `notify` configuration block is a listing of contact groups separated by severity
that should be notified when the `formula` evaluates to true based on the `raise_severity`
in the `formula`.

There are 5 severity levels supported and each one is a list of strings.  Each list contains
the CID of the contact group to notify at that severity level.

`sev1`, `sev2`, `sev3`, `sev4`, `sev5`: are the names of the `notify` lists.

### `formula` Configuration

A `formula` block contains 3 fields that indicate how to combine the member rule sets specified
in the `condition` blocks.  There can be only 1 `formula` block.  It has 3 fields:

* `expression` - (Required) The expression that combines the `condition`s into a boolean logical expression.
  See Formulas [here](https://login.circonus.com/resources/docs/user/Alerting/RuleGroups/Configure.html)
* `raise_severity` - (Required) The severity level to raise (see `notify` for who would be contacted), when
  the `expression` is true.
* `wait` - (Required) How long to wait before sending out the alert.

### `condition` Configuration

A `condition` block contains 3 fields that indicate what rule set, it's position in the expression, and
which severities of the original ruleset to pay attention to.  It has 3 fields:

* `index` - (Required) The position this condition has in the `formula`.`expression`.  A value of `1` maps
  to `A`, a value of `2` maps to `B`, etc..
* `rule_set` - (Required) The CID of the rule set to pay attention to.
* `matching_severities` - (Required) The list(string) of severities from that rule set to watch.


## Import Example

`circonus_rule_set_group` supports importing resources.  Supposing the following
Terraform (and that the referenced [`circonus_rule_set`](rule_set.html)s and [`circonus_contact_group`](contact_group.html)
have already been imported):

```hcl
resource "circonus_rule_set_group" "myrulesetgroup" {
  name = "Test Rule Set Group 1"
  notify {
    sev3 = [
      "${circonus_contact_group.myapp-owners.id}"
    ]
  }

  formula {
    expression = "A and B"
    raise_severity = 4
    wait = 0
  }

  condition {
    index = 1
    rule_set = circonus_rule_set.myapp-cert-ttl-alert.id
    matching_severities = ["3"]
  }

  condition {
    index = 2
    rule_set = circonus_rule_set.myapp-healthy-alert.id
    matching_severities = ["3"]
  }
}
```

It is possible to import a `circonus_rule_set_group` resource with the following command:

```
$ terraform import circonus_rule_set_group.myrulesetgroup ID
```

Where `ID` is the `_cid` or Circonus ID of the Rule Set Group
(e.g. `/rule_set_group/201285`) and `circonus_rule_set_group.myrulesetgroup` is
the name of the resource whose state will be populated as a result of the
command.
