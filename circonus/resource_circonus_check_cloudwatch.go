package circonus

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/circonus-labs/terraform-provider-circonus/internal/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// circonus_check.cloudwatch.* resource attribute names.
	checkCloudWatchAPIKeyAttr      = "api_key"
	checkCloudWatchAPISecretAttr   = "api_secret"
	checkCloudWatchDimmensionsAttr = "dimmensions"
	checkCloudWatchMetricAttr      = "metric"
	checkCloudWatchNamespaceAttr   = "namespace"
	checkCloudWatchURLAttr         = "url"
	checkCloudWatchVersionAttr     = "version"
)

var checkCloudWatchDescriptions = attrDescrs{
	checkCloudWatchAPIKeyAttr:      "The AWS API Key",
	checkCloudWatchAPISecretAttr:   "The AWS API Secret",
	checkCloudWatchDimmensionsAttr: "The dimensions to query for the metric",
	checkCloudWatchMetricAttr:      "One or more CloudWatch Metric attributes",
	checkCloudWatchNamespaceAttr:   "The namespace to pull telemetry from",
	checkCloudWatchURLAttr:         "The URL including schema and hostname for the Cloudwatch monitoring server. This value will be used to specify the region - for example, to pull from us-east-1, the URL would be https://monitoring.us-east-1.amazonaws.com.",
	checkCloudWatchVersionAttr:     "The version of the Cloudwatch API to use.",
}

var schemaCheckCloudWatch = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      hashCheckCloudWatch,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkCloudWatchDescriptions, map[schemaAttr]*schema.Schema{
			checkCloudWatchAPIKeyAttr: {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateRegexp(checkCloudWatchAPIKeyAttr, `[\S]+`),
				DefaultFunc:  schema.EnvDefaultFunc("AWS_ACCESS_KEY_ID", ""),
			},
			checkCloudWatchAPISecretAttr: {
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateRegexp(checkCloudWatchAPISecretAttr, `[\S]+`),
				DefaultFunc:  schema.EnvDefaultFunc("AWS_SECRET_ACCESS_KEY", ""),
			},
			checkCloudWatchDimmensionsAttr: {
				Type:         schema.TypeMap,
				Required:     true,
				Elem:         schema.TypeString,
				ValidateFunc: validateCheckCloudWatchDimmensions,
			},
			checkCloudWatchMetricAttr: {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Set:      schema.HashString,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateRegexp(checkCloudWatchMetricAttr, `^([\S]+)$`),
				},
			},
			checkCloudWatchNamespaceAttr: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp(checkCloudWatchNamespaceAttr, `.+`),
			},
			checkCloudWatchURLAttr: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateHTTPURL(checkCloudWatchURLAttr, urlIsAbs),
			},
			checkCloudWatchVersionAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      defaultCheckCloudWatchVersion,
				ValidateFunc: validateRegexp(checkCloudWatchVersionAttr, `^[\d]{4}-[\d]{2}-[\d]{2}$`),
			},
		}),
	},
}

// checkAPIToStateCloudWatch reads the Config data out of circonusCheck.CheckBundle into the
// statefile.
func checkAPIToStateCloudWatch(c *circonusCheck, d *schema.ResourceData) error {
	cloudwatchConfig := make(map[string]interface{}, len(c.Config))

	// swamp is a sanity check: it must be empty by the time this method returns
	swamp := make(map[config.Key]string, len(c.Config))
	for k, v := range c.Config {
		swamp[k] = v
	}

	saveStringConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if v, ok := c.Config[apiKey]; ok {
			cloudwatchConfig[string(attrName)] = v
		}

		delete(swamp, apiKey)
	}

	saveStringConfigToState(config.APIKey, checkCloudWatchAPIKeyAttr)
	saveStringConfigToState(config.APISecret, checkCloudWatchAPISecretAttr)

	dimmensions := make(map[string]interface{}, len(c.Config))
	dimmensionPrefixLen := len(config.DimPrefix)
	for k, v := range c.Config {
		if len(k) <= dimmensionPrefixLen {
			continue
		}

		if strings.Compare(string(k[:dimmensionPrefixLen]), string(config.DimPrefix)) == 0 {
			key := k[dimmensionPrefixLen:]
			dimmensions[string(key)] = v
		}
		delete(swamp, k)
	}
	cloudwatchConfig[string(checkCloudWatchDimmensionsAttr)] = dimmensions

	metricSet := schema.NewSet(schema.HashString, nil)
	metricList := strings.Split(c.Config[config.CloudwatchMetrics], ",")
	for _, m := range metricList {
		metricSet.Add(m)
	}
	cloudwatchConfig[string(checkCloudWatchMetricAttr)] = metricSet

	saveStringConfigToState(config.Namespace, checkCloudWatchNamespaceAttr)
	saveStringConfigToState(config.URL, checkCloudWatchURLAttr)
	saveStringConfigToState(config.Version, checkCloudWatchVersionAttr)

	whitelistedConfigKeys := map[config.Key]struct{}{
		config.ReverseSecretKey: {},
		config.SubmissionURL:    {},
	}

	for k := range swamp {
		if _, ok := whitelistedConfigKeys[k]; ok {
			delete(c.Config, k)
		}

		if _, ok := whitelistedConfigKeys[k]; !ok {
			log.Printf("[ERROR]: PROVIDER BUG: API Config not empty: %#v", swamp)
		}
	}

	if err := d.Set(checkCloudWatchAttr, schema.NewSet(hashCheckCloudWatch, []interface{}{cloudwatchConfig})); err != nil {
		return fmt.Errorf("Unable to store check %q attribute: %w", checkCloudWatchAttr, err)
	}

	return nil
}

