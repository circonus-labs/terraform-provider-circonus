---
layout: "circonus"
page_title: "Circonus: circonus_dashboard"
sidebar_current: "docs-circonus-resource-circonus_dashboard"
description: |-
  Manages a Circonus dashboard.
---

# circonus\_dashboard

The ``circonus_dashboard`` resource creates and manages a
[Circonus Dashboard](https://login.circonus.com/resources/docs/user/Dashboards.html).

https://login.circonus.com/resources/api/calls/dashboard.

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

resource "circonus_dashboard" "latency-dash" {
  title = "Latency Dashboard"
  shared = true
  account_default = false
  grid_layout = {
    width = 1
    height = 2
  }
  
  options {
    hide_grid = true
    text_size = 14
  }

  widget {
    active = true
    height = 1
    width = 1
    origin = "a0"
    type = "html"
    name = "HTML"
    widget_id = "w0"
    settings {
      markup = <<EOF
      <p>Some rando HTML</p>
EOF
    }
  }

  widget {
    active = true
    height = 1
    width = 1
    origin = "a1"
    type = "graph"
    widget_id = "w1"
    name = "Graph"

    settings {
      date_window = "global"
      graph_uuid = element(split("/", circonus_graph.latency_graph.id),2)
      show_flags = true
    }
  }
}

```

## Argument Reference

* `title` - (Required) String. The title of the dashboard

* `shared` - (Optional) Boolean. Whether this dash is shared with everyone in the 
  Circonus account or private to the user.
  
* `account_default` - (Optional) Boolean. Whether this dash is set as the default
  for the account
  
* `grid_layout` - (Required) Map. Has members `width` and `height` to control the 
  dimensions of the dashboard.
  
* `options` - (Required) Set of 1. The options for the dashboard, see below for description.

* `widget` - (Required) Set of N. The widgets that make up the dashboard.  See below
  for the description

## `options` Configuration

These options control rendering and other global settings about the dashboard.

* `full_screen_hide_title` - (Optional) Boolean.  Set the dashboard to fullscreen mode.

* `hide_grid` - (Optional) Boolean.  Hides the grid lines when rendering the dashboard.

* `scale_text` - (Optional) Boolean.  Scale the text up and down with the count of widgets.
  Defaults to `true`.
  
* `text_size` - (Optional) Integer.  The point size of the text in the dashboard (titles and such)
  Defaults to `12`.
  
## `widget` Configuration

* `active` - (Optional) Boolean.  Is this widget active.  Default to true

* `height` - (Required) Integer.  The number of dashboard rows this widget consumes.

* `width` - (Required) Integer.  The number of dashboard columns this widget consumes.

* `name` - (Required) String.  The name of the widget.  One of: `Graph`, `Chart`, `Gauge`,
  `Text`, `Cluster`, `HTML`, `Status`, `List`, `Alerts`, `Forecast`.
  
* `type` - (Required) String.  The type of the widget.  This is duplicative of `name` but must be 
  one of: `graph`, `chart`, `gauge`, `text`, `cluster`, `html`, `status`, `list`, `alerts`, `forecast`.

* `origin` - (Required) String. The cell ID of the left top corner of this widget. Where cell ID
  is the columnar letter (a->z lower case) followed by by the row number (zero indexed).  "a0" would
  mean first column, first row (top-left most cell).  "c3" would mean the 3rd column, 4th row.
  
* `settings` - (Required) Set. The settings specific to this widget type.  See below for descriptions.

* `widget_id` - (Required) String. Widget ID for this widget. Must be unique to this dashboard. 
  A common convention is to use 'w' + integer.
  
## `widget.settings` Configuration

These settings are all optional in practice but some are marked Required here because they
are required by the Circonus API.

### `Alert` type settings

* `account_id` - (Optional) String.  To access alert widgets from other circonus accounts supply the account_id here.
* `acknowledged` - (Optional) String.  Show ack'd alerts. "y", "n", or "all"
* `cleared` - (Optional) String.  Show cleared alerts. "y", "n", or "all"
* `dependents` - (Optional) String.  Show dependent alerts. "y", "n", or "all"
* `display` - (Optional) String.  Method to use when displaying. "list", "bar", or "sunburst"
* `maintenance` - (Optional) String.  Show alerts in maintenance. "y", "n", or "all"
* `min_age` - (Optional) String.  Show alerts of a certain age. A string matching the following regex: "(?:0|\d+[mhdwMy])"
* `search` - (Optional) String.  A search string to use when displaying alerts.
* `severity` - (Optional) String.  A string matching the following regex: "[1-5]{1,5}"
* `time_window` - (Optional) String.  A string matching the following regex: "(?:0|\d+[mhdwMy])"
* `title` - (Optional) String.  The title of the widget.

### `Chart` type settings

* `title` - (Optional) String.  The title of the widget.
* `chart_type` - (Optional) String.  One of: "pie", "bar"
* `datapoints` - (Optional) Set of N. Consisting of:

  * `_metric_type` - (Required) String. 
  * `_check_id` - (Required) Integer.
  * `label` - (Required) String.
  * `metric` - (Required) String.
  
### `Graph` type settings

* `title` - (Optional) String.  The title of the widget.
* `date_window` - (Optional) String. 'global' (follow the page datetool settings) | 
  <time_interval> (e.g. '30m', '6h', '2d', '1w', etc.) | <dual_time_intervals> (e.g. '6h:12h', '1w:1w', etc.)
* `graph_uuid` - (Required) String.  The uuid of the graph.
* `hide_xaxis` - (Optional) Boolean.  Whether to hide the x-axis labels.
* `hide_yaxis` - (Optional) Boolean.  Whether to hide the y-axis labels.
* `key_inline` - (Optional) Boolean.  Whether to show the legend when hovering.
* `key_loc` - (Optional) String.  Whether to show a persistent legend key beside the graph and where 
  to position it. 'noop' (don't show key) - 'in-tl' (inside graph, top left corner) 
  - 'out-r' (outside graph, on the right) - 'out-b' (outside graph, on the bottom) 
* `key_size` - (Optional) Integer.  The size of the persistent legend key (if `key_loc` is 'in-tl' 
  or 'out-r', this controls key width; if `key_loc` is 'out-b', this controls key height)
* `key_wrap` - (Optional) Boolean.  Whether to wrap text within the persistent legend key.
* `label` - (Optional) String.  The label to show at the top of the widget if you want it to show something 
  other than the graph title
* `period` - (Optional) Integer.  Realtime streaming update period, in milliseconds
* `real_time` - (Optional) Boolean.  Whether to plot streaming data in realtime instead of showing recent stored data
* `show_flags` - (Optional) Boolean.  Whether to show the legend upon mouse hover












  
  










  


