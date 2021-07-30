package circonus

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// circonus_check.snmp.* resource attribute names.
	checkSNMPAuthPassphrase    = "auth_passphrase"
	checkSNMPAuthProtocol      = "auth_protocol"
	checkSNMPCommunity         = "community"
	checkSNMPContextEngine     = "context_engine"
	checkSNMPContextName       = "context_name"
	checkSNMPOID               = "oid"
	checkSNMPOIDName           = "name"
	checkSNMPOIDPath           = "path"
	checkSNMPOIDType           = "type"
	checkSNMPPort              = "port"
	checkSNMPPrivacyPassphrase = "privacy_passphrase"
	checkSNMPPrivacyProtocol   = "privacy_protocol"
	checkSNMPSecurityEngine    = "security_engine"
	checkSNMPSecurityLevel     = "security_level"
	checkSNMPSecurityName      = "security_name"
	checkSNMPSeparateQueries   = "separate_queries"
	checkSNMPVersion           = "version"
)

var checkSNMPDescriptions = attrDescrs{

	checkSNMPAuthPassphrase:    "The authentication passphrase to use. Only applicaable to SNMP Version 3.",
	checkSNMPAuthProtocol:      "The authentication protocol to use. Only applicaable to SNMP Version 3.",
	checkSNMPCommunity:         "The SNMP community string providing read access.",
	checkSNMPContextEngine:     "The context engine hex value to use. Only applicaable to SNMP Version 3.",
	checkSNMPContextName:       "The context name to use. Only applicaable to SNMP Version 3.",
	checkSNMPOID:               "Defines a metric to query.",
	checkSNMPPort:              "The UDP port to which SNMP queries will be sent.",
	checkSNMPPrivacyPassphrase: "The privacy passphrase to use. Only applicaable to SNMP Version 3.",
	checkSNMPPrivacyProtocol:   "The privacy protocol to use. Only applicaable to SNMP Version 3.",
	checkSNMPSecurityEngine:    "The security engine hex value to use. Only applicaable to SNMP Version 3.",
	checkSNMPSecurityLevel:     "The security level to use for the SNMP session. Choices are \"authPriv\" (authenticated and encrypted), \"authNoPriv\" (authenticated and unencrypted) and \"noAuthNoPriv\" (unauthenticated and unencrypted). Only applicaable to SNMP Version 3.",
	checkSNMPSecurityName:      "The security name (or user name) to use. Only applicaable to SNMP Version 3.",
	checkSNMPSeparateQueries:   "Whether or not to query each OID separately.",
	checkSNMPVersion:           "The SNMP version used for queries.",
}

var checkSNMPOIDDescriptions = attrDescrs{
	checkSNMPOIDName: "Name of the metric produced by this MIB.",
	checkSNMPOIDPath: "The decimal notation or MIB name of this OID.",
	checkSNMPOIDType: "The metric type of this OID. The value can be either one of the single letter codes in the metric_type_t enum or the following string variants: guess, int32, uint32, int64, uint64, double, string.",
}

var schemaCheckSNMP = &schema.Schema{
	Type:     schema.TypeList,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkSNMPDescriptions, map[schemaAttr]*schema.Schema{
			checkSNMPAuthPassphrase: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSNMPAuthPassphrase, `.+`),
			},
			checkSNMPAuthProtocol: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSNMPAuthProtocol, `(MD5|SHA)`),
				Default:      "MD5",
			},
			checkSNMPCommunity: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp(checkSNMPCommunity, `.+`),
			},
			checkSNMPContextEngine: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSNMPContextEngine, `[0-9a-fA-F]+`),
			},
			checkSNMPContextName: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSNMPContextName, `.+`),
			},
			checkSNMPPort: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  161,
			},
			checkSNMPPrivacyPassphrase: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSNMPPrivacyPassphrase, `.+`),
			},
			checkSNMPPrivacyProtocol: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "DES",
				ValidateFunc: validateRegexp(checkSNMPPrivacyProtocol, `(DES|AES128|AES)`),
			},
			checkSNMPSecurityEngine: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSNMPSecurityEngine, `[0-9a-fA-F]+`),
			},
			checkSNMPSecurityLevel: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSNMPSecurityLevel, `(noAuthNoPriv|authNoPriv|authPriv)`),
				Default:      "authPriv",
			},
			checkSNMPSecurityName: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSNMPSecurityName, `.+`),
			},
			checkSNMPSeparateQueries: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			checkSNMPVersion: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp(checkSNMPVersion, `(1|2c|3)`),
			},
			checkSNMPOID: {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: convertToHelperSchema(checkSNMPOIDDescriptions, map[schemaAttr]*schema.Schema{
						checkSNMPOIDName: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateRegexp(checkSNMPOIDName, `^.+$`),
						},
						checkSNMPOIDPath: {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateRegexp(checkSNMPOIDPath, `^.+$`),
						},
						checkSNMPOIDType: {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateRegexp(checkSNMPOIDType, `^.+$`),
						},
					}),
				},
			},
		}),
	},
}

