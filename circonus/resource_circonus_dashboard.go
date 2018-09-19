package circonus

import (
	"fmt"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	dashboardNameAttr             = "name"
	dashboardSharedAttr           = "shared"
	dashboardAccountDefaultAttr   = "account_default"
	dashboardActiveAttr           = "active"
	dashboardGridLayoutAttr       = "grid_layout"
	dashboardGridLayoutHeightAttr = "height"
	dashboardGridLayoutWidthAttr  = "width"
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

		Schema: convertToHelperSchema(graphDescriptions, map[schemaAttr]*schema.Schema{
			dashboardNameAttr: {
				Type:     schema.TypeString,
				Required: true,
			},
			dashboardAccountDefaultAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  defaultDashboardAccountDefault,
			},
			dashboardSharedAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  defaultDashboardShared,
			},
			dashboardActiveAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  defaultDashboardActive,
			},
			dashboardGridLayoutAttr: {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: convertToHelperSchema(graphMetricClusterDescriptions, map[schemaAttr]*schema.Schema{
						dashboardGridLayoutHeightAttr: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  defaultDashboardGridHeight,
							ValidateFunc: validateFuncs(
								validateIntMin(dashboardGridLayoutHeightAttr, 4),
								validateIntMax(dashboardGridLayoutHeightAttr, 26),
							),
						},
						dashboardGridLayoutWidthAttr: {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  defaultDashboardGridWidth,
							ValidateFunc: validateFuncs(
								validateIntMin(dashboardGridLayoutWidthAttr, 4),
								validateIntMax(dashboardGridLayoutWidthAttr, 26),
							),
						},
					}),
				},
			},
		}),
	}
}

func dashboardCreate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	g := newDashboard()
	if err := g.ParseConfig(d); err != nil {
		return errwrap.Wrapf("error parsing dashboard schema during create: {{err}}", err)
	}

	if err := g.Create(ctxt); err != nil {
		return errwrap.Wrapf("error creating dashboard: {{err}}", err)
	}

	d.SetId(g.CID)

	return graphRead(d, meta)
}

func dashboardRead(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	g, err := loadDashboard(ctxt, api.CIDType(&cid))
	if err != nil {
		return err
	}

	d.SetId(g.CID)
	d.Set(dashboardActiveAttr, g.Active)
	d.Set(dashboardAccountDefaultAttr, g.AccountDefault)
	d.Set(dashboardSharedAttr, g.Shared)
	d.Set(dashboardNameAttr, g.Title)

	if err := d.Set(dashboardGridLayoutAttr, flattenDashboardGridLayout(g.GridLayout)); err != nil {
		return err
	}

	return nil
}

func dashboardUpdate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	g := newDashboard()
	if err := g.ParseConfig(d); err != nil {
		return err
	}

	g.CID = d.Id()
	if err := g.Update(ctxt); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to update dashboard %q: {{err}}", d.Id()), err)
	}

	return graphRead(d, meta)
}

func dashboardExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	g, err := ctxt.client.FetchDashboard(api.CIDType(&cid))
	if err != nil {
		if strings.Contains(err.Error(), defaultCirconus404ErrorString) {
			return false, nil
		}

		return false, err
	}

	if g.CID == "" {
		return false, nil
	}

	return true, nil
}

func dashboardDelete(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	if _, err := ctxt.client.DeleteDashboardByCID(api.CIDType(&cid)); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to delete dashboard %q: {{err}}", d.Id()), err)
	}

	d.SetId("")

	return nil
}

func (g *circonusDashboard) Update(ctxt *providerContext) error {
	_, err := ctxt.client.UpdateDashboard(&g.Dashboard)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to update dashboard %s: {{err}}", g.CID), err)
	}

	return nil
}

func (g *circonusDashboard) Create(ctxt *providerContext) error {
	ng, err := ctxt.client.CreateDashboard(&g.Dashboard)
	if err != nil {
		return err
	}

	g.CID = ng.CID

	return nil
}

func (g *circonusDashboard) ParseConfig(d *schema.ResourceData) error {
	g.Title = d.Get(dashboardNameAttr).(string)
	g.AccountDefault = d.Get(dashboardAccountDefaultAttr).(bool)
	g.Active = d.Get(dashboardActiveAttr).(bool)
	g.Shared = d.Get(dashboardSharedAttr).(bool)
	g.GridLayout = expandDashboardGridLayout(d)

	return nil
}

type circonusDashboard struct {
	api.Dashboard
}

func newDashboard() circonusDashboard {
	g := circonusDashboard{
		Dashboard: *api.NewDashboard(),
	}

	return g
}

func loadDashboard(ctxt *providerContext, cid api.CIDType) (circonusDashboard, error) {
	var g circonusDashboard
	ng, err := ctxt.client.FetchDashboard(cid)
	if err != nil {
		return circonusDashboard{}, err
	}
	g.Dashboard = *ng

	return g, nil
}

func expandDashboardGridLayout(d *schema.ResourceData) api.DashboardGridLayout {
	grid := d.Get(dashboardGridLayoutAttr).([]interface{})
	layout := grid[0].(map[string]interface{})

	gridLayout := api.DashboardGridLayout{
		Height: uint(layout[dashboardGridLayoutHeightAttr].(int)),
		Width:  uint(layout[dashboardGridLayoutWidthAttr].(int)),
	}

	return gridLayout
}

func flattenDashboardGridLayout(grid api.DashboardGridLayout) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	item := make(map[string]interface{})

	item[dashboardGridLayoutHeightAttr] = grid.Height
	item[dashboardGridLayoutWidthAttr] = grid.Width

	result = append(result, item)

	return result
}
