package circonus

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// circonus_check.redis.* resource attribute names
	checkRedisCommandAttr  = "command"
	checkRedisDbIndexAttr  = "db_index"
	checkRedisPasswordAttr = "password"
	checkRedisPortAttr     = "port"
)

var checkRedisDescriptions = attrDescrs{
	checkRedisCommandAttr:  "The redis command to run to gather stats, default: INFO.",
	checkRedisDbIndexAttr:  "The database index to query, defaults to zero",
	checkRedisPasswordAttr: "The pass required to run the command.",
	checkRedisPortAttr:     "Specifies the port on which the Redis instance can be reached.",
}

var schemaCheckRedis = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      hashCheckRedis,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkRedisDescriptions, map[schemaAttr]*schema.Schema{
			checkRedisCommandAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "INFO",
				ValidateFunc: validateRegexp(checkRedisCommandAttr, `.+`),
			},
			checkRedisDbIndexAttr: {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntMin(checkRedisDbIndexAttr, 0),
			},
			checkRedisPasswordAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validateRegexp(checkRedisPasswordAttr, `.+`),
			},
			checkTCPPortAttr: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  6379,
				ValidateFunc: validateFuncs(
					validateIntMin(checkTCPPortAttr, 0),
					validateIntMax(checkTCPPortAttr, 65535),
				),
			},
		}),
	},
}

func checkAPIToStateRedis(c *circonusCheck, d *schema.ResourceData) error {
	redisConfig := make(map[string]interface{}, len(c.Config))

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
			redisConfig[string(attrName)] = int(i)
		}

		delete(swamp, apiKey)
	}

	saveStringConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if v, ok := c.Config[apiKey]; ok {
			redisConfig[string(attrName)] = v
		}

		delete(swamp, apiKey)
	}

	saveStringConfigToState(config.Command, checkRedisCommandAttr)
	saveIntConfigToState(config.DBIndex, checkRedisDbIndexAttr)
	saveStringConfigToState(config.Password, checkRedisPasswordAttr)
	saveIntConfigToState(config.Port, checkRedisPortAttr)

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

	if err := d.Set(checkRedisAttr, schema.NewSet(hashCheckRedis, []interface{}{redisConfig})); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store check %q attribute: {{err}}", checkRedisAttr), err)
	}

	return nil
}

func hashCheckRedis(v interface{}) int {
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

	// Order writes to the buffer using lexically sorted list for easy visual
	// reconciliation with other lists.
	writeString(checkRedisCommandAttr)
	writeInt(checkRedisDbIndexAttr)
	writeString(checkRedisPasswordAttr)
	writeInt(checkRedisPortAttr)

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPIRedis(c *circonusCheck, l interfaceList) error {
	c.Type = string(apiCheckTypeRedis)

	// Iterate over all `tcp` attributes, even though we have a max of 1 in the
	// schema.
	for _, mapRaw := range l {
		redisConfig := newInterfaceMap(mapRaw)

		if v, found := redisConfig[checkRedisCommandAttr]; found {
			c.Config[config.Command] = v.(string)
		}

		if v, found := redisConfig[checkRedisDbIndexAttr]; found {
			c.Config[config.DBIndex] = fmt.Sprintf("%d", v.(int))
		}

		if v, found := redisConfig[checkRedisPasswordAttr]; found {
			c.Config[config.Password] = v.(string)
		}

		if v, found := redisConfig[checkRedisPortAttr]; found {
			c.Config[config.Port] = fmt.Sprintf("%d", v.(int))
		}

	}

	return nil
}
