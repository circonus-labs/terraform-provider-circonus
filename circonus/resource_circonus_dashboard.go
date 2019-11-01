package circonus

import (
	"fmt"
	"strings"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDashboard() *schema.Resource {
	return &schema.Resource{
		Create: dashboardCreate,
		Read:   dashboardRead,
		Update: dashboardUpdate,
		Delete: dashboardDelete,
		Exists: dashboardExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"title": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				StateFunc:   suppressWhitespace,
				Description: "The title of the dashboard.",
			},
			"uuid": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uuid of the dashboard.",
			},
			"shared": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"account_default": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"grid_layout": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
				Elem:     schema.TypeInt,
			},
			"options": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_screen_hide_title": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"hide_grid": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						// "linkages": &schema.Schema{
						// 	Type:     schema.TypeList,
						// 	Optional: true,
						// 	Elem: &schema.Resource{
						// 		Schema: &schema.Schema{
						// 			Type:     schema.TypeList,
						// 			Optional: true,
						// 			Elem:     schema.TypeString,
						// 		},
						// 	},
						// },
						"scale_text": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"text_size": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  12,
						},
						"access_configs": &schema.Schema{
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"black_dash": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"enabled": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"full_screen": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"full_screen_hide_title": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"nick_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"scale_text": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"shared_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"text_size": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
										Default:  16,
									},
								},
							},
						},
					},
				},
			},
			"widget": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"height": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"origin": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"settings": &schema.Schema{
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"acknowledged": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"algorithm": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"body_format": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"chart_type": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"check_uuid": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"cleared": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"cluster_id": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"cluster_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"content_type": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"datapoints": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"_metric_type": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"_check_id": &schema.Schema{
													Type:     schema.TypeInt,
													Required: true,
												},
												"label": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
												"metric": &schema.Schema{
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"dependents": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"disable_autoformat": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"display": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"format": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"formula": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"layout": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"limit": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"maintenance": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"markup": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"metric_display_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"metric_name": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"min_age": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"overlay_set_id": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"range_high": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"range_low": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"resource_limit": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"resource_usage": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"search": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"severity": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"size": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold": &schema.Schema{
										Type:     schema.TypeFloat,
										Optional: true,
									},
									"thresholds": &schema.Schema{
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"colors": &schema.Schema{
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"values": &schema.Schema{
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"flip": &schema.Schema{
													Type:     schema.TypeBool,
													Required: true,
												},
											},
										},
									},
									"time_window": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"title": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"title_format": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"trend": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"use_default": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"value_type": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"date_window": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"graph_uuid": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"hide_xaxis": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"hide_yaxis": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"key_inline": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"key_loc": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"key_size": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"key_wrap": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"label": &schema.Schema{
										Type:     schema.TypeString,
										Optional: true,
									},
									"period": &schema.Schema{
										Type:     schema.TypeInt,
										Optional: true,
									},
									"real_time": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
									"show_flags": &schema.Schema{
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"widget_id": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"width": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dashboardCreate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	dash := newDashboard()
	if err := dash.ParseConfig(d); err != nil {
		return errwrap.Wrapf("error parsing graph schema during create: {{err}}", err)
	}

	if err := dash.Create(ctxt); err != nil {
		return errwrap.Wrapf("error creating dashboard: {{err}}", err)
	}

	d.SetId(dash.CID)

	return dashboardRead(d, meta)
}

func dashboardExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	dash, err := ctxt.client.FetchDashboard(api.CIDType(&cid))
	if err != nil {
		if strings.Contains(err.Error(), defaultCirconus404ErrorString) {
			return false, nil
		}

		return false, err
	}

	if dash.CID == "" {
		return false, nil
	}

	return true, nil
}

