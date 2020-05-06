package circonus

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				StateFunc:   suppressWhitespace,
				Description: "The title of the dashboard.",
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uuid of the dashboard.",
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"account_default": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"grid_layout": {
				Type:     schema.TypeMap,
				Required: true,
				Elem:     schema.TypeInt,
			},
			"options": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"full_screen_hide_title": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"hide_grid": {
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
						"scale_text": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"text_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  12,
						},
						"access_configs": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"black_dash": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"enabled": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"full_screen": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"full_screen_hide_title": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"nick_name": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  "",
									},
									"scale_text": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  true,
									},
									"shared_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"text_size": {
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
			"widget": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"active": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"height": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"origin": {
							Type:     schema.TypeString,
							Required: true,
						},
						"settings": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"acknowledged": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"algorithm": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"autoformat": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"bad_rules": {
										Type:     schema.TypeList,
										Optional: true,
										MaxItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"value": {
													Type:     schema.TypeString,
													Required: true,
												},
												"criterion": {
													Type:     schema.TypeString,
													Required: true,
												},
												"color": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"body_format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"caql": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"chart_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"check_uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cleared": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"cluster_id": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"cluster_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"content_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"datapoints": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"_metric_type": {
													Type:     schema.TypeString,
													Required: true,
												},
												"_check_id": {
													Type:     schema.TypeInt,
													Required: true,
												},
												"label": {
													Type:     schema.TypeString,
													Required: true,
												},
												"metric": {
													Type:     schema.TypeString,
													Required: true,
												},
											},
										},
									},
									"dependents": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"disable_autoformat": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"display": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"display_markup": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"formula": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"good_color": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"layout": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"layout_style": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"limit": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"link_url": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"maintenance": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"markup": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"metric_display_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"metric_name": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"metric_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"min_age": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"overlay_set_id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"range_high": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"range_low": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"resource_limit": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"resource_usage": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"search": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"severity": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"show_value": {
										Type:     schema.TypeBool,
										Optional: true,
									},
									"size": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"text_align": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"threshold": {
										Type:     schema.TypeFloat,
										Optional: true,
										Default:  0,
									},
									"thresholds": {
										Type:     schema.TypeSet,
										Optional: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"colors": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"values": {
													Type:     schema.TypeList,
													Required: true,
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"flip": {
													Type:     schema.TypeBool,
													Required: true,
												},
											},
										},
									},
									"time_window": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"title": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"title_format": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"trend": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"use_default": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"value_type": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"date_window": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"graph_uuid": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"hide_xaxis": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"hide_yaxis": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"key_inline": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"key_loc": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"key_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"key_wrap": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"label": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"period": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
									},
									"real_time": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
									"show_flags": {
										Type:     schema.TypeBool,
										Optional: true,
										Default:  false,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"widget_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"width": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

type ByWidgetId []map[string]interface{}

func (a ByWidgetId) Len() int      { return len(a) }
func (a ByWidgetId) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByWidgetId) Less(i, j int) bool {
	x := a[i]
	y := a[j]
	return x["widget_id"].(string) < y["widget_id"].(string)
}

