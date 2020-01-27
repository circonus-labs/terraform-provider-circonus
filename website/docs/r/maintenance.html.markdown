---
layout: "circonus"
page_title: "Circonus: circonus_maintenance"
sidebar_current: "docs-circonus-resource-circonus_maintenance"
description: |-
  Manages a Circonus maintenance window.
---

# circonus\_maintenance

The ``circonus_maintenance`` resource creates and manages a
single [maintenance window resource](https://login.circonus.com/resources/docs/user/Alerting/Maintenance.html)
which can reference a check, ruleset, target or entire account.

## Usage

```hcl
resource "circonus_maintenance" "turn_off_account" {
  account  = "/account/12345"
  notes = "shutting sev1 alerts for the account for the weekend"
  severities = ["1"]
  start = "2020-01-25T19:00:00-05:00"
  stop = "2020-01-27T19:00:00-05:00"
  tags = {
    author = "terraform"
    source = "circonus"
  }
}
```

## Argument Reference

* `account` - (Optional) A string referencing the account CID to have maintenance on, mutually exclusive 
  with `check`, `rule_set`, and `target`.

* `check` - (Optional) A string referencing the check CID to have maintenance on, mutually exclusive 
  with `account`, `rule_set`, and `target`.

* `rule_set` - (Optional) A string referencing the rule_set CID to have maintenance on, mutually exclusive 
  with `account`, `check`, and `target`.
  
* `target` - (Optional) A string referencing the check target (host) to have maintenance on, mutually exclusive 
  with `account`, `rule_set`, and `check`.
  
* `severities` - (Required) A list of strings determining which severities to put into maintenance.  
  Must be in the range: "1"-"5"
  
* `start` - (Required) An RFC3339 timestamp string which indicates the start of the maintenance window.

* `stop` - (Required) An RFC3339 timestamp string which indicates the end of the maintenance window.
  
* `tags` - (Optional) A list of tags assigned to the maintenance window.

## Import Example

`circonus_maintenance` supports importing resources.  Supposing the following
Terraform:

```hcl
provider "circonus" {
  alias = "b8fec159-f9e5-4fe6-ad2c-dc1ec6751586"
}

resource "circonus_maintenance" "mine" {
  account  = "/account/12345"
  notes = "shutting sev1 alerts for the account for the weekend"
  severities = ["1"]
  start = "2020-01-25T19:00:00-05:00"
  stop = "2020-01-27T19:00:00-05:00"
  tags = {
    author = "terraform"
    source = "circonus"
  }
}
```

It is possible to import a `circonus_maintenance` resource with the following command:

```
$ terraform import circonus_maintenance.mine ID
```

Where `ID` is the CID of the matching maintenance window.
