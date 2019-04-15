package circonus

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/schema"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func resourceOverlaySet() *schema.Resource {
	return &schema.Resource{
		Create: overlaySetCreate,
		Read:   overlaySetRead,
		Update: overlaySetUpdate,
		Delete: overlaySetDelete,
		Exists: overlaySetExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"graph_cid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"title": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"data_opts": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"graph_title": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"graph_uuid": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"x_shift": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"ui_specs": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"decouple": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"label": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"z": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func overlaySetCreate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	o := newOverlaySet()

	if err := o.ParseConfig(d); err != nil {
		return errwrap.Wrapf("error parsing graph schema during create: {{err}}", err)
	}

	if err := o.Create(ctxt); err != nil {
		return errwrap.Wrapf("error creating graph: {{err}}", err)
	}

	return overlaySetRead(d, meta)
}

func overlaySetExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ctxt := meta.(*providerContext)

	if graph_cid, found := d.GetOk("graph_cid"); found {
		id := d.Id()
		s := graph_cid.(string)
		g, err := ctxt.client.FetchGraph(api.CIDType(&s))
		if err != nil {
			if strings.Contains(err.Error(), defaultCirconus404ErrorString) {
				return false, nil
			}

			return false, err
		}

		if g.CID == "" {
			return false, nil
		}

		if g.OverlaySets == nil {
			return false, nil
		}

		if _, ok := (*g.OverlaySets)[id]; ok {
			return true, nil
		}

		return false, nil
	}
	return false, nil
}

// graphRead pulls data out of the Graph object and stores it into the
// appropriate place in the statefile.
func overlaySetRead(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	id := d.Id()
	if graph_cid, found := d.GetOk("graph_cid"); found {
		s := graph_cid.(string)
		g, err := loadOverlaySet(ctxt, api.CIDType(&s), id)
		if err != nil {
			return err
		}

		d.SetId(g.GraphOverlaySet.ID)
		d.Set("graph_cid", graph_cid)
		d.Set("title", g.GraphOverlaySet.Title)

		uiSpecs := make(map[string]interface{}, 5)
		uiSpecs["decouple"] = g.GraphOverlaySet.UISpecs.Decouple
		uiSpecs["label"] = g.GraphOverlaySet.UISpecs.Label
		uiSpecs["type"] = g.GraphOverlaySet.UISpecs.Type
		if g.GraphOverlaySet.UISpecs.Z != nil {
			uiSpecs["z"] = fmt.Sprintf("%d", *g.GraphOverlaySet.UISpecs.Z)
		}

		set := make([]map[string]interface{}, 1)
		set[0] = uiSpecs
		d.Set("ui_specs", set)

		dataOpts := make(map[string]interface{}, 3)
		dataOpts["graph_title"] = g.GraphOverlaySet.DataOpts.GraphTitle
		dataOpts["graph_uuid"] = g.GraphOverlaySet.DataOpts.GraphUUID
		dataOpts["x_shift"] = g.GraphOverlaySet.DataOpts.XShift

		set = make([]map[string]interface{}, 1)
		set[0] = dataOpts
		d.Set("data_opt", set)

		return nil
	}
	return errors.New(fmt.Sprintf("graph_cid field is required for %q", d.Id()))
}

func overlaySetUpdate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	g := newOverlaySet()
	if err := g.ParseConfig(d); err != nil {
		return err
	}

	g.GraphOverlaySet.ID = d.Id()

	if err := g.Update(ctxt); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to update graph %q: {{err}}", d.Id()), err)
	}

	return overlaySetRead(d, meta)
}

func overlaySetDelete(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	if graph_cid, found := d.GetOk("graph_cid"); found {
		id := d.Id()
		s := graph_cid.(string)
		var graph *api.Graph
		var err error
		if graph, err = ctxt.client.FetchGraph(api.CIDType(&s)); err != nil {
			return errwrap.Wrapf(fmt.Sprintf("unable to delete overlay set %q: {{err}}", d.Id()), err)
		}

		if graph.OverlaySets != nil {
			delete(*graph.OverlaySets, id)
		}

		if _, err := ctxt.client.UpdateGraph(graph); err != nil {
			return errwrap.Wrapf(fmt.Sprintf("unable to delete overlay set %q: {{err}}", d.Id()), err)
		}

		d.SetId("")
	}

	return nil
}