// checkAPIToStateSNMP reads the Config data out of circonusCheck.CheckBundle into the
// statefile.
func checkAPIToStateSNMP(c *circonusCheck, d *schema.ResourceData) error {
	snmpConfig := make(map[string]interface{}, len(c.Config))

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
			snmpConfig[string(attrName)] = int(i)
		}

		delete(swamp, apiKey)
	}

	saveBoolConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if s, ok := c.Config[apiKey]; ok {
			switch s {
			case "true", "on":
				snmpConfig[string(attrName)] = true
			case "false", "off":
				snmpConfig[string(attrName)] = false
			default:
				log.Printf("PROVIDER BUG: unsupported value %q returned in key %q", s, apiKey)
			}
		}

		delete(swamp, apiKey)
	}

	saveStringConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if v, ok := c.Config[apiKey]; ok {
			snmpConfig[string(attrName)] = v
		}

		delete(swamp, apiKey)
	}

	saveStringConfigToState(config.AuthPassphrase, checkSNMPAuthPassphrase)
	saveStringConfigToState(config.AuthProtocol, checkSNMPAuthProtocol)
	saveStringConfigToState(config.Community, checkSNMPCommunity)
	saveStringConfigToState(config.ContextEngine, checkSNMPContextEngine)
	saveStringConfigToState(config.ContextName, checkSNMPContextName)
	saveIntConfigToState(config.Port, checkSNMPPort)
	saveStringConfigToState(config.PrivacyPassphrase, checkSNMPPrivacyPassphrase)
	saveStringConfigToState(config.PrivacyProtocol, checkSNMPPrivacyProtocol)
	saveStringConfigToState(config.SecurityEngine, checkSNMPSecurityEngine)
	saveStringConfigToState(config.SecurityLevel, checkSNMPSecurityLevel)
	saveStringConfigToState(config.SecurityName, checkSNMPSecurityName)
	saveBoolConfigToState(config.SeparateQueries, checkSNMPSeparateQueries)
	saveStringConfigToState(config.Version, checkSNMPVersion)

	// count the number of oids in the config so we can make our list
	oidCount := 0
	for k := range c.Config {
		key := string(k)
		if strings.HasPrefix(key, string(config.OIDPrefix)) {
			oidCount++
		}
	}
	oidList := make([]interface{}, 0)
	for k, v := range c.Config {
		key := string(k)
		if strings.HasPrefix(key, string(config.OIDPrefix)) {
			oidProps := make(map[string]interface{})
			name := key[4:]
			oidProps[string(checkSNMPOIDName)] = name
			oidProps[string(checkSNMPOIDPath)] = v

			t := string(config.TypePrefix) + name
			if tv, ok := c.Config[config.Key(t)]; ok {
				oidProps[string(checkSNMPOIDType)] = tv
				delete(swamp, config.Key(t))
			}
			delete(swamp, k)
			oidList = append(oidList, oidProps)
		}
	}

	sort.Slice(oidList, func(i, j int) bool {
		if oidList[i] != nil && oidList[j] != nil {
			y := oidList[i].(map[string]interface{})
			z := oidList[j].(map[string]interface{})
			return y[string(checkSNMPOIDName)].(string) < z[string(checkSNMPOIDName)].(string)
		}
		return true
	})
	snmpConfig[string(checkSNMPOID)] = oidList

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

	if err := d.Set(checkSNMPAttr, []interface{}{snmpConfig}); err != nil {
		return fmt.Errorf("Unable to store check %q attribute: %w", checkSNMPAttr, err)
	}

	return nil
}

func checkConfigToAPISNMP(c *circonusCheck, l interfaceList) error { //nolint:unparam
	c.Type = string(apiCheckTypeSNMP)

	// Iterate over all `snmp` attributes, even though we have a max of 1 in the
	// schema.
	for _, mapRaw := range l {
		snmpConfig := newInterfaceMap(mapRaw)

		if v, found := snmpConfig[checkSNMPAuthPassphrase]; found {
			c.Config[config.AuthPassphrase] = v.(string)
		}

		if v, found := snmpConfig[checkSNMPAuthProtocol]; found {
			c.Config[config.AuthProtocol] = v.(string)
		}

		if v, found := snmpConfig[checkSNMPCommunity]; found {
			c.Config[config.Community] = v.(string)
		}

		if v, found := snmpConfig[checkSNMPContextEngine]; found {
			c.Config[config.ContextEngine] = v.(string)
		}
		if v, found := snmpConfig[checkSNMPContextName]; found {
			c.Config[config.ContextName] = v.(string)
		}
		if v, found := snmpConfig[checkSNMPPort]; found {
			c.Config[config.Port] = fmt.Sprintf("%d", v.(int))
		}
		if v, found := snmpConfig[checkSNMPPrivacyPassphrase]; found {
			c.Config[config.PrivacyPassphrase] = v.(string)
		}
		if v, found := snmpConfig[checkSNMPPrivacyProtocol]; found {
			c.Config[config.PrivacyProtocol] = v.(string)
		}
		if v, found := snmpConfig[checkSNMPSecurityEngine]; found {
			c.Config[config.SecurityEngine] = v.(string)
		}
		if v, found := snmpConfig[checkSNMPSecurityLevel]; found {
			c.Config[config.SecurityLevel] = v.(string)
		}
		if v, found := snmpConfig[checkSNMPSecurityName]; found {
			c.Config[config.SecurityName] = v.(string)
		}
		if v, found := snmpConfig[checkSNMPSeparateQueries]; found {
			b := v.(bool)
			if b {
				c.Config[config.SeparateQueries] = fmt.Sprintf("%t", b)
			}
		}
		if v, found := snmpConfig[checkSNMPVersion]; found {
			c.Config[config.Version] = v.(string)
		}

		if v, found := snmpConfig[checkSNMPOID]; found {
			m := v.([]interface{})
			for _, ll := range m {
				if ll == nil {
					continue
				}
				n := ll.(map[string]interface{})
				c.Config[config.Key(fmt.Sprintf("oid_%s", n[string(checkSNMPOIDName)].(string)))] = n[string(checkSNMPOIDPath)].(string)
				c.Config[config.Key(fmt.Sprintf("type_%s", n[string(checkSNMPOIDName)].(string)))] = n[string(checkSNMPOIDType)].(string)
			}
		}
	}
	return nil
}
