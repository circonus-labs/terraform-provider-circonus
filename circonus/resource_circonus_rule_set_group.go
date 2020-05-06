package circonus

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strconv"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRuleSetGroup() *schema.Resource {
	return &schema.Resource{
		Create: ruleSetGroupCreate,
		Read:   ruleSetGroupRead,
		Update: ruleSetGroupUpdate,
		Delete: ruleSetGroupDelete,
		Exists: ruleSetGroupExists,
		Importer: &schema.ResourceImporter{
			State: importStatePassthroughUnescape,
		},
		Schema: map[string]*schema.Schema{
			"notify": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sev1": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sev2": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sev3": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sev4": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"sev5": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"formula": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"expression": {
							Type:     schema.TypeString,
							Required: true,
						},
						"raise_severity": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validateFuncs(
								validateIntMax("raise_severity", 5),
								validateIntMin("raise_severity", 1),
							),
						},
						"wait": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validateFuncs(
								validateIntMin("wait", 0),
							),
						},
					},
				},
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:     schema.TypeSet,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"index": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"rule_set": {
							Type:     schema.TypeString,
							Required: true,
						},
						"matching_severities": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
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

func ruleSetGroupCreate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	rsg := newRuleSetGroup()

	if err := rsg.ParseConfig(d); err != nil {
		return errwrap.Wrapf("error parsing rule set group schema during create: {{err}}", err)
	}

	if err := rsg.Create(ctxt); err != nil {
		return errwrap.Wrapf("error creating rule set group: {{err}}", err)
	}

	d.SetId(rsg.CID)

	return ruleSetGroupRead(d, meta)
}

func ruleSetGroupExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	rsg, err := ctxt.client.FetchRuleSetGroup(api.CIDType(&cid))
	if err != nil {
		return false, err
	}

	if rsg.CID == "" {
		return false, nil
	}

	return true, nil
}

// ruleSetRead pulls data out of the RuleSet object and stores it into the
// appropriate place in the statefile.
func ruleSetGroupRead(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	rsg, err := loadRuleSetGroup(ctxt, api.CIDType(&cid))
	if err != nil {
		return err
	}

	d.SetId(rsg.CID)

	formulas := make([]interface{}, 0, 1)
	for _, formula := range rsg.Formulas {
		f := make(map[string]interface{}, 3)
		f["expression"] = formula.Expression.(string)
		t := reflect.TypeOf(formula.RaiseSeverity)
		if t.String() == "uint" {
			f["raise_severity"] = int(formula.RaiseSeverity.(uint))
		} else if t.String() == "string" {
			s, _ := strconv.ParseInt(formula.RaiseSeverity.(string), 10, 32)
			f["raise_severity"] = s
		} else {
			f["raise_severity"] = int(formula.RaiseSeverity.(float64))
		}
		f["wait"] = int(formula.Wait)
		formulas = append(formulas, f)
	}
	_ = d.Set("formula", formulas)
	_ = d.Set("name", rsg.Name)
	n := make([]interface{}, 0)
	notify := make(map[string]interface{})
	notify["sev1"] = rsg.ContactGroups[1]
	notify["sev2"] = rsg.ContactGroups[2]
	notify["sev3"] = rsg.ContactGroups[3]
	notify["sev4"] = rsg.ContactGroups[4]
	notify["sev5"] = rsg.ContactGroups[5]
	n = append(n, notify)

	_ = d.Set("notify", n)
	conditions := make([]interface{}, 0, len(rsg.RuleSetConditions))
	for idx, c := range rsg.RuleSetConditions {
		cond := make(map[string]interface{}, 2)
		cond["index"] = idx + 1
		cond["rule_set"] = c.RuleSetCID
		cond["matching_severities"] = c.MatchingSeverities
		conditions = append(conditions, cond)
	}
	_ = d.Set("condition", conditions)

	tags := make([]interface{}, 0)
	if len(rsg.Tags) > 0 {
		for _, t := range rsg.Tags {
			tags = append(tags, t)
		}
	}
	_ = d.Set("tags", tags)

	return nil
}

func ruleSetGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)
	rs := newRuleSetGroup()

	if err := rs.ParseConfig(d); err != nil {
		return err
	}

	rs.CID = d.Id()

	if err := rs.Update(ctxt); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to update rule set group %q: {{err}}", d.Id()), err)
	}

	return ruleSetGroupRead(d, meta)
}

