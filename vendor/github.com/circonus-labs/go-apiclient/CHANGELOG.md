# v0.6.3

* upd: remove tests for invalid cids
* fix: validate cids on prefix only to compensate for breaking change to rule_set cid in public v2 api

# v0.6.2

* upd: dependencies (retryablehttp)

# v0.6.1

* add: full overlay test suite to `examples/graph/overlays`
* fix: incorrect attribute types in graph overlays (docs vs what api actually returns)

# v0.6.0

* fix: graph structures incorrectly represented nesting of overlay sets

# v0.5.4

* add: `search` (`*string`) attribute to graph datapoint
* upd: `cluster_ip` (`*string`) can be string OR null
* add: `cluster_ip` attribute to broker details

# v0.5.3

* upd: use std log for retryablehttp until dependency releases Logger interface

# v0.5.2

* upd: support any logging package with a `Printf` method via `Logger` interface rather than forcing `log.Logger` from standard log package
* upd: remove explicit log level classifications from logging messages
* upd: switch to errors package (for `errors.Wrap` et al.)
* upd: clarify error messages
* upd: refactor tests
* fix: `SearchCheckBundles` to use `*SearchFilterType` as its second argument
* fix: remove `NewAlert` - not applicable, alerts are not created via the API
* add: ensure all `Delete*ByCID` methods have CID corrections so short CIDs can be passed

# v0.5.1

* upd: retryablehttp to start using versions that are now available instead of tracking master

# v0.5.0

* Initial - promoted from github.com/circonus-labs/circonus-gometrics/api to an independant package
