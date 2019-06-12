package circonus

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/circonus-labs/circonus-gometrics/api/config"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	// circonus_check.jmx.* resource attribute names
	checkJMXMBeanDomainsAttr    = "mbean_domains"
	checkJMXMBeanPropertiesAttr = "mbean_properties"
	checkJMXPasswordAttr        = "password"
	checkJMXPortAttr            = "port"
	checkJMXHostAttr            = "host"
	checkJMXURIAttr             = "uri"
	checkJMXUsernameAttr        = "username"
)

var checkJMXDescriptions = attrDescrs{
	checkJMXMBeanDomainsAttr:    "The space separated list of domains to filter to",
	checkJMXMBeanPropertiesAttr: "The space separated list of properties to filter to",
	checkJMXPasswordAttr:        "JMX password",
	checkJMXHostAttr:            "JMX host",
	checkJMXPortAttr:            "JMX port",
	checkJMXURIAttr:             "JMX uri, defaults to '/jmxrmi'",
	checkJMXUsernameAttr:        "JMX username",
}

var schemaCheckJMX = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      hashCheckJMX,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkJMXDescriptions, map[schemaAttr]*schema.Schema{
			checkJMXMBeanDomainsAttr: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			checkJMXMBeanPropertiesAttr: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"index": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			checkJMXPasswordAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkJMXPasswordAttr, `.+`),
			},
			checkJMXURIAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "/jmxrmi",
				ValidateFunc: validateRegexp(checkJMXURIAttr, `.+`),
			},
			checkJMXUsernameAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkJMXUsernameAttr, `.+`),
			},
			checkJMXHostAttr: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp(checkJMXHostAttr, `.+`),
			},
			checkJMXPortAttr: {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validateFuncs(
					validateIntMin(checkJMXPortAttr, 0),
					validateIntMax(checkJMXPortAttr, 65535),
				),
			},
		}),
	},
}

// checkAPIToStateJMX reads the Config data out of circonusCheck.CheckBundle into the
// statefile.
func checkAPIToStateJMX(c *circonusCheck, d *schema.ResourceData) error {
	jmxConfig := make(map[string]interface{}, len(c.Config))

	// swamp is a sanity check: it must be empty by the time this method returns
	swamp := make(map[config.Key]string, len(c.Config))
	for k, v := range c.Config {
		swamp[k] = v
	}

	saveIntConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if v, ok := c.Config[apiKey]; ok {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				log.Printf("[ERROR]: Unable to convert %s to an integer: %v", apiKey, err)
				return
			}
			jmxConfig[string(attrName)] = int(i)
		}

		delete(swamp, apiKey)
	}

	saveStringConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if v, ok := c.Config[apiKey]; ok {
			jmxConfig[string(attrName)] = v
		}

		delete(swamp, apiKey)
	}

	saveIntConfigToState(config.Port, checkJMXPortAttr)
	saveStringConfigToState(config.Username, checkJMXUsernameAttr)
	saveStringConfigToState(config.Password, checkJMXPasswordAttr)
	saveStringConfigToState(config.URI, checkJMXURIAttr)
	jmxConfig[string(checkJMXHostAttr)] = c.Target

	l := make([]interface{}, 3)
	// deal with config.MBeanDomains into a list
	if v, ok := c.Config[config.MbeanDomains]; ok {
		log.Printf("Domains: %s", v)
		ll := strings.Split(v, " ")
		for _, i := range ll {
			log.Printf("piece: %s", i)
			l = append(l, i)
		}

		jmxConfig[string(checkJMXMBeanDomainsAttr)] = l
		delete(swamp, checkJMXMBeanDomainsAttr)
	}

	// deal with config.MBeanProperties
	jmxConfig[string(checkJMXMBeanPropertiesAttr)] = make([]interface{}, 1)
	for k, v := range c.Config {
		key := string(k)
		if strings.HasPrefix(key, "mbean_properties_") {
			beanProps := make(map[string]interface{}, 3)
			l := strings.Split(string(v), ",")
			for _, s := range l {
				t := strings.Split(s, "=")
				beanProps[t[0]] = t[1]
			}
			beanProps["index"] = strings.Split(string(k), "_")[2]
			jmxConfig[string(checkJMXMBeanPropertiesAttr)] = append(jmxConfig[string(checkJMXMBeanPropertiesAttr)].([]interface{}), beanProps)
			delete(swamp, k)
		}
	}

	whitelistedConfigKeys := map[config.Key]struct{}{
		config.ReverseSecretKey: {},
		config.SubmissionURL:    {},
	}

	for k := range swamp {
		if _, ok := whitelistedConfigKeys[k]; ok {
			delete(c.Config, k)
		}

		if _, ok := whitelistedConfigKeys[k]; !ok {
			log.Printf("[ERROR]: PROVIDER BUG: API Config not empty: %#v", swamp)
		}
	}

	if err := d.Set(checkJMXAttr, schema.NewSet(hashCheckJMX, []interface{}{jmxConfig})); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store check %q attribute: {{err}}", checkJMXAttr), err)
	}

	return nil
}