// dashboardRead pulls data out of the Dashboard object and stores it into the
// appropriate place in the statefile.
func dashboardRead(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	dash, err := loadDashboard(ctxt, api.CIDType(&cid))
	if err != nil {
		return err
	}

	d.SetId(dash.CID)

	widgets := make([]map[string]interface{}, len(dash.Widgets))
	for i, widget := range dash.Widgets {
		dashWidgetAttrs := make(map[string]interface{}, 8) // 8 == len(members in api.DashboardWidget)

		dashWidgetAttrs["active"] = widget.Active
		dashWidgetAttrs["height"] = int(widget.Height)
		dashWidgetAttrs["name"] = widget.Name
		dashWidgetAttrs["origin"] = widget.Origin
		dashWidgetAttrs["type"] = widget.Type
		dashWidgetAttrs["widget_id"] = widget.WidgetID
		dashWidgetAttrs["width"] = int(widget.Width)

		dashWidgetSettingsAttrs := make(map[string]interface{}, 62)
		dashWidgetSettingsAttrs["account_id"] = widget.Settings.AccountID
		dashWidgetSettingsAttrs["acknowledged"] = widget.Settings.Acknowledged
		dashWidgetSettingsAttrs["algorithm"] = widget.Settings.Algorithm
		dashWidgetSettingsAttrs["body_format"] = widget.Settings.BodyFormat
		dashWidgetSettingsAttrs["chart_type"] = widget.Settings.ChartType
		dashWidgetSettingsAttrs["check_uuid"] = widget.Settings.CheckUUID
		dashWidgetSettingsAttrs["cleared"] = widget.Settings.Cleared
		dashWidgetSettingsAttrs["cluster_id"] = int(widget.Settings.ClusterID)
		dashWidgetSettingsAttrs["cluster_name"] = widget.Settings.ClusterName
		dashWidgetSettingsAttrs["contact_groups"] = widget.Settings.ContactGroups
		dashWidgetSettingsAttrs["content_type"] = widget.Settings.ContentType
		dps := make([]map[string]interface{}, 0, len(widget.Settings.Datapoints))
		for _, dp := range widget.Settings.Datapoints {
			dpAttrs := make(map[string]interface{}, 4)
			dpAttrs["label"] = dp.Label
			dpAttrs["_metric_type"] = dp.MetricType
			dpAttrs["metric"] = dp.Metric
			dpAttrs["_check_id"] = int(dp.CheckID)
			dps = append(dps, dpAttrs)
		}
		dashWidgetSettingsAttrs["datapoints"] = dps
		dashWidgetSettingsAttrs["dependents"] = widget.Settings.Dependents
		dashWidgetSettingsAttrs["disable_autoformat"] = widget.Settings.DisableAutoformat
		dashWidgetSettingsAttrs["display"] = widget.Settings.Display
		dashWidgetSettingsAttrs["format"] = widget.Settings.Format
		dashWidgetSettingsAttrs["formula"] = widget.Settings.Formula
		dashWidgetSettingsAttrs["layout"] = widget.Settings.Layout
		dashWidgetSettingsAttrs["limit"] = int(widget.Settings.Limit)
		dashWidgetSettingsAttrs["maintenance"] = widget.Settings.Maintenance
		dashWidgetSettingsAttrs["markup"] = widget.Settings.Markup
		dashWidgetSettingsAttrs["metric_display_name"] = widget.Settings.MetricDisplayName
		dashWidgetSettingsAttrs["metric_name"] = widget.Settings.MetricName
		dashWidgetSettingsAttrs["min_age"] = widget.Settings.MinAge
		dashWidgetSettingsAttrs["off_hours"] = widget.Settings.OffHours
		dashWidgetSettingsAttrs["overlay_set_id"] = widget.Settings.OverlaySetID
		dashWidgetSettingsAttrs["range_high"] = int(widget.Settings.RangeHigh)
		dashWidgetSettingsAttrs["range_low"] = int(widget.Settings.RangeLow)
		dashWidgetSettingsAttrs["resource_limit"] = widget.Settings.ResourceLimit
		dashWidgetSettingsAttrs["resource_usage"] = widget.Settings.ResourceUsage
		dashWidgetSettingsAttrs["search"] = widget.Settings.Search
		dashWidgetSettingsAttrs["severity"] = widget.Settings.Severity
		dashWidgetSettingsAttrs["size"] = widget.Settings.Size
		dashWidgetSettingsAttrs["tag_filter_set"] = widget.Settings.TagFilterSet
		dashWidgetSettingsAttrs["threshold"] = widget.Settings.Threshold
		if widget.Settings.Thresholds != nil {
			t := make(map[string]interface{}, 3)
			t["colors"] = widget.Settings.Thresholds.Colors
			t["values"] = widget.Settings.Thresholds.Values
			t["flip"] = widget.Settings.Thresholds.Flip
			dashWidgetSettingsAttrs["thresholds"] = t
		}
		dashWidgetSettingsAttrs["time_window"] = widget.Settings.TimeWindow
		dashWidgetSettingsAttrs["title"] = widget.Settings.Title
		dashWidgetSettingsAttrs["title_format"] = widget.Settings.TitleFormat
		dashWidgetSettingsAttrs["trend"] = widget.Settings.Trend
		dashWidgetSettingsAttrs["type"] = widget.Settings.Type
		dashWidgetSettingsAttrs["use_default"] = widget.Settings.UseDefault
		dashWidgetSettingsAttrs["value_type"] = widget.Settings.ValueType
		dashWidgetSettingsAttrs["week_days"] = widget.Settings.WeekDays
		dashWidgetSettingsAttrs["date_window"] = widget.Settings.DateWindow
		dashWidgetSettingsAttrs["graph_uuid"] = widget.Settings.GraphUUID
		dashWidgetSettingsAttrs["hide_xaxis"] = widget.Settings.HideXAxis
		dashWidgetSettingsAttrs["hide_yaxis"] = widget.Settings.HideYAxis
		dashWidgetSettingsAttrs["key_inline"] = widget.Settings.KeyInline
		dashWidgetSettingsAttrs["key_loc"] = widget.Settings.KeyLoc
		dashWidgetSettingsAttrs["key_size"] = int(widget.Settings.KeySize)
		dashWidgetSettingsAttrs["key_wrap"] = widget.Settings.KeyWrap
		dashWidgetSettingsAttrs["label"] = widget.Settings.Label
		dashWidgetSettingsAttrs["period"] = int(widget.Settings.Period)
		dashWidgetSettingsAttrs["realtime"] = widget.Settings.Realtime
		dashWidgetSettingsAttrs["show_flags"] = widget.Settings.ShowFlags

		settings := make([]map[string]interface{}, 1)
		settings[0] = dashWidgetSettingsAttrs

		dashWidgetAttrs["settings"] = settings

		widgets[i] = dashWidgetAttrs
	}
	d.Set("widget", widgets)

	options := make([]map[string]interface{}, 1)
	optionsAttrs := make(map[string]interface{}, 6)
	optionsAttrs["full_screen_hide_title"] = dash.Options.FullscreenHideTitle
	optionsAttrs["hide_grid"] = dash.Options.HideGrid
	optionsAttrs["scale_text"] = dash.Options.ScaleText
	optionsAttrs["text_size"] = dash.Options.TextSize

	accessConfigs := make([]map[string]interface{}, 0, len(dash.Options.AccessConfigs))
	for _, ac := range dash.Options.AccessConfigs {
		acAttrs := make(map[string]interface{}, 8)
		acAttrs["black_dash"] = ac.BlackDash
		acAttrs["enabled"] = ac.Enabled
		acAttrs["full_screen"] = ac.Fullscreen
		acAttrs["full_screen_hide_title"] = ac.FullscreenHideTitle
		acAttrs["nick_name"] = ac.Nickname
		acAttrs["scale_text"] = ac.ScaleText
		acAttrs["shared_id"] = ac.SharedID
		acAttrs["text_size"] = ac.TextSize
		accessConfigs = append(accessConfigs, acAttrs)
	}
	optionsAttrs["access_configs"] = accessConfigs
	options[0] = optionsAttrs
	d.Set("options", options)

	gridLayoutAttrs := make(map[string]interface{}, 2)
	gridLayoutAttrs["width"] = dash.GridLayout.Width
	gridLayoutAttrs["height"] = dash.GridLayout.Height
	d.Set("grid_layout", gridLayoutAttrs)
	d.Set("account_default", dash.AccountDefault)
	d.Set("shared", dash.Shared)
	d.Set("title", dash.Title)
	d.Set("uuid", dash.UUID)

	return nil
}

func dashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	dash := newDashboard()
	if err := dash.ParseConfig(d); err != nil {
		return err
	}

	dash.CID = d.Id()
	if err := dash.Update(ctxt); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to update dashboard %q: {{err}}", d.Id()), err)
	}

	return dashboardRead(d, meta)
}

func dashboardDelete(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	if _, err := ctxt.client.DeleteDashboardByCID(api.CIDType(&cid)); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to delete dashboard %q: {{err}}", d.Id()), err)
	}

	d.SetId("")
	d.Set("uuid", "")

	return nil
}

type circonusDashboard struct {
	api.Dashboard
}

func newDashboard() circonusDashboard {
	dash := circonusDashboard{
		Dashboard: *api.NewDashboard(),
	}

	return dash
}

func loadDashboard(ctxt *providerContext, cid api.CIDType) (circonusDashboard, error) {
	var dash circonusDashboard
	ng, err := ctxt.client.FetchDashboard(cid)
	if err != nil {
		return circonusDashboard{}, err
	}
	dash.Dashboard = *ng

	return dash, nil
}

// ParseConfig reads Terraform config data and stores the information into a
// Circonus Dashboard object.  ParseConfig and dashboardRead() must be kept in sync.
func (dash *circonusDashboard) ParseConfig(d *schema.ResourceData) error {
	dash.Widgets = make([]api.DashboardWidget, 0, defaultDashboardWidgets)

	if v, found := d.GetOk("title"); found {
		dash.Title = v.(string)
	}
	if v, found := d.GetOk("shared"); found {
		dash.Shared = v.(bool)
	}
	if v, found := d.GetOk("account_default"); found {
		dash.AccountDefault = v.(bool)
	}

	if v, found := d.GetOk("grid_layout"); found {
		listRaw := v.(map[string]interface{})
		glMap := make(map[string]interface{}, len(listRaw))
		for k, v := range listRaw {
			glMap[k] = v
		}
		if v, ok := glMap["width"]; ok {
			dash.GridLayout.Width = uint(v.(int))
		}
		if v, ok := glMap["height"]; ok {
			dash.GridLayout.Height = uint(v.(int))
		}

	}
	if v, found := d.GetOk("options"); found {
		optionList := v.(*schema.Set).List()
		for _, optionElem := range optionList {
			optionsMap := newInterfaceMap(optionElem.(map[string]interface{}))

			if v, ok := optionsMap["full_screen_hide_title"]; ok {
				dash.Options.FullscreenHideTitle = v.(bool)
			}
			if v, ok := optionsMap["hide_grid"]; ok {
				dash.Options.HideGrid = v.(bool)
			}
			if v, ok := optionsMap["scale_text"]; ok {
				dash.Options.ScaleText = v.(bool)
			}
			if v, ok := optionsMap["text_size"]; ok {
				dash.Options.TextSize = uint(v.(int))
			}

			dash.Options.AccessConfigs = make([]api.DashboardAccessConfig, 0)
			if listRaw, found := optionsMap["access_configs"]; found {
				accessConfigsList := listRaw.(*schema.Set).List()
				for _, aclElem := range accessConfigsList {
					acAttrs := aclElem.(map[string]interface{})
					acl := api.DashboardAccessConfig{}

					if v, found := acAttrs["black_dash"]; found {
						acl.BlackDash = v.(bool)
					}

					if v, found := acAttrs["enabled"]; found {
						acl.Enabled = v.(bool)
					}

					if v, found := acAttrs["full_screen"]; found {
						acl.Fullscreen = v.(bool)
					}

					if v, found := acAttrs["full_screen_hide_title"]; found {
						acl.FullscreenHideTitle = v.(bool)
					}

					if v, found := acAttrs["nick_name"]; found {
						acl.Nickname = (v.(string))
					}

					if v, found := acAttrs["scale_text"]; found {
						acl.ScaleText = (v.(bool))
					}

					if v, found := acAttrs["shared_id"]; found {
						acl.SharedID = (v.(string))
					}

					if v, found := acAttrs["text_size"]; found {
						acl.TextSize = v.(uint)
					}
					dash.Options.AccessConfigs = append(dash.Options.AccessConfigs, acl)
				}
			}
			dash.Options.Linkages = make([][]string, 0)
		}

	}

	if listRaw, found := d.GetOk("widget"); found {
		widgetList := listRaw.(*schema.Set).List()
		for _, widgetListElem := range widgetList {
			wAttrs := widgetListElem.(map[string]interface{})

			w := api.DashboardWidget{}

			if v, found := wAttrs["active"]; found {
				w.Active = v.(bool)
			}

			if v, found := wAttrs["height"]; found {
				w.Height = uint(v.(int))
			}

			if v, found := wAttrs["name"]; found {
				w.Name = (v.(string))
			}

			if v, found := wAttrs["origin"]; found {
				w.Origin = (v.(string))
			}

			if v, found := wAttrs["type"]; found {
				w.Type = (v.(string))
			}

			if v, found := wAttrs["widget_id"]; found {
				w.WidgetID = (v.(string))
			}

			if v, found := wAttrs["width"]; found {
				w.Width = uint(v.(int))
			}

			if mapRaw, found := wAttrs["settings"]; found {
				listRaw := mapRaw.(*schema.Set).List()
				for _, settingElem := range listRaw {
					sMap := settingElem.(map[string]interface{})

					if v, found := sMap["account_id"]; found {
						w.Settings.AccountID = (v.(string))
					}
					if v, found := sMap["acknowledged"]; found {
						w.Settings.Acknowledged = (v.(string))
					}
					if v, found := sMap["algorithm"]; found {
						w.Settings.Algorithm = (v.(string))
					}
					if v, found := sMap["body_format"]; found {
						w.Settings.BodyFormat = (v.(string))
					}
					if v, found := sMap["chart_type"]; found {
						w.Settings.ChartType = (v.(string))
					}
					if v, found := sMap["check_uuid"]; found {
						w.Settings.CheckUUID = (v.(string))
					}
					if v, found := sMap["cleared"]; found {
						w.Settings.Cleared = (v.(string))
					}
					if v, found := sMap["cluster_id"]; found {
						w.Settings.ClusterID = uint(v.(int))
					}
					if v, found := sMap["cluster_name"]; found {
						w.Settings.ClusterName = (v.(string))
					}
					// if v, found := sMap[string(dashWidgetSettingsContactGroupsAttr)]; found {
					// 	w.Settings.ContactGroups = (v.(string))
					// }
					if v, found := sMap["content_type"]; found {
						w.Settings.ContentType = (v.(string))
					}
					if v, found := sMap["date_window"]; found {
						w.Settings.DateWindow = (v.(string))
					}
					// if v, found := sMap[string(dashWidgetSettingsDefinitionAttr)]; found {
					// 	w.Settings.Definition = (v.(string))
					// }
					if v, found := sMap["datapoints"]; found {
						w.Settings.Datapoints = make([]api.ChartTextWidgetDatapoint, 0)

						datapointList := v.(*schema.Set).List()
						for _, dpElem := range datapointList {
							dpAttrs := dpElem.(map[string]interface{})
							dp := api.ChartTextWidgetDatapoint{}
							if vv, found := dpAttrs["label"]; found {
								dp.Label = (vv.(string))
							}
							if vv, found := dpAttrs["_metric_type"]; found {
								dp.MetricType = (vv.(string))
							}
							if vv, found := dpAttrs["_check_id"]; found {
								dp.CheckID = uint(vv.(int))
							}
							if vv, found := dpAttrs["metric"]; found {
								dp.Metric = (vv.(string))
							}
							w.Settings.Datapoints = append(w.Settings.Datapoints, dp)
						}
					}
					if v, found := sMap["dependents"]; found {
						w.Settings.Dependents = (v.(string))
					}

					if v, found := sMap["disable_autoformat"]; found {
						w.Settings.DisableAutoformat = v.(bool)
					}
					if v, found := sMap["display"]; found {
						w.Settings.Display = (v.(string))
					}
					if v, found := sMap["format"]; found {
						w.Settings.Format = (v.(string))
					}
					if v, found := sMap["formula"]; found {
						w.Settings.Formula = (v.(string))
					}
					if v, found := sMap["graph_uuid"]; found {
						w.Settings.GraphUUID = (v.(string))
					}
					if v, found := sMap["hide_xaxis"]; found {
						w.Settings.HideXAxis = v.(bool)
					}
					if v, found := sMap["hide_yaxis"]; found {
						w.Settings.HideYAxis = v.(bool)
					}
					if v, found := sMap["key_inline"]; found {
						w.Settings.KeyInline = v.(bool)
					}
					if v, found := sMap["key_loc"]; found {
						w.Settings.KeyLoc = (v.(string))
					}
					if v, found := sMap["key_size"]; found {
						w.Settings.KeySize = uint(v.(int))
					}
					if v, found := sMap["key_wrap"]; found {
						w.Settings.KeyWrap = v.(bool)
					}
					if v, found := sMap["label"]; found {
						w.Settings.Label = (v.(string))
					}
					if v, found := sMap["layout"]; found {
						w.Settings.Layout = (v.(string))
					}
					if v, found := sMap["limit"]; found {
						w.Settings.Limit = uint(v.(int))
					}
					if v, found := sMap["maintenance"]; found {
						w.Settings.Maintenance = (v.(string))
					}
					if v, found := sMap["markup"]; found {
						w.Settings.Markup = (v.(string))
					}
					if v, found := sMap["metric_display_name"]; found {
						w.Settings.MetricDisplayName = (v.(string))
					}
					if v, found := sMap["metric_name"]; found {
						w.Settings.MetricName = (v.(string))
					}
					if v, found := sMap["min_age"]; found {
						w.Settings.MinAge = (v.(string))
					}
					if v, found := sMap["overlay_set_id"]; found {
						w.Settings.OverlaySetID = (v.(string))
					}
					if v, found := sMap["period"]; found {
						w.Settings.Period = uint(v.(int))
					}
					if v, found := sMap["range_high"]; found {
						w.Settings.RangeHigh = v.(int)
					}
					if v, found := sMap["range_low"]; found {
						w.Settings.RangeLow = v.(int)
					}
					if v, found := sMap["realtime"]; found {
						w.Settings.Realtime = v.(bool)
					}
					if v, found := sMap["resource_limit"]; found {
						w.Settings.ResourceLimit = (v.(string))
					}
					if v, found := sMap["resource_usage"]; found {
						w.Settings.ResourceUsage = (v.(string))
					}
					if v, found := sMap["search"]; found {
						w.Settings.Search = (v.(string))
					}
					if v, found := sMap["severity"]; found {
						w.Settings.Severity = (v.(string))
					}
					if v, found := sMap["show_flags"]; found {
						w.Settings.ShowFlags = v.(bool)
					}
					if v, found := sMap["size"]; found {
						w.Settings.Size = (v.(string))
					}
					if v, found := sMap["threshold"]; found {
						w.Settings.Threshold = float32((v.(float64)))
					}
					if v, found := sMap["thresholds"]; found {
						t := api.ForecastGaugeWidgetThresholds{}

						// there will be only 1
						tList := v.(*schema.Set).List()
						for _, tElem := range tList {
							tAttrs := tElem.(map[string]interface{})
							if vv, found := tAttrs["colors"]; found {
								t.Colors = make([]string, len(vv.([]interface{})))
								for i, x := range vv.([]interface{}) {
									t.Colors[i] = (x.(string))
								}
							}
							if vv, found := tAttrs["values"]; found {
								t.Values = make([]string, len(vv.([]interface{})))
								for i, x := range vv.([]interface{}) {
									t.Values[i] = (x.(string))
								}
							}
							if vv, found := tAttrs["flip"]; found {
								t.Flip = vv.(bool)
							}
						}
						w.Settings.Thresholds = &t
					}
					if v, found := sMap["time_window"]; found {
						w.Settings.TimeWindow = (v.(string))
					}
					if v, found := sMap["title"]; found {
						w.Settings.Title = (v.(string))
					}
					if v, found := sMap["title_format"]; found {
						w.Settings.TitleFormat = (v.(string))
					}
					if v, found := sMap["trend"]; found {
						w.Settings.Trend = (v.(string))
					}
					if v, found := sMap["type"]; found {
						w.Settings.Type = (v.(string))
					}
					if v, found := sMap["use_default"]; found {
						w.Settings.UseDefault = v.(bool)
					}
					if v, found := sMap["value_type"]; found {
						w.Settings.ValueType = (v.(string))
					}
				}
				dash.Widgets = append(dash.Widgets, w)
			}
		}
	}

	if err := dash.Validate(); err != nil {
		return err
	}

	return nil
}