type circonusOverlaySet struct {
	GraphCID        string
	GraphOverlaySet api.GraphOverlaySet
}

func newOverlaySet() circonusOverlaySet {
	g := circonusOverlaySet{
		GraphCID:        "",
		GraphOverlaySet: api.GraphOverlaySet{},
	}

	return g
}

func loadOverlaySet(ctxt *providerContext, graph_cid api.CIDType, set_id string) (circonusOverlaySet, error) {
	var g circonusOverlaySet
	ng, err := ctxt.client.FetchGraph(graph_cid)
	if err != nil {
		return circonusOverlaySet{}, err
	}
	if ng.OverlaySets == nil {
		return circonusOverlaySet{}, nil
	}

	g.GraphOverlaySet = (*ng.OverlaySets)[set_id]
	g.GraphCID = *graph_cid
	return g, nil
}

// ParseConfig reads Terraform config data and stores the information into a
// Circonus OverlaySet object.  ParseConfig and graphRead() must be kept in sync.
func (g *circonusOverlaySet) ParseConfig(d *schema.ResourceData) error {

	if v, found := d.GetOk("title"); found {
		g.GraphOverlaySet.Title = v.(string)
	}
	if v, found := d.GetOk("ui_specs"); found {
		uiSpecsList := v.(*schema.Set).List()
		for _, uiSpec := range uiSpecsList {
			uiSpecMap := newInterfaceMap(uiSpec.(map[string]interface{}))
			if v, found := uiSpecMap["decouple"]; found {
				g.GraphOverlaySet.UISpecs.Decouple = v.(bool)
			}
			if v, found := uiSpecMap["label"]; found {
				g.GraphOverlaySet.UISpecs.Label = v.(string)
			}
			if v, found := uiSpecMap["type"]; found {
				g.GraphOverlaySet.UISpecs.Type = v.(string)
			}
			if v, found := uiSpecMap["z"]; found {
				i64, _ := strconv.ParseInt(v.(string), 10, 64)
				g.GraphOverlaySet.UISpecs.Z = new(int)
				*g.GraphOverlaySet.UISpecs.Z = int(i64)
			}
		}
	}
	if v, found := d.GetOk("data_opts"); found {
		dataOptsList := v.(*schema.Set).List()
		for _, dataOpt := range dataOptsList {
			dataOptMap := newInterfaceMap(dataOpt.(map[string]interface{}))
			if v, found := dataOptMap["graph_title"]; found {
				g.GraphOverlaySet.DataOpts.GraphTitle = v.(string)
			}
			if v, found := dataOptMap["graph_uuid"]; found {
				g.GraphOverlaySet.DataOpts.GraphUUID = v.(string)
			}
			if v, found := dataOptMap["x_shift"]; found {
				g.GraphOverlaySet.DataOpts.XShift = v.(string)
			}
		}
	}
	if v, found := d.GetOk("graph_cid"); found {
		g.GraphCID = v.(string)
	}
	return nil
}

func (g *circonusOverlaySet) Create(ctxt *providerContext) error {

	gg, err := ctxt.client.FetchGraph(api.CIDType(&g.GraphCID))
	if err != nil {
		return err
	}

	set := make(map[string]api.GraphOverlaySet)

	g.GraphOverlaySet.ID = randStringBytes(6)
	g.GraphOverlaySet.UISpecs.ID = g.GraphOverlaySet.ID

	if gg.OverlaySets == nil {
		gg.OverlaySets = &set
	}

	(*gg.OverlaySets)[g.GraphOverlaySet.ID] = g.GraphOverlaySet

	_, err = ctxt.client.UpdateGraph(gg)
	if err != nil {
		return err
	}

	return nil
}

func (g *circonusOverlaySet) Update(ctxt *providerContext) error {
	gg, err := ctxt.client.FetchGraph(api.CIDType(&g.GraphCID))
	if err != nil {
		return err
	}

	g.GraphOverlaySet.UISpecs.ID = g.GraphOverlaySet.ID

	(*gg.OverlaySets)[g.GraphOverlaySet.ID] = g.GraphOverlaySet

	_, err = ctxt.client.UpdateGraph(gg)
	if err != nil {
		return err
	}

	return nil
}

func (g *circonusOverlaySet) Validate() error {
	return nil
}
