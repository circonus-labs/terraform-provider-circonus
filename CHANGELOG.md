## 0.3.0 (Unreleased)

IMPROVEMENTS:

* Switch from deprecated circonus-gometrics api sub-package to the current [go-apiclient](https://github.com/circonus-labs/go-apiclient) package

BUG FIXES:

* Fix CAQL datapoint support in graphs
* Incorporate fix in [go-apiclient v0.6.3](https://github.com/circonus-labs/go-apiclient/releases/tag/v0.6.3) to address breaking change to `rule_set` CID format in public API
* Incorporate fix in [go-apiclient v0.6.4](https://github.com/circonus-labs/go-apiclient/releases/tag/v0.6.4) graph.datapoint.alpha - doc:floating point number, api:string

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
