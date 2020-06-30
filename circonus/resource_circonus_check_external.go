package circonus

import (
	"fmt"
	"strings"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// circonus_check.http.* resource attribute names
	checkCommandAttr       = "command"
	checkOutputExtractAttr = "output_extract"
	checkArg1Attr          = "arg1"
	checkArg2Attr          = "arg2"
	checkArg3Attr          = "arg3"
	checkArg4Attr          = "arg4"
	checkArg5Attr          = "arg5"
	checkArg6Attr          = "arg6"
	checkArg7Attr          = "arg7"
	checkArg8Attr          = "arg8"
	checkArg9Attr          = "arg9"
	checkArg10Attr         = "arg10"
	checkExternalEnvAttr   = "env"
)

var checkExternalDescriptions = attrDescrs{
	checkCommandAttr:       "The full path to the command to run",
	checkOutputExtractAttr: "The output extraction method: JSON or NAGIOS, otherwise treated as regexp",
	checkArg1Attr:          "The 1st argument to the command",
	checkArg2Attr:          "The 2nd argument to the command",
	checkArg3Attr:          "The 3rd argument to the command",
	checkArg4Attr:          "The 4th argument to the command",
	checkArg5Attr:          "The 5th argument to the command",
	checkArg6Attr:          "The 6th argument to the command",
	checkArg7Attr:          "The 7th argument to the command",
	checkArg8Attr:          "The 8th argument to the command",
	checkArg9Attr:          "The 9th argument to the command",
	checkArg10Attr:         "The 10th argument to the command",
	checkExternalEnvAttr:   "The map of environment vars",
}

var schemaCheckExternal = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkExternalDescriptions, map[schemaAttr]*schema.Schema{
			checkOutputExtractAttr: {
				Type:     schema.TypeString,
				Required: true,
				// anything other than JSON or NAGIOS treated as regexp
				// ValidateFunc: validateRegexp(checkOutputExtractAttr, `^(?:JSON|NAGIOS)$`),
			},
			checkCommandAttr: {
				Type:     schema.TypeString,
				Required: true,
			},
			checkArg1Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg2Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg3Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg4Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg5Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg6Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg7Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg8Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg9Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkArg10Attr: {
				Type:     schema.TypeString,
				Optional: true,
			},
			checkExternalEnvAttr: {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     schema.TypeString,
			},
		}),
	},
}

// checkAPIToStateExternal reads the Config data out of circonusCheck.CheckBundle into the
// statefile.
func checkAPIToStateExternal(c *circonusCheck, d *schema.ResourceData) error {
	externalConfig := make(map[string]interface{}, len(c.Config))
	envs := make(map[string]interface{})

	// swamp is a sanity check: it must be empty by the time this method returns
	swamp := make(map[config.Key]string, len(c.Config))
	for k, v := range c.Config {
		swamp[k] = v
	}

	saveStringConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if v, ok := c.Config[apiKey]; ok {
			externalConfig[string(attrName)] = v
		}

		delete(swamp, apiKey)
	}

	saveStringConfigToState("command", checkCommandAttr)
	saveStringConfigToState("output_extract", checkOutputExtractAttr)
	saveStringConfigToState("arg1", checkArg1Attr)
	saveStringConfigToState("arg2", checkArg2Attr)
	saveStringConfigToState("arg3", checkArg3Attr)
	saveStringConfigToState("arg4", checkArg4Attr)
	saveStringConfigToState("arg5", checkArg5Attr)
	saveStringConfigToState("arg6", checkArg6Attr)
	saveStringConfigToState("arg7", checkArg7Attr)
	saveStringConfigToState("arg8", checkArg8Attr)
	saveStringConfigToState("arg9", checkArg9Attr)
	saveStringConfigToState("arg10", checkArg10Attr)
	saveStringConfigToState("arg10", checkArg10Attr)

	// env vars
	for k, v := range c.Config {
		if len(k) <= 4 {
			continue
		}

		// Handle all of the prefix variable headers, like `header_`
		if strings.Compare(string(k[:4]), "env_") == 0 {
			key := k[4:]
			envs[string(key)] = v
			delete(swamp, k)
		}
	}

	externalConfig[string(checkExternalEnvAttr)] = envs

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

	if err := d.Set(checkExternalAttr, []interface{}{externalConfig}); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store check %q attribute: {{err}}", checkExternalAttr), err)
	}

	return nil
}

func checkConfigToAPIExternal(c *circonusCheck, l interfaceList) error {
	c.Type = string(apiCheckTypeExternal)

	// Iterate over all `http` attributes, even though we have a max of 1 in the
	// schema.
	for _, mapRaw := range l {
		externalConfig := newInterfaceMap(mapRaw)

		if v, found := externalConfig[checkCommandAttr]; found {
			c.Config["command"] = v.(string)
		}

		if v, found := externalConfig[checkOutputExtractAttr]; found {
			c.Config["output_extract"] = v.(string)
		}

		if v, found := externalConfig[checkArg1Attr]; found {
			c.Config["arg1"] = v.(string)
		}

		if v, found := externalConfig[checkArg2Attr]; found {
			c.Config["arg2"] = v.(string)
		}

		if v, found := externalConfig[checkArg3Attr]; found {
			c.Config["arg3"] = v.(string)
		}

		if v, found := externalConfig[checkArg4Attr]; found {
			c.Config["arg4"] = v.(string)
		}

		if v, found := externalConfig[checkArg5Attr]; found {
			c.Config["arg5"] = v.(string)
		}

		if v, found := externalConfig[checkArg6Attr]; found {
			c.Config["arg6"] = v.(string)
		}

		if v, found := externalConfig[checkArg7Attr]; found {
			c.Config["arg7"] = v.(string)
		}

		if v, found := externalConfig[checkArg8Attr]; found {
			c.Config["arg8"] = v.(string)
		}

		if v, found := externalConfig[checkArg9Attr]; found {
			c.Config["arg9"] = v.(string)
		}

		if v, found := externalConfig[checkArg10Attr]; found {
			c.Config["arg10"] = v.(string)
		}

		for k, v := range externalConfig.CollectMap(checkExternalEnvAttr) {
			h := config.Key("env_") + config.Key(k)
			c.Config[h] = v
		}
	}

	return nil
}
