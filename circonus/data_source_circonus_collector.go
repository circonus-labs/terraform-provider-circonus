package circonus

import (
	"context"
	"fmt"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	collectorCNAttr           = "cn"
	collectorIDAttr           = "id"
	collectorDetailsAttr      = "details"
	collectorExternalHostAttr = "external_host"
	collectorExternalPortAttr = "external_port"
	collectorIPAttr           = "ip"
	collectorLatitudeAttr     = "latitude"
	collectorLongitudeAttr    = "longitude"
	collectorMinVersionAttr   = "min_version"
	collectorModulesAttr      = "modules"
	collectorNameAttr         = "name"
	collectorPortAttr         = "port"
	collectorSkewAttr         = "skew"
	collectorStatusAttr       = "status"
	collectorTagsAttr         = "tags"
	collectorTypeAttr         = "type"
	collectorVersionAttr      = "version"
)

var collectorDescription = map[schemaAttr]string{
	collectorDetailsAttr: "Details associated with individual collectors (a.k.a. broker)",
	collectorTagsAttr:    "Tags assigned to a collector",
}

func dataSourceCirconusCollector() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCirconusCollectorRead,

		Schema: map[string]*schema.Schema{
			collectorDetailsAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: collectorDescription[collectorDetailsAttr],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						collectorCNAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorCNAttr],
						},
						collectorExternalHostAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorExternalHostAttr],
						},
						collectorExternalPortAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: collectorDescription[collectorExternalPortAttr],
						},
						collectorIPAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorIPAttr],
						},
						collectorMinVersionAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: collectorDescription[collectorMinVersionAttr],
						},
						collectorModulesAttr: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: collectorDescription[collectorModulesAttr],
						},
						collectorPortAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: collectorDescription[collectorPortAttr],
						},
						collectorSkewAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorSkewAttr],
						},
						collectorStatusAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorStatusAttr],
						},
						collectorVersionAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: collectorDescription[collectorVersionAttr],
						},
					},
				},
			},
			collectorIDAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateRegexp(collectorIDAttr, config.BrokerCIDRegex),
				Description:  collectorDescription[collectorIDAttr],
			},
			collectorLatitudeAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: collectorDescription[collectorLatitudeAttr],
			},
			collectorLongitudeAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: collectorDescription[collectorLongitudeAttr],
			},
			collectorNameAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: collectorDescription[collectorNameAttr],
			},
			collectorTagsAttr: tagMakeConfigSchema(collectorTagsAttr),
			collectorTypeAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: collectorDescription[collectorTypeAttr],
			},
		},
	}
}

func dataSourceCirconusCollectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*providerContext).client
	var diags diag.Diagnostics

	var collector *api.Broker
	var err error
	cid := d.Id()
	if cidRaw, ok := d.GetOk(collectorIDAttr); ok {
		cid = cidRaw.(string)
	}
	collector, err = client.FetchBroker(api.CIDType(&cid))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Error fetching brokers",
			Detail:   fmt.Sprintf("Unable to fetch brokers: %s", err),
		})
		return diags
	}

	d.SetId(collector.CID)

	if err := d.Set(collectorDetailsAttr, collectorDetailsToState(collector)); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to store broker details",
			Detail:   fmt.Sprintf("details (%q) attribute: %s", collectorDetailsAttr, err),
		})
		return diags
	}

	_ = d.Set(collectorIDAttr, collector.CID)
	_ = d.Set(collectorLatitudeAttr, collector.Latitude)
	_ = d.Set(collectorLongitudeAttr, collector.Longitude)
	_ = d.Set(collectorNameAttr, collector.Name)
	_ = d.Set(collectorTagsAttr, collector.Tags)
	_ = d.Set(collectorTypeAttr, collector.Type)

	return nil
}

func collectorDetailsToState(c *api.Broker) []interface{} {
	details := make([]interface{}, 0, len(c.Details))

	for _, collector := range c.Details {
		collectorDetails := make(map[string]interface{}, defaultCollectorDetailAttrs)

		collectorDetails[collectorCNAttr] = collector.CN

		if collector.ExternalHost != nil {
			collectorDetails[collectorExternalHostAttr] = *collector.ExternalHost
		}

		if collector.ExternalPort != 0 {
			collectorDetails[collectorExternalPortAttr] = collector.ExternalPort
		}

		if collector.IP != nil {
			collectorDetails[collectorIPAttr] = *collector.IP
		}

		if collector.MinVer != 0 {
			collectorDetails[collectorMinVersionAttr] = collector.MinVer
		}

		if len(collector.Modules) > 0 {
			collectorDetails[collectorModulesAttr] = collector.Modules
		}

		if collector.Port != nil {
			collectorDetails[collectorPortAttr] = *collector.Port
		}

		if collector.Skew != nil {
			collectorDetails[collectorSkewAttr] = *collector.Skew
		}

		if collector.Status != "" {
			collectorDetails[collectorStatusAttr] = collector.Status
		}

		if collector.Version != nil {
			collectorDetails[collectorVersionAttr] = *collector.Version
		}

		details = append(details, collectorDetails)
	}

	return details
}
