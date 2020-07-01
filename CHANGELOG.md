## 0.11.0 (Unreleased)
## 0.10.0 (July 01, 2020)

IMPROVEMENTS:

* go-apiclient v0.7.7
* Environment vars to external checks 
* Support for new windowing_min_duration in rulesets

FIXES:

* update ruleset parent pattern
* remove support for deprecated metric.tags and metric.units (use streamtags)
* remove support for deprecated metric_cluster

## 0.9.0 (April 03, 2020)

IMPROVEMENTS:

* go-apiclient 0.7.6
* Fix to prevent noop updates to SNMP checks 
* The metrics in a check should be a list, not a set, to preserve order (noop updates again)
* Add a test for metric_filters in a check
* Add validation for a ruleset to prevent metric_type -> Rule.Criteria mismatches
* Prevent unnecessary updates to dashboards
* Do exponential backoff when API returns 503's
* Fix slack contact groups to have reasonable slack default text
* Fix external check output extract: JSON, NAGIOS, otherwise treated as regexp


## 0.8.0 (February 24, 2020)

FEATURES:

* New support for check types:
  - dns

IMPROVEMENTS:

* Circonus go-apiclient 0.7.3
* Change `circonus_check.metric_filters` a List to preserve order of filters
* Add `state` widget support to the `circonus_dashboard`

## 0.7.0 (February 05, 2020)

FEATURES:

* New support for check types:
  - redis

IMPROVEMENTS:

* Change to `circonus_rule_set` to eliminate false differences due to time conversion.
  Now all times in a `circonus_rule_set` are seconds.

## 0.6.0 (January 30, 2020)

* **New Resource:** `circonus_maintenance` 

IMPROVEMENTS:

* Small fixes for ruleset processing

## 0.5.0 (January 23, 2020)

FEATURES:

* **New Resource:** `circonus_overlay_set` 
* **New Resource:** `circonus_rule_set_group` 
* **New Resource:** `circonus_dashboard` 

* New support for check types:
  - external
  - jmx
  - memcached
  - promtext
  - snmp

IMPROVEMENTS:

* Circonus go-apiclient 0.7.2
* Support for Metric allow/deny filters in Checks
* Support for guide lines in graphs
* Support for search based datapoints graphs
* Support for pattern based Rule sets
* Support for query order in worksheets

## 0.4.0 (December 02, 2019)

IMPROVEMENTS:

* Update dependencies
* Switch to Terraform Plugin SDK

NOTES:

* Minimum version of Go required to build the provider is now 1.13

## 0.3.0 (November 14, 2019)

IMPROVEMENTS:

* Provider: Migrate from deprecated [circonus-gometrics/api](https://github.com/circonus-labs/circonus-gometrics) to [go-apiclient](https://github.com/circonus-labs/go-apiclient)
* Provider: Support new attributes for contact group (`group_type` and `always_send_clear`)

## 0.2.0 (October 01, 2018)

FEATURES:

* **New Resource:** `circonus_worksheet` ([#17](https://github.com/terraform-providers/terraform-provider-circonus/pull/17))

IMPROVEMENTS:

* Provider: Accept the `CIRCONUS_API_URL` environment variable to configure the API URL ([#18](https://github.com/terraform-providers/terraform-provider-circonus/pull/18))
* Provider: Upgrade `circonus-gometrics` to `2.2.4` ([#22](https://github.com/terraform-providers/terraform-provider-circonus/pull/22))

NOTES:

* Minimum version of Go required to build the provider is now 1.10
* Deprecated `govendor` in favor of using `dep` for vendor management ([#19](https://github.com/terraform-providers/terraform-provider-circonus/pull/19))

## 0.1.1 (September 19, 2018)

BUG FIXES:

* `resource/circonus_rule_set`: Change the default `severity` to `0` to allow clearing of an alert ([#15](https://github.com/terraform-providers/terraform-provider-circonus/pull/15))
* `resource/circonus_rule_set`: Changes to the `check` or `metric_name` will result in a new ruleset being created ([#10](https://github.com/terraform-providers/terraform-provider-circonus/pull/10))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
