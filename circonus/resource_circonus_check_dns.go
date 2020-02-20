package circonus

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// circonus_check.dns.* resource attribute names
	checkDNSCTypeAttr      = "ctype"
	checkDNSNameserverAttr = "nameserver"
	checkDNSQueryAttr      = "query"
	checkDNSRTypeAttr      = "rtype"
)

var checkDNSDescriptions = attrDescrs{
	checkDNSCTypeAttr:      "The DNS class of the query. IN: Internet, CH: Chaos, HS: Hesoid.",
	checkDNSNameserverAttr: "The domain name server to query. If the name of the check is in-addr.arpa, the system default nameserver is used. Otherwise, the nameserver is the %[target] of the the check.",
	checkDNSQueryAttr:      "The query to send. If the name of the check is in-addr.arpa, the reverse IP octet notation of in-addr.arpa syntax is synthesized by default. Otherwise the default query is the name of the check itself.",
	checkDNSRTypeAttr:      "The DNS resource record type of the query. If the name of the check is in-addr.arpa, the default is PTR, otherwise it is A.",
}

var schemaCheckDNS = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      hashCheckDNS,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkDNSDescriptions, map[schemaAttr]*schema.Schema{
			checkDNSCTypeAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IN",
				ValidateFunc: validateStringIn(checkDNSCTypeAttr, validStringValues{"IN", "CH", "HS"}),
			},
			checkDNSNameserverAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "%[target]",
				ValidateFunc: validateRegexp(checkDNSNameserverAttr, ".+"),
			},
			checkDNSQueryAttr: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateRegexp(checkDNSNameserverAttr, ".+"),
			},
			checkDNSRTypeAttr: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "A",
				ValidateFunc: validateStringIn(checkDNSCTypeAttr, validStringValues{
					"A",
					"AAAA",
					"TXT",
					"MX",
					"SOA",
					"CNAME",
					"PTR",
					"NS",
					"MB",
					"MD",
					"MF",
					"MG",
					"MR",
				}),
			},
		}),
	},
}

// checkAPIToStateDNS reads the Config data out of circonusCheck.CheckBundle
// into the statefile.
func checkAPIToStateDNS(c *circonusCheck, d *schema.ResourceData) error {
	dnsConfig := make(map[string]interface{}, len(c.Config))

	if ctype, ok := c.Config[config.CType]; ok {
		dnsConfig[string(checkDNSCTypeAttr)] = ctype
	}

	if ns, ok := c.Config[config.Nameserver]; ok {
		dnsConfig[string(checkDNSNameserverAttr)] = ns
	}

	if q, ok := c.Config[config.Query]; ok {
		dnsConfig[string(checkDNSQueryAttr)] = q
	}

	if rtype, ok := c.Config[config.RType]; ok {
		dnsConfig[string(checkDNSRTypeAttr)] = rtype
	}

	if err := d.Set(checkDNSAttr, schema.NewSet(hashCheckDNS, []interface{}{dnsConfig})); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("Unable to store check %q attribute: {{err}}", checkDNSAttr), err)
	}

	return nil
}

// hashCheckICMPPing creates a stable hash of the normalized values
func hashCheckDNS(v interface{}) int {
	m := v.(map[string]interface{})
	b := &bytes.Buffer{}
	b.Grow(defaultHashBufSize)

	writeString := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok && v.(string) != "" {
			fmt.Fprint(b, strings.TrimSpace(v.(string)))
		}
	}

	writeString(checkDNSCTypeAttr)
	writeString(checkDNSNameserverAttr)
	writeString(checkDNSQueryAttr)
	writeString(checkDNSRTypeAttr)

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPIDNS(c *circonusCheck, l interfaceList) error {
	c.Type = string(apiCheckTypeDNS)

	mapRaw := l[0]
	dnsConfig := newInterfaceMap(mapRaw)

	if v, found := dnsConfig[checkDNSCTypeAttr]; found && v.(string) != "" {
		c.Config[config.CType] = v.(string)
	}

	if v, found := dnsConfig[checkDNSNameserverAttr]; found && v.(string) != "" {
		c.Config[config.Nameserver] = v.(string)
	}

	if v, found := dnsConfig[checkDNSQueryAttr]; found && v.(string) != "" {
		c.Config[config.Query] = v.(string)
	}

	if v, found := dnsConfig[checkDNSRTypeAttr]; found && v.(string) != "" {
		c.Config[config.RType] = v.(string)
	}

	return nil
}
