package circonus

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	// circonus_check.json.* resource attribute names
	checkPromTextPortAttr = "port"
	checkPromTextURLAttr  = "url"
)

var checkPromTextDescriptions = attrDescrs{
	checkPromTextPortAttr: "Specifies the port on which the prometheus metrics can be scraped",
	checkPromTextURLAttr:  "The URL to use as the target of the check",
}

var schemaCheckPromText = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      checkPromTextConfigChecksum,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkPromTextDescriptions, map[schemaAttr]*schema.Schema{
			checkPromTextPortAttr: {
				Type:     schema.TypeInt,
				Default:  443,
				Optional: true,
				ValidateFunc: validateFuncs(
					validateIntMin(checkPromTextPortAttr, 0),
					validateIntMax(checkPromTextPortAttr, 65535),
				),
			},
			checkPromTextURLAttr: {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateFuncs(
					validateHTTPURL(checkPromTextURLAttr, urlIsAbs),
				),
			},
		}),
	},
}

// checkAPIToStateJSON reads the Config data out of circonusCheck.CheckBundle into
// the statefile.
func checkAPIToStatePromText(c *circonusCheck, d *schema.ResourceData) error {
	ptConfig := make(map[string]interface{}, len(c.Config))

	// swamp is a sanity check: it must be empty by the time this method returns
	swamp := make(map[config.Key]string, len(c.Config))
	for k, s := range c.Config {
		swamp[k] = s
	}

	saveStringConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if s, ok := c.Config[apiKey]; ok && s != "" {
			ptConfig[string(attrName)] = s
		}

		delete(swamp, apiKey)
	}

	saveIntConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if s, ok := c.Config[apiKey]; ok && s != "0" {
			i, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Printf("[ERROR]: Unable to convert %s to an integer: %v", apiKey, err)
				return
			}
			ptConfig[string(attrName)] = int(i)
		}

		delete(swamp, apiKey)
	}

	saveIntConfigToState(config.Port, checkPromTextPortAttr)
	saveStringConfigToState(config.URL, checkPromTextURLAttr)

	whitelistedConfigKeys := map[config.Key]struct{}{
		config.ReverseSecretKey: {},
		config.SubmissionURL:    {},
	}

	for k := range swamp {
		if _, ok := whitelistedConfigKeys[k]; ok {
			delete(c.Config, k)
		}

		if _, ok := whitelistedConfigKeys[k]; !ok {
			return fmt.Errorf("PROVIDER BUG: API Config not empty: %#v", swamp)
		}
	}

	if err := d.Set(checkPromTextAttr, schema.NewSet(checkPromTextConfigChecksum, []interface{}{ptConfig})); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store check %q attribute: {{err}}", checkPromTextAttr), err)
	}

	return nil
}

// checkJSONConfigChecksum creates a stable hash of the normalized values found
// in a user's Terraform config.
func checkPromTextConfigChecksum(v interface{}) int {
	m := v.(map[string]interface{})
	b := &bytes.Buffer{}
	b.Grow(defaultHashBufSize)

	writeInt := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok && v.(int) != 0 {
			fmt.Fprintf(b, "%x", v.(int))
		}
	}

	writeString := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok && v.(string) != "" {
			fmt.Fprint(b, strings.TrimSpace(v.(string)))
		}
	}

	// Order writes to the buffer using lexically sorted list for easy visual
	// reconciliation with other lists.
	writeInt(checkPromTextPortAttr)
	writeString(checkPromTextURLAttr)

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPIPromText(c *circonusCheck, l interfaceList) error {
	c.Type = string(apiCheckTypePromText)

	// Iterate over all `promtext` attributes, even though we have a max of 1 in the
	// schema.
	for _, mapRaw := range l {
		ptConfig := newInterfaceMap(mapRaw)

		if v, found := ptConfig[checkPromTextPortAttr]; found {
			i := v.(int)
			if i != 0 {
				c.Config[config.Port] = fmt.Sprintf("%d", i)
			}
		}

		if v, found := ptConfig[checkPromTextURLAttr]; found {
			c.Config[config.URL] = v.(string)

			u, _ := url.Parse(v.(string))
			hostInfo := strings.SplitN(u.Host, ":", 2)
			if len(c.Target) == 0 {
				c.Target = hostInfo[0]
			}

			if len(hostInfo) > 1 && c.Config[config.Port] == "" {
				c.Config[config.Port] = hostInfo[1]
			}
		}

	}

	return nil
}
