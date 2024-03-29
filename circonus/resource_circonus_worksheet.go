package circonus

import (
	"context"
	"fmt"
	"strings"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	workspaceTitleAttr        = "title"
	workspaceDescriptionAttr  = "description"
	workspaceFavoriteAttr     = "favorite"
	workspaceNotesAttr        = "notes"
	workspaceTagsAttr         = "tags"
	workspaceGraphsAttr       = "graphs"
	workspaceSmartQueriesAttr = "smart_queries"

	queryNameAttr  = "name"
	queryQueryAttr = "query"
	queryOrderAttr = "order"
)

var worksheetDescriptions = attrDescrs{
	workspaceTitleAttr:        "",
	workspaceDescriptionAttr:  "",
	workspaceFavoriteAttr:     "",
	workspaceNotesAttr:        "",
	workspaceTagsAttr:         "",
	workspaceGraphsAttr:       "",
	workspaceSmartQueriesAttr: "",
}

var worksheetSmartQueryDescriptions = attrDescrs{
	queryNameAttr:  "",
	queryQueryAttr: "",
	queryOrderAttr: "",
}

func resourceWorksheet() *schema.Resource {
	return &schema.Resource{
		CreateContext: worksheetCreate,
		ReadContext:   worksheetRead,
		UpdateContext: worksheetUpdate,
		DeleteContext: worksheetDelete,
		Exists:        worksheetExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: convertToHelperSchema(worksheetDescriptions, map[schemaAttr]*schema.Schema{
			workspaceTitleAttr: {
				Type:     schema.TypeString,
				Required: true,
			},

			workspaceDescriptionAttr: {
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: suppressWhitespace,
			},

			workspaceFavoriteAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  defaultWorkspaceFavourite,
			},

			workspaceNotesAttr: {
				Type:     schema.TypeString,
				Optional: true,
			},

			workspaceGraphsAttr: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			workspaceSmartQueriesAttr: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: convertToHelperSchema(worksheetSmartQueryDescriptions, map[schemaAttr]*schema.Schema{
						queryNameAttr: {
							Type:     schema.TypeString,
							Required: true,
						},
						queryQueryAttr: {
							Type:     schema.TypeString,
							Required: true,
						},
						queryOrderAttr: {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					}),
				},
			},
			workspaceTagsAttr: tagMakeConfigSchema(workspaceTagsAttr),
		}),
	}
}

func worksheetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctxt := meta.(*providerContext)

	g := newWorksheet()
	if err := g.ParseConfig(d); err != nil {
		return diag.FromErr(fmt.Errorf("parsing worksheet schema during create: %w", err))
	}

	if err := g.Create(ctxt); err != nil {
		return diag.FromErr(fmt.Errorf("creating worksheet: %w", err))
	}

	d.SetId(g.CID)

	return worksheetRead(ctx, d, meta)
}

func worksheetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	ctxt := meta.(*providerContext)

	cid := d.Id()
	w, err := loadWorksheet(ctxt, api.CIDType(&cid))
	if err != nil {
		return diag.FromErr(fmt.Errorf("load worksheet: %w", err))
	}

	d.SetId(w.CID)

	_ = d.Set(workspaceTitleAttr, w.Title)
	_ = d.Set(workspaceDescriptionAttr, w.Description)
	_ = d.Set(workspaceFavoriteAttr, w.Favorite)
	_ = d.Set(workspaceNotesAttr, w.Notes)

	if err := d.Set(workspaceGraphsAttr, worksheetGraphsToState(apiToWorksheetGraphs(w.Graphs))); err != nil {
		return diag.FromErr(fmt.Errorf("unable to store worksheet %q attribute: %w", workspaceTagsAttr, err))
	}

	if err := d.Set(workspaceTagsAttr, tagsToState(apiToTags(w.Tags))); err != nil {
		return diag.FromErr(fmt.Errorf("unable to store worksheet %q attribute: %w", workspaceTagsAttr, err))
	}

	smartQueries := make([]map[string]interface{}, 0, len(w.SmartQueries))

	for _, query := range w.SmartQueries {
		newQuery := map[string]interface{}{
			"name":  query.Name,
			"query": query.Query,
			"order": query.Order,
		}

		smartQueries = append(smartQueries, newQuery)
	}

	if err := d.Set(workspaceSmartQueriesAttr, smartQueries); err != nil {
		return diag.FromErr(fmt.Errorf("unable to store worksheet %q attribute: %w", workspaceSmartQueriesAttr, err))
	}

	return diags
}

func worksheetExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	w, err := ctxt.client.FetchWorksheet(api.CIDType(&cid))
	if err != nil {
		if strings.Contains(err.Error(), defaultCirconus404ErrorString) {
			return false, nil
		}

		return false, err
	}

	if w.CID == "" {
		return false, nil
	}

	return true, nil
}

func worksheetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctxt := meta.(*providerContext)
	w := newWorksheet()
	if err := w.ParseConfig(d); err != nil {
		return diag.FromErr(fmt.Errorf("parse worksheet config: %w", err))
	}

	w.CID = d.Id()
	if err := w.Update(ctxt); err != nil {
		return diag.FromErr(fmt.Errorf("unable to update worksheet %q: %w", d.Id(), err))
	}

	return worksheetRead(ctx, d, meta)
}

func worksheetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	ctxt := meta.(*providerContext)

	cid := d.Id()
	if _, err := ctxt.client.DeleteWorksheetByCID(api.CIDType(&cid)); err != nil {
		return diag.FromErr(fmt.Errorf("unable to delete worksheet %q: %w", d.Id(), err))
	}

	d.SetId("")

	return diags
}

func (w *circonusWorksheet) Create(ctxt *providerContext) error {
	nw, err := ctxt.client.CreateWorksheet(&w.Worksheet)
	if err != nil {
		return err
	}

	w.CID = nw.CID

	return nil
}

func (w *circonusWorksheet) Update(ctxt *providerContext) error {
	_, err := ctxt.client.UpdateWorksheet(&w.Worksheet)
	if err != nil {
		return fmt.Errorf("Unable to update worksheet %s: %w", w.CID, err)
	}

	return nil
}

func loadWorksheet(ctxt *providerContext, cid api.CIDType) (circonusWorksheet, error) {
	var w circonusWorksheet
	nw, err := ctxt.client.FetchWorksheet(cid)
	if err != nil {
		return circonusWorksheet{}, err
	}
	w.Worksheet = *nw

	return w, nil
}

func (w *circonusWorksheet) ParseConfig(d *schema.ResourceData) error {
	w.Title = d.Get(workspaceTitleAttr).(string)

	if v, ok := d.GetOk(workspaceDescriptionAttr); ok {
		desc := v.(string)
		w.Description = &desc
	}

	if v, ok := d.GetOk(workspaceNotesAttr); ok {
		notes := v.(string)
		w.Notes = &notes
	}

	if v, found := d.GetOk(workspaceTagsAttr); found {
		w.Tags = derefStringList(flattenSet(v.(*schema.Set)))
	}

	if v, found := d.GetOk(workspaceGraphsAttr); found {
		graphs := derefStringList(flattenSet(v.(*schema.Set)))
		var workspaceGraphs []api.WorksheetGraph
		for _, graph := range graphs {
			workspaceGraphs = append(workspaceGraphs, api.WorksheetGraph{
				GraphCID: graph,
			})
		}

		w.Graphs = workspaceGraphs
	}

	if v, found := d.GetOk(workspaceSmartQueriesAttr); found {
		queriesList := v.(*schema.Set).List()
		smaryQueries := make([]api.WorksheetSmartQuery, 0, len(queriesList))

		for _, queryListRaw := range queriesList {
			var query api.WorksheetSmartQuery
			queryAttrs := queryListRaw.(map[string]interface{})

			if v, found := queryAttrs[queryNameAttr]; found {
				query.Name = v.(string)
			}

			if v, found := queryAttrs[queryQueryAttr]; found {
				query.Query = v.(string)
			}

			if v, found := queryAttrs[queryOrderAttr]; found {
				orderList := v.(*schema.Set).List()
				query.Order = make([]string, len(orderList))
				for _, s := range orderList {
					query.Order = append(query.Order, s.(string))
				}
			}
			smaryQueries = append(smaryQueries, query)
		}

		w.SmartQueries = smaryQueries
	}

	return nil
}

type circonusWorksheet struct {
	api.Worksheet
}

func newWorksheet() circonusWorksheet {
	w := circonusWorksheet{
		Worksheet: *api.NewWorksheet(),
	}

	return w
}

func apiToWorksheetGraphs(graphs []api.WorksheetGraph) []string {
	workSheetGraphs := make([]string, 0, len(graphs))

	for _, v := range graphs {
		workSheetGraphs = append(workSheetGraphs, v.GraphCID)
	}
	return workSheetGraphs
}

func worksheetGraphsToState(graphs []string) *schema.Set {
	graphsSet := schema.NewSet(schema.HashString, nil)
	for i := range graphs {
		graphsSet.Add(graphs[i])
	}
	return graphsSet
}