// hashCheckCloudWatch creates a stable hash of the normalized values.
func hashCheckCloudWatch(v interface{}) int {
	m := v.(map[string]interface{})
	b := &bytes.Buffer{}
	b.Grow(defaultHashBufSize)

	writeString := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok && v.(string) != "" {
			fmt.Fprint(b, strings.TrimSpace(v.(string)))
		}
	}

	// Order writes to the buffer using lexically sorted list for easy visual
	// reconciliation with other lists.
	writeString(checkCloudWatchAPIKeyAttr)
	writeString(checkCloudWatchAPISecretAttr)

	if dimmensionsRaw, ok := m[string(checkCloudWatchDimmensionsAttr)]; ok {
		dimmensionMap := dimmensionsRaw.(map[string]interface{})
		dimmensions := make([]string, 0, len(dimmensionMap))
		for k := range dimmensionMap {
			dimmensions = append(dimmensions, k)
		}

		sort.Strings(dimmensions)
		for i := range dimmensions {
			fmt.Fprint(b, dimmensions[i])
		}
	}

	if metricsRaw, ok := m[string(checkCloudWatchMetricAttr)]; ok {
		metricListRaw := flattenSet(metricsRaw.(*schema.Set))
		for i := range metricListRaw {
			if metricListRaw[i] == nil {
				continue
			}
			fmt.Fprint(b, *metricListRaw[i])
		}
	}

	writeString(checkCloudWatchNamespaceAttr)
	writeString(checkCloudWatchURLAttr)
	writeString(checkCloudWatchVersionAttr)

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPICloudWatch(c *circonusCheck, l interfaceList) error { //nolint:unparam
	c.Type = string(apiCheckTypeCloudWatchAttr)

	// Iterate over all `cloudwatch` attributes, even though we have a max of 1 in the
	// schema.
	for _, mapRaw := range l {
		cloudwatchConfig := newInterfaceMap(mapRaw)

		if v, found := cloudwatchConfig[checkCloudWatchAPIKeyAttr]; found {
			c.Config[config.APIKey] = v.(string)
		}

		if v, found := cloudwatchConfig[checkCloudWatchAPISecretAttr]; found {
			c.Config[config.APISecret] = v.(string)
		}

		for k, v := range cloudwatchConfig.CollectMap(checkCloudWatchDimmensionsAttr) {
			dimKey := config.DimPrefix + config.Key(k)
			c.Config[dimKey] = v
		}

		if v, found := cloudwatchConfig[checkCloudWatchMetricAttr]; found {
			metricsRaw := v.(*schema.Set).List()
			metrics := make([]string, 0, len(metricsRaw))
			for _, m := range metricsRaw {
				metrics = append(metrics, m.(string))
			}
			sort.Strings(metrics)
			c.Config[config.CloudwatchMetrics] = strings.Join(metrics, ",")
		}

		if v, found := cloudwatchConfig[checkCloudWatchNamespaceAttr]; found {
			c.Config[config.Namespace] = v.(string)
		}

		if v, found := cloudwatchConfig[checkCloudWatchURLAttr]; found {
			c.Config[config.URL] = v.(string)
		}

		if v, found := cloudwatchConfig[checkCloudWatchVersionAttr]; found {
			c.Config[config.Version] = v.(string)
		}
	}

	return nil
}
