package circonus

// The `circonus_metric` type is a synthetic, top-level resource that doesn't
// actually exist within Circonus.  The `circonus_check` resource uses
// `circonus_metric` as input to its `metric` attribute.  The `circonus_check`
// resource can, if configured, override various parameters in the
// `circonus_metric` resource if no value was set.

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// circonus_metric.* resource attribute names.
	metricActiveAttr = "active"
	metricIDAttr     = "id"
	metricNameAttr   = "name"
	metricTypeAttr   = "type"

	// CheckBundle.Metric.Status can be one of these values.
	metricStatusActive    = "active"
	metricStatusAvailable = "available"
)

var metricDescriptions = attrDescrs{
	metricActiveAttr: "Enables or disables the metric",
	metricNameAttr:   "Name of the metric",
	metricTypeAttr:   "Type of metric (e.g. numeric, histogram, text)",
}

func resourceMetric() *schema.Resource {
	return &schema.Resource{
		Create: metricCreate,
		Read:   metricRead,
		Update: metricUpdate,
		Delete: metricDelete,
		Exists: metricExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: convertToHelperSchema(metricDescriptions, map[schemaAttr]*schema.Schema{
			metricActiveAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			metricNameAttr: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp(metricNameAttr, `[\S]+`),
			},
			metricTypeAttr: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringIn(metricTypeAttr, validMetricTypes),
			},
		}),
	}
}

func metricCreate(d *schema.ResourceData, meta interface{}) error {
	m := newMetric()

	id := d.Id()
	if id == "" {
		var err error
		id, err = newMetricID()
		if err != nil {
			return fmt.Errorf("metric ID creation failed: %w", err)
		}
	}

	if err := m.ParseConfig(id, d); err != nil {
		return fmt.Errorf("error parsing metric schema during create: %w", err)
	}

	if err := m.Create(d); err != nil {
		return fmt.Errorf("error creating metric: %w", err)
	}

	return metricRead(d, meta)
}

func metricRead(d *schema.ResourceData, meta interface{}) error {
	m := newMetric()

	if err := m.ParseConfig(d.Id(), d); err != nil {
		return fmt.Errorf("error parsing metric schema during read: %w", err)
	}

	if err := m.SaveState(d); err != nil {
		return fmt.Errorf("error saving metric during read: %w", err)
	}

	return nil
}

func metricUpdate(d *schema.ResourceData, meta interface{}) error {
	m := newMetric()

	if err := m.ParseConfig(d.Id(), d); err != nil {
		return fmt.Errorf("error parsing metric schema during update: %w", err)
	}

	if err := m.Update(d); err != nil {
		return fmt.Errorf("error updating metric: %w", err)
	}

	return nil
}

func metricDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")

	return nil
}

func metricExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	if id := d.Id(); id != "" {
		return true, nil
	}

	return false, nil
}
