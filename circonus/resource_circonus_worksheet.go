package circonus

import (
	"fmt"
	"strings"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	workspaceTitleAttr        = "title"
	workspaceDescriptionAttr  = "description"
	workspaceFavouriteAttr    = "favourite"
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
	workspaceFavouriteAttr:    "",
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
		Create: worksheetCreate,
		Read:   worksheetRead,
		Update: worksheetUpdate,
		Delete: worksheetDelete,
		Exists: worksheetExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

			workspaceFavouriteAttr: {
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

func worksheetCreate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	g := newWorksheet()
	if err := g.ParseConfig(d); err != nil {
		return errwrap.Wrapf("error parsing graph schema during create: {{err}}", err)
	}

	if err := g.Create(ctxt); err != nil {
		return errwrap.Wrapf("error creating graph: {{err}}", err)
	}

	d.SetId(g.CID)

	return worksheetRead(d, meta)
}

func worksheetRead(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	w, err := loadWorksheet(ctxt, api.CIDType(&cid))
	if err != nil {
		return err
	}

	d.SetId(w.CID)

	d.Set(workspaceTitleAttr, w.Title)
	d.Set(workspaceDescriptionAttr, w.Description)
	d.Set(workspaceFavouriteAttr, w.Favorite)
	d.Set(workspaceNotesAttr, w.Notes)

	if err := d.Set(workspaceGraphsAttr, worksheetGraphsToState(apiToWorksheetGraphs(w.Graphs))); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store workspace %q attribute: {{err}}", workspaceTagsAttr), err)
	}

	if err := d.Set(workspaceTagsAttr, tagsToState(apiToTags(w.Tags))); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store workspace %q attribute: {{err}}", workspaceTagsAttr), err)
	}

	var smartQueries []map[string]interface{}

	for _, query := range w.SmartQueries {

		newQuery := map[string]interface{}{
			"name":  query.Name,
			"query": query.Query,
			"order": query.Order,
		}

		smartQueries = append(smartQueries, newQuery)
	}

	if err := d.Set(workspaceSmartQueriesAttr, smartQueries); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to store worksheet %q attribute: {{err}}", workspaceSmartQueriesAttr), err)
	}

	return nil
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

func worksheetUpdate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	w := newWorksheet()
	if err := w.ParseConfig(d); err != nil {
		return err
	}

	w.CID = d.Id()
	if err := w.Update(ctxt); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to update worksheet %q: {{err}}", d.Id()), err)
	}

	return worksheetRead(d, meta)
}

func worksheetDelete(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	if _, err := ctxt.client.DeleteWorksheetByCID(api.CIDType(&cid)); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to delete worksheet %q: {{err}}", d.Id()), err)
	}

	d.SetId("")

	return nil
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
		return errwrap.Wrapf(fmt.Sprintf("Unable to update worksheet %s: {{err}}", w.CID), err)
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
	var workSheetGraphs []string

	for _, v := range graphs {
		workSheetGraphs = append(workSheetGraphs, v.GraphCID)
	}
	return workSheetGraphs
}

func worksheetGraphsToState(graphs []string) *schema.Set {
	graphsSet := schema.NewSet(schema.HashString, nil)
	for i := range graphs {
		graphsSet.Add(string(graphs[i]))
	}
	return graphsSet
}
