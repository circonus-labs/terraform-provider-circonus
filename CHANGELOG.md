## 0.1.1 (September 19, 2018)

BUG FIXES:

* `resource/circonus_rule_set`: Change the default `severity` to `0` to allow clearing of an alert ([#15](https://github.com/terraform-providers/terraform-provider-circonus/pull/15))
* `resource/circonus_rule_set`: Changes to the `check` or `metric_name` will result in a new ruleset being created ([#10](https://github.com/terraform-providers/terraform-provider-circonus/pull/10))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
