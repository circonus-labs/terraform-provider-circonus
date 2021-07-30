// Package circonus defines the Terraform Circonus provider
package circonus

import (
	"context"
	"fmt"
	"log"
	"os"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultCirconus404ErrorString        = "API response code 404:"
	defaultCirconusAggregationWindow     = "300s"
	defaultCirconusAlertMinEscalateAfter = "300s"
	defaultCirconusCheckPeriodMax        = "300s"
	defaultCirconusCheckPeriodMin        = "10s"
	defaultCirconusHTTPFormat            = "json"
	defaultCirconusHTTPMethod            = "POST"
	defaultCirconusSlackUsername         = "Circonus"
	defaultCirconusTimeoutMax            = "300s"
	defaultCirconusTimeoutMin            = "0s"
	maxSeverity                          = 5
	minSeverity                          = 0
)

var providerDescription = map[string]string{
	providerAPIURLAttr:  "URL of the Circonus API",
	providerAutoTagAttr: "Signals that the provider should automatically add a tag to all API calls denoting that the resource was created by Terraform",
	providerKeyAttr:     "API token used to authenticate with the Circonus API",
}

// Constants that want to be a constant but can't in Go.
var (
	validContactHTTPFormats = validStringValues{"json", "params"}
	validContactHTTPMethods = validStringValues{"GET", "POST"}
)

type contactMethods string

// globalAutoTag controls whether or not the provider should automatically add a
// tag to each resource.
//
// NOTE(sean): This is done as a global variable because the diff suppress
// functions does not have access to the providerContext, only the key, old, and
// new values.
var globalAutoTag bool //nolint:unused

type providerContext struct {
	// Circonus API client
	client *api.API
	// defaultTag make up the tag to be used when autoTag tags a tag.
	defaultTag circonusTag
	// autoTag, when true, automatically appends defaultCirconusTag
	autoTag bool
}

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			providerAPIURLAttr: {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("CIRCONUS_API_URL", "https://api.circonus.com/v2"),
				Description: providerDescription[providerAPIURLAttr],
			},
			providerAutoTagAttr: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     defaultAutoTag,
				Description: providerDescription[providerAutoTagAttr],
			},
			providerKeyAttr: {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("CIRCONUS_API_TOKEN", nil),
				Description: providerDescription[providerKeyAttr],
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"circonus_account":   dataSourceCirconusAccount(),
			"circonus_collector": dataSourceCirconusCollector(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"circonus_check":          resourceCheck(),
			"circonus_contact_group":  resourceContactGroup(),
			"circonus_graph":          resourceGraph(),
			"circonus_overlay_set":    resourceOverlaySet(),
			"circonus_dashboard":      resourceDashboard(),
			"circonus_maintenance":    resourceMaintenance(),
			"circonus_metric":         resourceMetric(),
			"circonus_rule_set":       resourceRuleSet(),
			"circonus_rule_set_group": resourceRuleSetGroup(),
			"circonus_worksheet":      resourceWorksheet(),
		},

		ConfigureContextFunc: providerConfigure,
	}

	return p
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	globalAutoTag = d.Get(providerAutoTagAttr).(bool)

	envLevel := os.Getenv("TF_LOG")
	debug := false
	if envLevel != "" {
		debug = true
	}

	config := &api.Config{
		URL:      d.Get(providerAPIURLAttr).(string),
		TokenKey: d.Get(providerKeyAttr).(string),
		TokenApp: "terraform-provider-circonus",
	}

	if debug {
		config.Debug = true
		config.Log = log.New(log.Writer(), "", log.LstdFlags)
	}

	// turn on logging if terraform log level set to debug
	if os.Getenv("TF_LOG") == "DEBUG" {
		config.Debug = true
		config.Log = log.New(log.Writer(), "", log.LstdFlags)
	}

	var diags diag.Diagnostics

	client, err := api.NewAPI(config)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error initializing Circonus",
			Detail:   fmt.Sprintf("Unable to initialize Circonus API client: %s", err),
		})
		return nil, diags
	}

	client.EnableExponentialBackoff()

	return &providerContext{
		client:     client,
		autoTag:    d.Get(providerAutoTagAttr).(bool),
		defaultTag: defaultCirconusTag,
	}, diags
}
