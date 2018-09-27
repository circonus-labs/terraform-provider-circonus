---
layout: "circonus"
page_title: "Circonus: circonus_worksheet"
sidebar_current: "docs-circonus-resource-circonus_worksheet"
description: |-
  Manages a Circonus worksheet.
---

# circonus_worksheet

The ``circonus_worksheet`` resource creates and manages a
[Circonus Worksheet](https://login.circonus.com/resources/api/calls/worksheet).

## Usage

```hcl
variable "myapp-tags" {
  type    = "list"
  default = [ "app:myapp", "owner:myteam" ]
}

resource "circonus_graph" "latency-graph" {
  name        = "Latency Graph"
  description = "A sample graph showing off two data points"
  notes       = "Misc notes about this graph"
  graph_style = "line"
  line_style  = "stepped"

  metric {
    check       = "${circonus_check.api_latency.checks[0]}"
    metric_name = "maximum"
    metric_type = "numeric"
    name        = "Maximum Latency"
    axis        = "left"
    color       = "#657aa6"
  }

  metric {
    check       = "${circonus_check.api_latency.checks[0]}"
    metric_name = "minimum"
    metric_type = "numeric"
    name        = "Minimum Latency"
    axis        = "right"
    color       = "#0000ff"
  }

  tags = [ "${var.myapp-tags}" ]
}

resource "circonus_worksheet" "latency_worksheet" {
  title = "%s"
  graphs = [
    "${circonus_graph.latency-graph.id}",
  ]
}
```

## Argument Reference

* `title` - (Required) The title of the worksheet.

* `description` - (Optional) Description of what the worksheet is for.

* `favourite` - (Optional) Mark (star) this worksheet as a favorite. Default is `false`.

* `notes` - (Optional) A place to store notes about this worksheet.

* `graphs` - (Optional) A list of graphs that compose this worksheet.

* `smart_queries` - (Optional) The smart queries that will be displayed on this worksheet. See below for details on how to configure a `smart_query`.

* `tags` - (Optional) A list of tags assigned to this worksheet.

### `smart_queries` Attributes

* `name` - (Required) The name (heading) for the smart graph section in the worksheet.

* `query` - (Required) A search query that determines which graphs will be shown..

## Import Example

It is possible to import a `circonus_worksheet` resource with the following command:

```
$ terraform import circonus_worksheet.icmp-latency ID
```

Where `ID` is the `_cid` or Circonus ID of the worksheet
(e.g. `worksheets/45640239-bb81-4ecb-81e6-b5c6015e5dd5`) and `circonus_worksheet.icmp-latency` is
the name of the resource whose state will be populated as a result of the
command.