func hashWidgets(vv interface{}) int {
	b := &bytes.Buffer{}
	b.Grow(defaultHashBufSize)

	writeBool := func(m map[string]interface{}, attrName string) {
		if v, ok := m[attrName]; ok {
			fmt.Fprintf(b, "%t", v.(bool))
		}
	}

	writeFloat := func(m map[string]interface{}, attrName string) {
		if v, ok := m[attrName]; ok {
			fmt.Fprintf(b, "%f", v.(float32))
		}
	}

	writeInt := func(m map[string]interface{}, attrName string) {
		if v, ok := m[attrName]; ok {
			fmt.Fprintf(b, "%x", v.(int))
		}
	}

	writeString := func(m map[string]interface{}, attrName string) {
		if v, ok := m[attrName]; ok && v.(string) != "" {
			fmt.Fprint(b, strings.TrimSpace(v.(string)))
		}
	}

	widgetsRaw := vv.([]map[string]interface{})
	sort.Sort(ByWidgetId(widgetsRaw))
	for _, m := range widgetsRaw {

		// Order writes to the buffer using lexically sorted list for easy visual
		// reconciliation with other lists.
		writeBool(m, "active")
		writeInt(m, "height")
		writeString(m, "name")
		writeString(m, "origin")
		writeString(m, "type")
		writeString(m, "widget_id")
		writeInt(m, "width")

		if settingsRaw, ok := m["settings"]; ok {
			settingsMap := settingsRaw.([]map[string]interface{})[0]
			writeString(settingsMap, "account_id")
			writeString(settingsMap, "acknowledged")
			writeString(settingsMap, "algorithm")
			writeBool(settingsMap, "autoformat")
			if badRulesRaw, ok := settingsMap["bad_rules"]; ok {
				badRulesList := badRulesRaw.([]map[string]interface{})
				for i := range badRulesList {
					br := badRulesList[i]
					writeString(br, "value")
					writeString(br, "criterion")
					writeString(br, "color")
				}
			}
			writeString(settingsMap, "body_format")
			writeString(settingsMap, "caql")
			writeString(settingsMap, "chart_type")
			writeString(settingsMap, "check_uuid")
			writeString(settingsMap, "cleared")
			writeInt(settingsMap, "cluster_id")
			writeString(settingsMap, "cluster_name")
			writeString(settingsMap, "content_type")

			if datapointsRaw, ok := settingsMap["datapoints"]; ok {
				datapointsListRaw := datapointsRaw.([]map[string]interface{})
				for i := range datapointsListRaw {
					if datapointsListRaw[i] == nil {
						continue
					}
					dp := datapointsListRaw[i]
					writeString(dp, "_metric_type")
					writeInt(dp, "_check_id")
					writeString(dp, "label")
					writeString(dp, "metric")
				}
			}

			writeString(settingsMap, "dependents")
			writeBool(settingsMap, "disable_autoformat")
			writeString(settingsMap, "display")
			writeString(settingsMap, "display_markup")
			writeString(settingsMap, "format")
			writeString(settingsMap, "formula")
			writeString(settingsMap, "good_color")
			writeString(settingsMap, "layout")
			writeString(settingsMap, "layout_style")
			writeInt(settingsMap, "limit")
			writeString(settingsMap, "link_url")
			writeString(settingsMap, "maintenance")
			writeString(settingsMap, "markup")
			writeString(settingsMap, "metric_display_name")
			writeString(settingsMap, "metric_name")
			writeString(settingsMap, "metric_type")
			writeString(settingsMap, "min_age")
			writeString(settingsMap, "overlay_set_id")
			writeInt(settingsMap, "range_high")
			writeInt(settingsMap, "range_low")
			writeString(settingsMap, "resource_limit")
			writeString(settingsMap, "resource_usage")
			writeString(settingsMap, "search")
			writeString(settingsMap, "severity")
			writeBool(settingsMap, "show_value")
			writeString(settingsMap, "size")
			writeString(settingsMap, "text_align")
			writeFloat(settingsMap, "threshold")

			if thresholdsRaw, ok := settingsMap["thresholds"]; ok {
				thresholdsListRaw := thresholdsRaw.([]map[string]interface{})
				for i := range thresholdsListRaw {
					if thresholdsListRaw[i] == nil {
						continue
					}
					t := thresholdsListRaw[i]
					colors := t["colors"].([]string)
					for c := range colors {
						if colors[c] != "" {
							fmt.Fprint(b, strings.TrimSpace(colors[c]))
						}
					}
					values := t["values"].([]string)
					for v := range values {
						if values[v] != "" {
							fmt.Fprint(b, strings.TrimSpace(values[v]))
						}
					}

					writeBool(t, "flip")
				}
			}

			writeString(settingsMap, "time_window")
			writeString(settingsMap, "title")
			writeString(settingsMap, "title_format")
			writeString(settingsMap, "trend")
			writeString(settingsMap, "type")
			writeBool(settingsMap, "use_default")
			writeString(settingsMap, "value_type")
			writeString(settingsMap, "date_window")
			writeString(settingsMap, "graph_uuid")
			writeBool(settingsMap, "hide_xaxis")
			writeBool(settingsMap, "hide_yaxis")
			writeBool(settingsMap, "key_inline")
			writeString(settingsMap, "key_loc")
			writeInt(settingsMap, "key_size")
			writeBool(settingsMap, "key_wrap")
			writeString(settingsMap, "label")
			writeInt(settingsMap, "period")
			writeBool(settingsMap, "real_time")
			writeBool(settingsMap, "show_flags")
		}
	}

	s := b.String()
	return hashcode.String(s)
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

		dashWidgetSettingsAttrs := make(map[string]interface{}, 63)
		dashWidgetSettingsAttrs["account_id"] = widget.Settings.AccountID
		dashWidgetSettingsAttrs["acknowledged"] = widget.Settings.Acknowledged
		dashWidgetSettingsAttrs["algorithm"] = widget.Settings.Algorithm
		dashWidgetSettingsAttrs["autoformat"] = widget.Settings.Autoformat
		brs := make([]map[string]interface{}, 0)
		for _, br := range widget.Settings.BadRules {
			brAttrs := make(map[string]interface{}, 3)
			brAttrs["value"] = br.Value
			brAttrs["criterion"] = br.Criterion
			brAttrs["color"] = br.Color
			brs = append(brs, brAttrs)
		}
		dashWidgetSettingsAttrs["bad_rules"] = brs
		dashWidgetSettingsAttrs["body_format"] = widget.Settings.BodyFormat
		dashWidgetSettingsAttrs["caql"] = widget.Settings.Caql
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
		dashWidgetSettingsAttrs["display_markup"] = widget.Settings.DisplayMarkup
		dashWidgetSettingsAttrs["format"] = widget.Settings.Format
		dashWidgetSettingsAttrs["formula"] = widget.Settings.Formula
		dashWidgetSettingsAttrs["good_color"] = widget.Settings.GoodColor
		dashWidgetSettingsAttrs["layout"] = widget.Settings.Layout
		dashWidgetSettingsAttrs["layout_style"] = widget.Settings.LayoutStyle
		dashWidgetSettingsAttrs["limit"] = int(widget.Settings.Limit)
		dashWidgetSettingsAttrs["link_url"] = widget.Settings.LinkUrl
		dashWidgetSettingsAttrs["maintenance"] = widget.Settings.Maintenance
		dashWidgetSettingsAttrs["markup"] = widget.Settings.Markup
		dashWidgetSettingsAttrs["metric_display_name"] = widget.Settings.MetricDisplayName
		dashWidgetSettingsAttrs["metric_name"] = widget.Settings.MetricName
		dashWidgetSettingsAttrs["metric_type"] = widget.Settings.MetricType
		dashWidgetSettingsAttrs["min_age"] = widget.Settings.MinAge
		dashWidgetSettingsAttrs["off_hours"] = widget.Settings.OffHours
		dashWidgetSettingsAttrs["overlay_set_id"] = widget.Settings.OverlaySetID
		if widget.Settings.RangeHigh != nil {
			dashWidgetSettingsAttrs["range_high"] = int(*widget.Settings.RangeHigh)
		}
		if widget.Settings.RangeLow != nil {
			dashWidgetSettingsAttrs["range_low"] = int(*widget.Settings.RangeLow)
		}
		dashWidgetSettingsAttrs["resource_limit"] = widget.Settings.ResourceLimit
		dashWidgetSettingsAttrs["resource_usage"] = widget.Settings.ResourceUsage
		dashWidgetSettingsAttrs["search"] = widget.Settings.Search
		dashWidgetSettingsAttrs["severity"] = widget.Settings.Severity
		if widget.Settings.ShowValue != nil {
			dashWidgetSettingsAttrs["show_value"] = *widget.Settings.ShowValue
		}
		dashWidgetSettingsAttrs["size"] = widget.Settings.Size
		dashWidgetSettingsAttrs["text_align"] = widget.Settings.TextAlign
		dashWidgetSettingsAttrs["tag_filter_set"] = widget.Settings.TagFilterSet
		dashWidgetSettingsAttrs["threshold"] = widget.Settings.Threshold
		if widget.Settings.Thresholds != nil {
			t := make([]map[string]interface{}, 0, 1)
			th := make(map[string]interface{}, 3)
			th["colors"] = widget.Settings.Thresholds.Colors
			th["values"] = widget.Settings.Thresholds.Values
			th["flip"] = widget.Settings.Thresholds.Flip
			t = append(t, th)
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
	_ = d.Set("widget", schema.NewSet(hashWidgets, []interface{}{widgets}))

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
	_ = d.Set("options", options)

	gridLayoutAttrs := make(map[string]interface{}, 2)
	gridLayoutAttrs["width"] = dash.GridLayout.Width
	gridLayoutAttrs["height"] = dash.GridLayout.Height
	_ = d.Set("grid_layout", gridLayoutAttrs)
	_ = d.Set("account_default", dash.AccountDefault)
	_ = d.Set("shared", dash.Shared)
	_ = d.Set("title", dash.Title)
	_ = d.Set("uuid", dash.UUID)

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
	_ = d.Set("uuid", "")

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
			w.Settings.RangeHigh = nil
			w.Settings.RangeLow = nil

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
				w.Settings.ShowValue = nil
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
					if v, found := sMap["autoformat"]; found {
						w.Settings.Autoformat = v.(bool)
					}
					if v, found := sMap["bad_rules"]; found {
						w.Settings.BadRules = make([]api.StateWidgetBadRulesSettings, 0)

						brList := v.([]interface{})
						for _, brElem := range brList {
							brAttrs := brElem.(map[string]interface{})
							br := api.StateWidgetBadRulesSettings{}
							if vv, found := brAttrs["value"]; found {
								br.Value = (vv.(string))
							}
							if vv, found := brAttrs["criterion"]; found {
								br.Criterion = (vv.(string))
							}
							if vv, found := brAttrs["color"]; found {
								br.Color = (vv.(string))
							}
							w.Settings.BadRules = append(w.Settings.BadRules, br)
						}
					}
					if v, found := sMap["body_format"]; found {
						w.Settings.BodyFormat = (v.(string))
					}
					if v, found := sMap["caql"]; found {
						w.Settings.Caql = (v.(string))
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
					if v, found := sMap["display_markup"]; found {
						w.Settings.DisplayMarkup = (v.(string))
					}
					if v, found := sMap["format"]; found {
						w.Settings.Format = (v.(string))
					}
					if v, found := sMap["formula"]; found {
						w.Settings.Formula = (v.(string))
					}
					if v, found := sMap["good_color"]; found {
						w.Settings.GoodColor = (v.(string))
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
					if v, found := sMap["layout_style"]; found {
						w.Settings.LayoutStyle = (v.(string))
					}
					if v, found := sMap["link_url"]; found {
						w.Settings.LinkUrl = (v.(string))
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
					if v, found := sMap["metric_type"]; found {
						w.Settings.MetricType = (v.(string))
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
						y, ok := wAttrs["type"]
						if ok && y.(string) == "gauge" {
							x := v.(int)
							w.Settings.RangeHigh = &x
						} else {
							w.Settings.RangeHigh = nil
						}
					} else {
						w.Settings.RangeHigh = nil
					}
					if v, found := sMap["range_low"]; found {
						y, ok := wAttrs["type"]
						if ok && y.(string) == "gauge" {
							x := v.(int)
							w.Settings.RangeLow = &x
						} else {
							w.Settings.RangeLow = nil
						}
					} else {
						w.Settings.RangeLow = nil
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
					if w.Type == "state" {
						if v, found := sMap["show_value"]; found {
							x := v.(bool)
							w.Settings.ShowValue = &x
						} else {
							x := false
							w.Settings.ShowValue = &x
						}
					} else {
						w.Settings.ShowValue = nil
					}
					if v, found := sMap["size"]; found {
						w.Settings.Size = (v.(string))
					}
					if v, found := sMap["text_align"]; found {
						w.Settings.TextAlign = (v.(string))
					}
					if v, found := sMap["threshold"]; found {
						w.Settings.Threshold = float32((v.(float64)))
					}
					if v, found := sMap["thresholds"]; found {
						t := api.ForecastGaugeWidgetThresholds{}

						// there will be only 1
						tList := v.(*schema.Set).List()
						if len(tList) > 0 {
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
						} else {
							w.Settings.Thresholds = nil
						}
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