func ruleSetGroupDelete(d *schema.ResourceData, meta interface{}) error {
	ctxt := meta.(*providerContext)

	cid := d.Id()
	if _, err := ctxt.client.DeleteRuleSetGroupByCID(api.CIDType(&cid)); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to delete rule set group %q: {{err}}", d.Id()), err)
	}

	d.SetId("")

	return nil
}

type circonusRuleSetGroup struct {
	api.RuleSetGroup
}

func newRuleSetGroup() circonusRuleSetGroup {
	rsg := circonusRuleSetGroup{
		RuleSetGroup: *api.NewRuleSetGroup(),
	}

	rsg.ContactGroups = make(map[uint8][]string, config.NumSeverityLevels)
	for i := uint8(0); i < config.NumSeverityLevels; i++ {
		rsg.ContactGroups[i+1] = make([]string, 0, 1)
	}
	rsg.Tags = make([]string, 0)

	return rsg
}

func loadRuleSetGroup(ctxt *providerContext, cid api.CIDType) (circonusRuleSetGroup, error) {
	var rs circonusRuleSetGroup
	crs, err := ctxt.client.FetchRuleSetGroup(cid)
	if err != nil {
		return circonusRuleSetGroup{}, err
	}
	log.Printf("RuleSetGroup: %v\n", *crs)
	rs.RuleSetGroup = *crs

	return rs, nil
}

// func ruleSetGroupNotifyChecksum(v interface{}) int {
// 	b := &bytes.Buffer{}
// 	b.Grow(defaultHashBufSize)

// 	writeStringArray := func(m map[string]interface{}, attrName string) {
// 		if v, found := m[attrName]; found {
// 			a := v.([]string)
// 			if a != nil {
// 				sort.Strings(a)
// 				for _, s := range a {
// 					fmt.Fprint(b, strings.TrimSpace(s))
// 				}
// 			}
// 		}
// 	}

// 	m := v.(map[string]interface{})

// 	writeStringArray(m, "sev1")
// 	writeStringArray(m, "sev2")
// 	writeStringArray(m, "sev3")
// 	writeStringArray(m, "sev4")
// 	writeStringArray(m, "sev5")

// 	s := b.String()
// 	return hashcode.String(s)
// }

// func ruleSetGroupFormulasChecksum(v interface{}) int {
// 	b := &bytes.Buffer{}
// 	b.Grow(defaultHashBufSize)

// 	writeInt := func(m map[string]interface{}, attrName string) {
// 		if v, found := m[attrName]; found {
// 			i := v.(int)
// 			if i != 0 {
// 				fmt.Fprintf(b, "%x", i)
// 			}
// 		}
// 	}

// 	writeString := func(m map[string]interface{}, attrName string) {
// 		if v, found := m[attrName]; found {
// 			s := strings.TrimSpace(v.(string))
// 			if s != "" {
// 				fmt.Fprint(b, s)
// 			}
// 		}
// 	}

// 	m := v.([]map[string]interface{})
// 	for _, f := range m {
// 		writeString(f, "expression")
// 		writeInt(f, "raise_severity")
// 		writeInt(f, "wait")
// 	}

// 	s := b.String()
// 	return hashcode.String(s)
// }

// func ruleSetGroupConditionsChecksum(v interface{}) int {
// 	b := &bytes.Buffer{}
// 	b.Grow(defaultHashBufSize)

// 	writeInt := func(m map[string]interface{}, attrName string) {
// 		if v, found := m[attrName]; found {
// 			i := v.(int)
// 			if i != 0 {
// 				fmt.Fprintf(b, "%x", i)
// 			}
// 		}
// 	}

// 	writeString := func(m map[string]interface{}, attrName string) {
// 		if v, found := m[attrName]; found {
// 			s := strings.TrimSpace(v.(string))
// 			if s != "" {
// 				fmt.Fprint(b, s)
// 			}
// 		}
// 	}

// 	writeStringArray := func(m map[string]interface{}, attrName string) {
// 		if v, found := m[attrName]; found {
// 			a := v.([]string)
// 			if a != nil {
// 				sort.Strings(a)
// 				for _, s := range a {
// 					fmt.Fprint(b, strings.TrimSpace(s))
// 				}
// 			}
// 		}
// 	}

