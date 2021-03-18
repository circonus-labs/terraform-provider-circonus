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
	// circonus_check.ntp.* resource attribute names
	checkNTPPortAttr       = "port"
	checkNTPUseControlAttr = "use_control"
)

var checkNTPDescriptions = attrDescrs{
	checkNTPPortAttr:       "The port to talk to NTP over (default: 123)",
	checkNTPUseControlAttr: "Control protocol means that the agent will request the NTP telemetry of the target regarding its preferred peer, (default: false)",
}

var schemaCheckNTP = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      hashCheckNTP,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkNTPDescriptions, map[schemaAttr]*schema.Schema{
			checkNTPPortAttr: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  123,
			},
			checkNTPUseControlAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		}),
	},
}

// checkAPIToStateNTP reads the Config data out of circonusCheck.CheckBundle
// into the statefile.
func checkAPIToStateNTP(c *circonusCheck, d *schema.ResourceData) error {
	ntpConfig := make(map[string]interface{}, len(c.Config))

	if port, ok := c.Config[config.Port]; ok {
		ntpConfig[string(checkNTPPortAttr)], _ = strconv.Atoi(port)
	}

	if control, ok := c.Config[config.Control]; ok {
		ntpConfig[string(checkNTPUseControlAttr)], _ = strconv.ParseBool(control)
	}

	if err := d.Set(checkNTPAttr, schema.NewSet(hashCheckNTP, []interface{}{ntpConfig})); err != nil {
		return fmt.Errorf("Unable to store check %q attribute: %w", checkNTPAttr, err)
	}

	return nil
}

// hashCheckNTP creates a stable hash of the normalized values
func hashCheckNTP(v interface{}) int {
	m := v.(map[string]interface{})
	b := &bytes.Buffer{}
	b.Grow(defaultHashBufSize)

	writeBool := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok {
			fmt.Fprintf(b, "%t", v.(bool))
		}
	}

	writeInt := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok {
			fmt.Fprintf(b, "%x", v.(int))
		}
	}

	writeInt(checkNTPPortAttr)
	writeBool(checkNTPUseControlAttr)

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPINTP(c *circonusCheck, l interfaceList) error {
	c.Type = string(apiCheckTypeNTP)

	mapRaw := l[0]
	ntpConfig := newInterfaceMap(mapRaw)

	if v, found := ntpConfig[checkNTPPortAttr]; found && v.(int) != 0 {
		c.Config[config.Port] = strconv.Itoa(v.(int))
	}

	if v, found := ntpConfig[checkNTPUseControlAttr]; found {
		c.Config[config.Control] = fmt.Sprintf("%t", v.(bool))
	}

	return nil
}
