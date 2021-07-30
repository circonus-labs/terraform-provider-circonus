package circonus

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/circonus-labs/terraform-provider-circonus/internal/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// circonus_check.memcached.* resource attribute names.
	checkMemcachedPortAttr = "port"
)

var checkMemcachedDescriptions = attrDescrs{
	checkMemcachedPortAttr: `The port the memcached instance is listenening on, default 11211`,
}

var schemaCheckMemcached = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      hashCheckMemcached,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkMemcachedDescriptions, map[schemaAttr]*schema.Schema{
			checkMemcachedPortAttr: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  11211,
				ValidateFunc: validateFuncs(
					validateIntMin(checkMemcachedPortAttr, 1),
					validateIntMax(checkMemcachedPortAttr, 65535),
				),
			},
		}),
	},
}

// checkAPIToStateMemcached reads the Config data out of circonusCheck.CheckBundle
// into the statefile.
func checkAPIToStateMemcached(c *circonusCheck, d *schema.ResourceData) error {
	memcachedConfig := make(map[string]interface{}, len(c.Config))

	port, err := strconv.ParseInt(c.Config[config.Port], 10, 64)
	if err != nil {
		return fmt.Errorf("unable to parse %s: %w", config.Port, err)
	}

	memcachedConfig[string(checkMemcachedPortAttr)] = int(port)

	if err := d.Set(checkMemcachedAttr, schema.NewSet(hashCheckMemcached, []interface{}{memcachedConfig})); err != nil {
		return fmt.Errorf("Unable to store check %q attribute: %w", checkMemcachedAttr, err)
	}

	return nil
}

// hashCheckICMPPing creates a stable hash of the normalized values.
func hashCheckMemcached(v interface{}) int {
	m := v.(map[string]interface{})
	b := &bytes.Buffer{}
	b.Grow(defaultHashBufSize)

	writeInt := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok {
			fmt.Fprintf(b, "%x", v.(int))
		}
	}

	// Order writes to the buffer using lexically sorted list for easy visual
	// reconciliation with other lists.
	writeInt(checkMemcachedPortAttr)

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPIMemcached(c *circonusCheck, l interfaceList) error { //nolint:unparam
	c.Type = string(apiCheckTypeMemcached)

	// Iterate over all `memcached` attributes, even though we have a max of 1 in
	// the schema.
	for _, mapRaw := range l {
		memcachedConfig := newInterfaceMap(mapRaw)

		if v, found := memcachedConfig[checkMemcachedPortAttr]; found {
			c.Config[config.Port] = fmt.Sprintf("%d", v.(int))
		}
	}

	return nil
}
