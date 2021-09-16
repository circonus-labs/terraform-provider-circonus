package circonus

import (
	"context"

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
	collectorDetailNameAttr   = "name"
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
			// _cid
			collectorIDAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateRegexp(collectorIDAttr, config.BrokerCIDRegex),
				Description:  collectorDescription[collectorIDAttr],
			},
			// _details
			collectorDetailsAttr: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: collectorDescription[collectorDetailsAttr],
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// cn
						collectorCNAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorCNAttr],
						},
						// name
						collectorDetailNameAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorDetailNameAttr],
						},
						// external_host
						collectorExternalHostAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorExternalHostAttr],
						},
						// external_port
						collectorExternalPortAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: collectorDescription[collectorExternalPortAttr],
						},
						// ipaddress
						collectorIPAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorIPAttr],
						},
						// minimum_version_required
						collectorMinVersionAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: collectorDescription[collectorMinVersionAttr],
						},
						// modules
						collectorModulesAttr: {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: collectorDescription[collectorModulesAttr],
						},
						// port
						collectorPortAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: collectorDescription[collectorPortAttr],
						},
						// skew
						collectorSkewAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorSkewAttr],
						},
						// status
						collectorStatusAttr: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: collectorDescription[collectorStatusAttr],
						},
						// version
						collectorVersionAttr: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: collectorDescription[collectorVersionAttr],
						},
					},
				},
			},
			// _latitude
			collectorLatitudeAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: collectorDescription[collectorLatitudeAttr],
			},
			// _longitude
			collectorLongitudeAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: collectorDescription[collectorLongitudeAttr],
			},
			// name
			collectorNameAttr: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: collectorDescription[collectorNameAttr],
			},
			// _tags
			collectorTagsAttr: tagMakeConfigSchema(collectorTagsAttr),
			// _type
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

	cid := d.Id()
	if cidRaw, ok := d.GetOk(collectorIDAttr); ok {
		cid = cidRaw.(string)
	}
	broker, err := client.FetchBroker(api.CIDType(&cid))
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(broker.CID)
	if err := d.Set(collectorIDAttr, broker.CID); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(collectorDetailsAttr, collectorDetailsToState(broker)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(collectorLatitudeAttr, broker.Latitude); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(collectorLongitudeAttr, broker.Longitude); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(collectorNameAttr, broker.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(collectorTagsAttr, broker.Tags); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set(collectorTypeAttr, broker.Type); err != nil {
		return diag.FromErr(err)
	}

	return diags
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