// 	m := v.([]interface{})
// 	for _, c := range m {
// 		writeInt(c.(map[string]interface{}), "index")
// 		writeString(c.(map[string]interface{}), "rule_set")
// 		writeStringArray(c.(map[string]interface{}), "matching_severities")
// 	}

// 	s := b.String()
// 	return hashcode.String(s)
// }

type conditionSorter struct {
	conditions []interface{}
}

func (s *conditionSorter) Len() int {
	return len(s.conditions)
}
func (s *conditionSorter) Swap(i, j int) {
	s.conditions[i], s.conditions[j] = s.conditions[j], s.conditions[i]
}
func (s *conditionSorter) Less(i, j int) bool {
	m := s.conditions[i].(map[string]interface{})
	n := s.conditions[j].(map[string]interface{})
	return m["index"].(int) < n["index"].(int)
}

// ParseConfig reads Terraform config data and stores the information into a
// Circonus RuleSetGroup object.  ParseConfig, ruleSetGroupRead(), and ruleSetGroupChecksum
// must be kept in sync.
func (rsg *circonusRuleSetGroup) ParseConfig(d *schema.ResourceData) error {

	if v, found := d.GetOk("name"); found {
		rsg.Name = v.(string)
	}

	if v, found := d.GetOk("notify"); found {
		y := v.(*schema.Set)
		x := y.List()
		m := x[0].(map[string]interface{})
		for i := 1; i <= 5; i++ {
			s := fmt.Sprintf("sev%d", i)
			sevList := m[s].([]interface{})
			for _, cg := range sevList {
				rsg.ContactGroups[uint8(i)] = append(rsg.ContactGroups[uint8(i)], cg.(string))
			}
		}
	}

	if v, found := d.GetOk("formula"); found {
		y := v.(*schema.Set)
		x := y.List()
		rsg.Formulas = make([]api.RuleSetGroupFormula, 0, len(x))
		for _, f := range x {
			m := f.(map[string]interface{})
			rsgf := api.RuleSetGroupFormula{}
			rsgf.Expression = m["expression"].(string)
			rsgf.RaiseSeverity = uint(m["raise_severity"].(int))
			rsgf.Wait = uint(m["wait"].(int))
			rsg.Formulas = append(rsg.Formulas, rsgf)
		}
	}
	if v, found := d.GetOk("condition"); found {
		y := v.(*schema.Set)
		x := y.List()
		cs := &conditionSorter{
			conditions: x,
		}
		sort.Sort(cs)

		rsg.RuleSetConditions = make([]api.RuleSetGroupCondition, 0, len(x))
		for _, m := range x {
			c := m.(map[string]interface{})
			cond := api.RuleSetGroupCondition{}
			sevs := c["matching_severities"].([]interface{})
			cond.MatchingSeverities = make([]string, 0)
			for _, s := range sevs {
				cond.MatchingSeverities = append(cond.MatchingSeverities, s.(string))
			}
			cond.RuleSetCID = c["rule_set"].(string)
			rsg.RuleSetConditions = append(rsg.RuleSetConditions, cond)
		}
	}

	if v, found := d.GetOk("tags"); found {
		rsg.Tags = derefStringList(flattenSet(v.(*schema.Set)))
	}

	log.Printf("RuleSetGroup: %v\n", rsg)
	if err := rsg.Validate(); err != nil {
		return err
	}

	return nil
}

func (rsg *circonusRuleSetGroup) Create(ctxt *providerContext) error {
	crs, err := ctxt.client.CreateRuleSetGroup(&rsg.RuleSetGroup)
	if err != nil {
		return err
	}

	rsg.CID = crs.CID

	return nil
}

func (rsg *circonusRuleSetGroup) Update(ctxt *providerContext) error {
	_, err := ctxt.client.UpdateRuleSetGroup(&rsg.RuleSetGroup)
	if err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to update rule set group %s: {{err}}", rsg.CID), err)
	}

	return nil
}

func (rsg *circonusRuleSetGroup) Validate() error {

	log.Printf("RuleSetGroup: %v\n", rsg)

	return nil
}