func (dash *circonusDashboard) Create(ctxt *providerContext) error {
	ctxt.client.Debug = true
	ng, err := ctxt.client.CreateDashboard(&dash.Dashboard)
	if err != nil {
		return err
	}

	dash.CID = ng.CID
	dash.UUID = ng.UUID

	return nil
}

func (dash *circonusDashboard) Update(ctxt *providerContext) error {
	_, err := ctxt.client.UpdateDashboard(&dash.Dashboard)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to update dashboard %s: {{err}}", dash.CID), err)
	}

	return nil
}

func (g *circonusDashboard) Validate() error {
	// for i, datapoint := range g.Datapoints {
	// 	if *g.Style == apiGraphStyleLine && datapoint.Alpha != nil && *datapoint.Alpha != 0 {
	// 		return fmt.Errorf("%s can not be set on graphs with style %s", graphMetricAlphaAttr, apiGraphStyleLine)
	// 	}

	// 	if datapoint.CheckID != 0 && datapoint.MetricName == "" {
	// 		return fmt.Errorf("Error with %s[%d] name=%q: %s is set, missing attribute %s must also be set", graphMetricAttr, i, datapoint.Name, graphMetricCheckAttr, graphMetricNameAttr)
	// 	}

	// 	if datapoint.CheckID == 0 && datapoint.MetricName != "" {
	// 		return fmt.Errorf("Error with %s[%d] name=%q: %s is set, missing attribute %s must also be set", graphMetricAttr, i, datapoint.Name, graphMetricNameAttr, graphMetricCheckAttr)
	// 	}

	// 	if datapoint.CAQL != nil && (datapoint.CheckID != 0 || datapoint.MetricName != "") {
	// 		return fmt.Errorf("Error with %s[%d] name=%q: %q attribute is mutually exclusive with attributes %s or %s or %s", graphMetricAttr, i, datapoint.Name, graphMetricCAQLAttr, graphMetricNameAttr, graphMetricCheckAttr, graphMetricSearchAttr)
	// 	}

	// 	if datapoint.Search != nil && (datapoint.CheckID != 0 || datapoint.MetricName != "") {
	// 		return fmt.Errorf("Error with %s[%d] name=%q: %q attribute is mutually exclusive with attributes %s or %s or %s", graphMetricAttr, i, datapoint.Name, graphMetricSearchAttr, graphMetricNameAttr, graphMetricCheckAttr, graphMetricCAQLAttr)
	// 	}

	// 	if datapoint.MetricType == "text" && datapoint.Derive != nil {
	// 		v := datapoint.Derive
	// 		switch v.(type) {
	// 		case bool:
	// 		default:
	// 			return fmt.Errorf("Error with %s[%d] (name=%q): attribute %q is mutually exclusive when %s=%q", graphMetricAttr, i, datapoint.Name, graphMetricFunctionAttr, graphMetricMetricTypeAttr, "text")
	// 		}
	// 	}
	// }

	// for i, mc := range g.MetricClusters {
	// 	if mc.AggregateFunc != "" && (mc.Color == nil || *mc.Color == "") {
	// 		return fmt.Errorf("Error with %s[%d] name=%q: %s is a required attribute for graphs with %s set", graphMetricClusterAttr, i, mc.Name, graphMetricClusterColorAttr, graphMetricClusterAggregateAttr)
	// 	}
	// }

	return nil
}
