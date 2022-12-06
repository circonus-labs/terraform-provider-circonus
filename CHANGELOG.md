## 0.12.14 (December 6, 2022)

CHANGES:

* upd: Updates the behavior of the HTTP check resource to scan through the
list of check type interfaces to find the correct configuration values.

## 0.12.13 (November 16, 2022)

CHANGES:

* upd: Updates the behavior of the graph resource to allow values removed
from terraform files to be correctly cleared by the terraform provider.

## 0.12.10 (June 8, 2022)

CHANGES:

* upd: Removes support from the circonus_contact_group resource for external
contact methods that are no longer supported: irc and xmpp. Updates the
documentation to remove these methods as well.

## 0.12.9 (June 6, 2022)

CHANGES:

* add: Adds the `default_dashboard_uuid` and `default_dashboard_type`
attributes to the Account data source schema to support the new changes to
the dashboard defaulting process.
* upd: Removes the `account_default` attribute from the Dashboard resource
schema to reflect the changes to the dashboard defaulting process since this
field is no longer returned by the API.

UPDATES:

* upd: Updates the version of go-apiclient to v0.7.16.

## 0.12.8 (June 6, 2022)

UPDATES:

* upd: Bumps github.com/hashicorp/terraform-plugin-sdk/v2 from 2.16.0 to 2.17.0.

## 0.12.7 (June 1, 2022)

IMPROVEMENTS:

* add: Added support for the ssh2 check type.

## 0.12.4 (November 22, 2021)

CHANGES:

* fix: mismatch in http contact types (not reading format field and normalize url field)
* upd: Added environment variable settings to the acceptance test suite so that
it can be run in other environments.

## 0.12.3 (October 6, 2021)

CHANGES:

* upd: deprecate tags on rule_set (they are just check tags, api ignores tags on rule_set objects)
* upd: using `over` on rule_set REQUIRES all three attributes to be set

UPDATES:

* upd: terraform-plugin-sdk v2.8.0
* add: rule_set_group test
* upd: use diagnostics instead of errors in several resources (ongoing updates)
* upd: use *Context methods in several resources (ongoing updates)

## 0.12.2 (September 15, 2021)

CHANGES:

* fix: rule set group expression can be string or numeric,
  use interface for api, convert to string for state file

## 0.12.1 (September 15, 2021)

CHANGES:

* fix: rule set group expression can be string or numeric, use interface
* upd: var for update ruleset test
* upd: var names for graph update test
* upd: api debug log if TF_LOG DEBUG or TRACE
* upd: lint errors/warnings
* add: lint config

UPDATES:

* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.6.1 to 2.7.0
* build(deps): bump github.com/circonus-labs/go-apiclient from 0.7.14 to 0.7.15
* build(deps): bump github.com/circonus-labs/go-apiclient from 0.7.13 to 0.7.14
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.5.0 to 2.6.1
* build(deps): bump github.com/circonus-labs/go-apiclient from 0.7.11 to 0.7.13
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.4.4 to 2.5.0
* build(deps): bump github.com/hashicorp/terraform-plugin-sdk/v2 from 2.4.3 to 2.4.4

## 0.12.0 (March 19, 2021)

CHANGES:

* upd: terraform-plugin-sdk v2 (Terraform v0.12.0+ required)

FIXES:

* fix: regex for ruleset.metric_name, metric names can indeed have spaces

## 0.11.8 (March 5, 2021)

FIXES:

* fix spurious snmp diffs, test for snmp checks

CHANGES:

* implement new ruleset.user_json field

## 0.11.7 (February 17, 2021)

FIXES:

* fix ruleset.parent regex to allow old form and new form

## 0.11.6 (February 16, 2021)

CHANGES:

* add support for eq_value and neq_value for numeric rulesets

## 0.11.5 (February 3, 2021)

CHANGES:

* add Name field to ruleset resource

FIXES:

* resolve ruleset notify list ordering issue
* deal with `if.then.after == 0`, deal with `severity = 0` notify lists
* spurious diffs in ruleset
* fix clear ruleset rules (severity 0)
* fixes for metric.caql whitespace causing spurious diffs
* merge ntp test fixes from ntp_check branch
* fix ruleset absence test with type switch since it can be string or number
* fix graph test, make check name and graph name static across create/modification tests

## 0.11.4 (Novermber 19, 2020)

FIXES:

* http check updating gets two configs in []interface{} list - the first is the valid updated config and the second is empty. This results in the check bundle Config being overwritten with blank values for each attribute, then the API complains about missing attributes.

## 0.11.3 (October 28, 2020)

CHANGES:

* upgrade dependencies (go-apiclient,retryablehttp) to use Retry-After header on 429 API responses

## 0.11.2 (September 2, 2020)

FIXES:

* add http check `redirects` attribute

CHANGES:

* deprecate `irc` contact type

## 0.11.1 (August 05, 2020)

FIXES:

* use `{{ .Tag }}` in binary name to get v prefixed version in binary file name using `{{ .Version }}` resulted in `x.y.z`; using `v{{ .Version }}` resulted in `vvx.y.z`

## 0.11.0 (August 05, 2020)

IMPROVEMENTS:

* add NTP check support
* add SMTP check support

## 0.10.1 (August 04, 2020)

* transfer repository to circonus-labs
* add goreleaser configuration
* initial build/release cycle

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
  * dns

IMPROVEMENTS:

* Circonus go-apiclient 0.7.3
* Change `circonus_check.metric_filters` a List to preserve order of filters
* Add `state` widget support to the `circonus_dashboard`

## 0.7.0 (February 05, 2020)

FEATURES:

* New support for check types:
  * redis

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
  * external
  * jmx
  * memcached
  * promtext
  * snmp

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