// hashCheckJMX creates a stable hash of the normalized values
func hashCheckJMX(v interface{}) int {
	m := v.(map[string]interface{})
	b := &bytes.Buffer{}
	b.Grow(defaultHashBufSize)

	writeInt := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok {
			fmt.Fprintf(b, "%x", v.(int))
		}
	}

	writeString := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok && v.(string) != "" {
			fmt.Fprint(b, strings.TrimSpace(v.(string)))
		}
	}

	writeString(checkJMXPasswordAttr)
	writeString(checkJMXUsernameAttr)
	writeString(checkJMXURIAttr)
	writeString(checkJMXHostAttr)
	writeInt(checkJMXPortAttr)

	list := m[string(checkJMXMBeanDomainsAttr)].([]interface{})
	for _, s := range list {
		if s != nil {
			fmt.Fprint(b, strings.TrimSpace(s.(string)))
		}
	}

	x := m[string(checkJMXMBeanPropertiesAttr)].([]interface{})
	sort.Slice(x, func(i, j int) bool {
		if x[i] != nil && x[j] != nil {
			y := x[i].(map[string]interface{})
			z := x[j].(map[string]interface{})
			return y["index"].(string) < z["index"].(string)
		}
		return true
	})

	for _, s := range x {
		if s != nil {
			t := s.(map[string]interface{})
			fmt.Fprintf(b, "%s%s%s", strings.TrimSpace(t["index"].(string)), strings.TrimSpace(t["name"].(string)), strings.TrimSpace(t["type"].(string)))
		}
	}

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPIJMX(c *circonusCheck, l interfaceList) error {
	c.Type = string(apiCheckTypeJMX)

	// Iterate over all `tcp` attributes, even though we have a max of 1 in the
	// schema.
	for _, mapRaw := range l {
		jmxConfig := newInterfaceMap(mapRaw)

		if v, found := jmxConfig[checkJMXPasswordAttr]; found {
			c.Config[config.Password] = v.(string)
		}

		if v, found := jmxConfig[checkJMXUsernameAttr]; found {
			c.Config[config.Username] = v.(string)
		}

		if v, found := jmxConfig[checkJMXURIAttr]; found {
			c.Config[config.URI] = v.(string)
		}
		if v, found := jmxConfig[checkJMXHostAttr]; found {
			c.Config[config.Host] = v.(string)
		}

		if v, found := jmxConfig[checkJMXPortAttr]; found {
			c.Config[config.Port] = fmt.Sprintf("%d", v.(int))
		}

		if v, found := jmxConfig[checkJMXMBeanDomainsAttr]; found {
			ll := v.([]interface{})
			var strs []string
			for _, x := range ll {
				s := x.(string)
				strs = append(strs, s)
			}
			mbeans := strings.Join(strs, " ")
			c.Config[config.MbeanDomains] = mbeans
		}

		if v, found := jmxConfig[checkJMXMBeanPropertiesAttr]; found {
			m := v.([]interface{})
			for _, ll := range m {
				n := ll.(map[string]interface{})
				c.Config[config.Key(fmt.Sprintf("mbean_properties_%s", n["index"].(string)))] = fmt.Sprintf("name=%s,type=%s", n["name"].(string), n["type"].(string))
			}
		}
	}

	return nil
}
