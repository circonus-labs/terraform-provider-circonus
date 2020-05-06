package circonus

import (
	"fmt"
	"time"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceMaintenance() *schema.Resource {
	return &schema.Resource{
		Create: maintenanceCreate,
		Read:   maintenanceRead,
		Update: maintenanceUpdate,
		Delete: maintenanceDelete,
		Exists: maintenanceExists,
		Importer: &schema.ResourceImporter{
			State: importStatePassthroughUnescape,
		},

		Schema: map[string]*schema.Schema{
			"account": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"check", "rule_set", "target"},
			},
			"check": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"account", "rule_set", "target"},
			},
			"rule_set": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"check", "account", "target"},
			},
			"target": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"check", "rule_set", "account"},
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severities": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateRegexp(ruleSetContainsAttr, `[1-5]`),
				},
			},
			"start": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},
			"stop": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.ValidateRFC3339TimeString,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func maintenanceCreate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	m := newMaintenance()

	if err := m.ParseConfig(d); err != nil {
		return errwrap.Wrapf("error parsing maintenance schema during create: {{err}}", err)
	}

	if err := m.Create(ctxt); err != nil {
		return errwrap.Wrapf("error creating maintenance: {{err}}", err)
	}

	d.SetId(m.CID)

	return maintenanceRead(d, meta)
}

func maintenanceExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	m, err := ctxt.client.FetchMaintenanceWindow(api.CIDType(&cid))
	if err != nil {
		return false, err
	}

	if m.CID == "" {
		return false, nil
	}

	return true, nil
}

func maintenanceRead(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	m, err := loadMaintenance(ctxt, api.CIDType(&cid))
	if err != nil {
		return err
	}

	d.SetId(m.CID)
	if m.Type == "account" {
		_ = d.Set("account", m.Item)
	} else if m.Type == "rule_set" {
		_ = d.Set("rule_set", m.Item)
	} else if m.Type == "check" {
		_ = d.Set("check", m.Item)
	} else if m.Type == "host" {
		_ = d.Set("target", m.Item)
	}

	_ = d.Set("notes", m.Notes)
	_ = d.Set("severities", m.Severities.([]interface{}))
	start := time.Unix(int64(m.Start), 0)
	stop := time.Unix(int64(m.Stop), 0)

	_ = d.Set("start", start.Format(time.RFC3339))
	_ = d.Set("stop", stop.Format(time.RFC3339))
	tags := make([]interface{}, 0)
	if len(m.Tags) > 0 {
		for _, t := range m.Tags {
			tags = append(tags, t)
		}
	}
	_ = d.Set("tags", tags)
	return nil
}

func maintenanceUpdate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	m := newMaintenance()

	if err := m.ParseConfig(d); err != nil {
		return err
	}

	m.CID = d.Id()

	if err := m.Update(ctxt); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to update maintenance %q: {{err}}", d.Id()), err)
	}

	return maintenanceRead(d, meta)
}

func maintenanceDelete(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	if _, err := ctxt.client.DeleteMaintenanceWindowByCID(api.CIDType(&cid)); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to delete rule set %q: {{err}}", d.Id()), err)
	}

	d.SetId("")

	return nil
}

type circonusMaintenance struct {
	api.Maintenance
}

func newMaintenance() circonusMaintenance {
	m := circonusMaintenance{
		Maintenance: *api.NewMaintenanceWindow(),
	}

	return m
}

func loadMaintenance(ctxt *providerContext, cid api.CIDType) (circonusMaintenance, error) {
	var m circonusMaintenance
	cm, err := ctxt.client.FetchMaintenanceWindow(cid)
	if err != nil {
		return circonusMaintenance{}, err
	}
	m.Maintenance = *cm

	return m, nil
}

func (m *circonusMaintenance) ParseConfig(d *schema.ResourceData) error {

	if v, found := d.GetOk("account"); found && v.(string) != "" {
		m.Item = v.(string)
		m.Type = "account"
	}

	if v, found := d.GetOk("check"); found && v.(string) != "" {
		m.Item = v.(string)
		m.Type = "check"
	}

	if v, found := d.GetOk("rule_set"); found && v.(string) != "" {
		m.Item = v.(string)
		m.Type = "rule_set"
	}

	if v, found := d.GetOk("target"); found && v.(string) != "" {
		m.Item = v.(string)
		m.Type = "host"
	}

	if v, found := d.GetOk("notes"); found && v.(string) != "" {
		m.Notes = v.(string)
	}

	if v, found := d.GetOk("severities"); found && len(v.([]interface{})) > 0 {
		m.Severities = make([]string, 0)
		for _, s := range v.([]interface{}) {
			m.Severities = append(m.Severities.([]string), s.(string))
		}
	}

	if v, found := d.GetOk("start"); found && v.(string) != "" {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err == nil {
			m.Start = uint(t.Unix())
		}
	}

	if v, found := d.GetOk("stop"); found && v.(string) != "" {
		t, err := time.Parse(time.RFC3339, v.(string))
		if err == nil {
			m.Stop = uint(t.Unix())
		}
	}

	if v, found := d.GetOk("tags"); found && len(v.([]string)) > 0 {
		m.Tags = derefStringList(flattenSet(v.(*schema.Set)))
	}

	if err := m.Validate(); err != nil {
		return err
	}

	return nil
}

func (m *circonusMaintenance) Create(ctxt *providerContext) error {
	cm, err := ctxt.client.CreateMaintenanceWindow(&m.Maintenance)
	if err != nil {
		return err
	}

	m.CID = cm.CID

	return nil
}

func (m *circonusMaintenance) Update(ctxt *providerContext) error {
	_, err := ctxt.client.UpdateMaintenanceWindow(&m.Maintenance)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to update maintenance %s: {{err}}", m.CID), err)
	}

	return nil
}

func (m *circonusMaintenance) Validate() error {
	return nil
}
