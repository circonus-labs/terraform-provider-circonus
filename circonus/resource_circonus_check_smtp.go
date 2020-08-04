package circonus

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	checkSMTPEhloAttr               = "ehlo"
	checkSMTPFromAttr               = "from"
	checkSMTPPayloadAttr            = "payload"
	checkSMTPPortAttr               = "port"
	checkSMTPProxyDestAddressAttr   = "proxy_dest_address"
	checkSMTPProxyDestPortAttr      = "proxy_dest_port"
	checkSMTPProxyFamilyAttr        = "proxy_family"
	checkSMTPProxyProtocolAttr      = "proxy_protocol"
	checkSMTPProxySourceAddressAttr = "proxy_source_address"
	checkSMTPProxySourcePortAttr    = "proxy_source_port"
	checkSMTPSaslAuthIDAttr         = "sasl_auth_id"
	checkSMTPSaslAuthenticationAttr = "sasl_authentication"
	checkSMTPSaslPasswordAttr       = "sasl_password"
	checkSMTPSaslUserAttr           = "sasl_user"
	checkSMTPStartTLSAttr           = "starttls"
	checkSMTPToAttr                 = "to"
)

var checkSMTPDescriptions = attrDescrs{
	checkSMTPEhloAttr:               "Specifies the EHLO parameter. (default: noit.local)",
	checkSMTPFromAttr:               "Specifies the envelope sender.",
	checkSMTPPayloadAttr:            "Specifies the payload sent (on the wire). CR LF DOT CR LF is appended automatically. (default: Subject: Testing)",
	checkSMTPPortAttr:               "Specifies the TCP port to connect to. (default: 25)",
	checkSMTPProxyDestAddressAttr:   "The IP (or string) to use as the destination address portion of the PROXY protocol. More on the proxy protocol here: http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt",
	checkSMTPProxyDestPortAttr:      "The port to use as the dest port portion of the PROXY protocol. Defaults to the port setting or 25",
	checkSMTPProxyFamilyAttr:        "The protocol family to send in the PROXY header. (default: TCP4)",
	checkSMTPProxyProtocolAttr:      "Test MTA responses to a PROXY protocol header by setting this to true. (default: false)",
	checkSMTPProxySourceAddressAttr: "The IP (or string) to use as the source address portion of the PROXY protocol. More on the proxy protocol here: http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt",
	checkSMTPProxySourcePortAttr:    "The port to use as the source port portion of the PROXY protocol. Defaults to the actual source port of the connection to the target_ip.",
	checkSMTPSaslAuthIDAttr:         "The SASL Authorization Identity.",
	checkSMTPSaslAuthenticationAttr: "Specifies the type of SASL Authentication to use. (default: off)",
	checkSMTPSaslPasswordAttr:       "The SASL Authentication password.",
	checkSMTPSaslUserAttr:           "The SASL Authentication username.",
	checkSMTPStartTLSAttr:           "Specified if the client should attempt a STARTTLS upgrade. (default: false)",
	checkSMTPToAttr:                 "Specifies the envelope recipient.",
}

var schemaCheckSMTP = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      hashCheckSMTP,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkSMTPDescriptions, map[schemaAttr]*schema.Schema{
			checkSMTPEhloAttr: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "noit.local",
			},
			checkSMTPFromAttr: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			checkSMTPPayloadAttr: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Subect: Testing",
			},
			checkSMTPPortAttr: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  25,
			},
			checkSMTPProxyDestAddressAttr: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			checkSMTPProxyDestPortAttr: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			checkSMTPProxyFamilyAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "TCP4",
				ValidateFunc: validateRegexp(checkSMTPProxyFamilyAttr, `^TCP(4|6)$`),
			},
			checkSMTPProxyProtocolAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			checkSMTPProxySourceAddressAttr: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			checkSMTPProxySourcePortAttr: {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			checkSMTPSaslAuthIDAttr: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			checkSMTPSaslAuthenticationAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "off",
				ValidateFunc: validateRegexp(checkSMTPSaslAuthenticationAttr, `^(off|login|plain)$`),
			},
			checkSMTPSaslPasswordAttr: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			checkSMTPSaslUserAttr: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			checkSMTPStartTLSAttr: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			checkSMTPToAttr: {
				Type:     schema.TypeString,
				Optional: false,
				Default:  "",
			},
		}),
	},
}

// checkAPIToStateSMTP reads the Config data out of circonusCheck.CheckBundle
// into the statefile
func checkAPIToStateSMTP(c *circonusCheck, d *schema.ResourceData) error {
	smtpConfig := make(map[string]interface{}, len(c.Config))

	if ehlo, ok := c.Config[config.EHLO]; ok {
		smtpConfig[string(checkSMTPEhloAttr)] = ehlo
	}

	if from, ok := c.Config[config.From]; ok {
		smtpConfig[string(checkSMTPFromAttr)] = from
	}

	if payload, ok := c.Config[config.Payload]; ok {
		smtpConfig[string(checkSMTPPayloadAttr)] = payload
	}

	if port, ok := c.Config[config.Port]; ok {
		smtpConfig[string(checkSMTPPortAttr)], _ = strconv.Atoi(port)
	}

	if proxyDestAddr, ok := c.Config[config.ProxyDestAddress]; ok {
		smtpConfig[string(checkSMTPProxyDestAddressAttr)] = proxyDestAddr
	}

	if proxyDestPort, ok := c.Config[config.ProxyDestPort]; ok {
		p, _ := strconv.Atoi(proxyDestPort)
		if p > 0 {
			smtpConfig[string(checkSMTPProxyDestPortAttr)] = p
		}
	}

	if proxyFamily, ok := c.Config[config.ProxyFamily]; ok {
		smtpConfig[string(checkSMTPProxyFamilyAttr)] = proxyFamily
	}

	if proxyProto, ok := c.Config[config.ProxyProtocol]; ok {
		smtpConfig[string(checkSMTPProxyProtocolAttr)], _ = strconv.ParseBool(proxyProto)
	}

	if proxySrcAddr, ok := c.Config[config.ProxySourceAddress]; ok {
		smtpConfig[string(checkSMTPProxySourceAddressAttr)] = proxySrcAddr
	}

	if proxySrcPort, ok := c.Config[config.ProxySourcePort]; ok {
		p, _ := strconv.Atoi(proxySrcPort)
		if p > 0 {
			smtpConfig[string(checkSMTPProxySourcePortAttr)] = p
		}
	}

	if saslAuthID, ok := c.Config[config.SASLAuthID]; ok {
		smtpConfig[string(checkSMTPSaslAuthIDAttr)] = saslAuthID
	}

	if saslAuth, ok := c.Config[config.SASLAuthentication]; ok {
		smtpConfig[string(checkSMTPSaslAuthenticationAttr)] = saslAuth
	}

	if saslPassword, ok := c.Config[config.SASLPassword]; ok {
		smtpConfig[string(checkSMTPSaslPasswordAttr)] = saslPassword
	}

	if saslUser, ok := c.Config[config.SASLUser]; ok {
		smtpConfig[string(checkSMTPSaslUserAttr)] = saslUser
	}

	if startTLS, ok := c.Config[config.StartTLS]; ok {
		smtpConfig[string(checkSMTPStartTLSAttr)], _ = strconv.ParseBool(startTLS)
	}

	if to, ok := c.Config[config.To]; ok {
		smtpConfig[string(checkSMTPToAttr)] = to
	}

	if err := d.Set(checkSMTPAttr, schema.NewSet(hashCheckSMTP, []interface{}{smtpConfig})); err != nil {
		return errwrap.Wrapf(fmt.Sprintf("unable to store check %q attribute: {{err}}", checkSMTPAttr), err)
	}

	return nil
}

// hashCheckSMTP creates a stable hash of the normalized values
func hashCheckSMTP(v interface{}) int {
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
			if v.(int) > 0 {
				fmt.Fprintf(b, "%x", v.(int))
			}
		}
	}

	writeString := func(attrName schemaAttr) {
		if v, ok := m[string(attrName)]; ok {
			fmt.Fprintf(b, "%s", v.(string))
		}
	}

	writeString(checkSMTPEhloAttr)
	writeString(checkSMTPFromAttr)
	writeString(checkSMTPPayloadAttr)
	writeInt(checkSMTPPortAttr)
	writeString(checkSMTPProxyDestAddressAttr)
	writeInt(checkSMTPProxyDestPortAttr)
	writeString(checkSMTPProxyFamilyAttr)
	writeBool(checkSMTPProxyProtocolAttr)
	writeString(checkSMTPProxySourceAddressAttr)
	writeInt(checkSMTPProxySourcePortAttr)
	writeString(checkSMTPSaslAuthIDAttr)
	writeString(checkSMTPSaslAuthenticationAttr)
	writeString(checkSMTPSaslPasswordAttr)
	writeString(checkSMTPSaslUserAttr)
	writeBool(checkSMTPStartTLSAttr)
	writeString(checkSMTPToAttr)

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPISMTP(c *circonusCheck, l interfaceList) error {
	c.Type = string(apiCheckTypeSMTP)

	mapRaw := l[0]
	smtpConfig := newInterfaceMap(mapRaw)

	if v, found := smtpConfig[checkSMTPEhloAttr]; found && v.(string) != "" {
		c.Config[config.EHLO] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPFromAttr]; found && v.(string) != "" {
		c.Config[config.From] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPPayloadAttr]; found && v.(string) != "" {
		c.Config[config.Payload] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPPortAttr]; found && v.(int) > 0 {
		c.Config[config.Port] = strconv.Itoa(v.(int))
	}

	if v, found := smtpConfig[checkSMTPProxyDestAddressAttr]; found && v.(string) != "" {
		c.Config[config.ProxyDestAddress] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPProxyDestPortAttr]; found && v.(int) > 0 {
		c.Config[config.ProxyDestPort] = strconv.Itoa(v.(int))
	}

	if v, found := smtpConfig[checkSMTPProxyFamilyAttr]; found && v.(string) != "" {
		c.Config[config.ProxyFamily] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPProxyProtocolAttr]; found {
		c.Config[config.ProxyProtocol] = fmt.Sprintf("%t", v.(bool))
	}

	if v, found := smtpConfig[checkSMTPProxySourceAddressAttr]; found && v.(string) != "" {
		c.Config[config.ProxySourceAddress] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPProxySourcePortAttr]; found && v.(int) > 0 {
		c.Config[config.ProxySourcePort] = strconv.Itoa(v.(int))
	}

	if v, found := smtpConfig[checkSMTPSaslAuthIDAttr]; found && v.(string) != "" {
		c.Config[config.SASLAuthID] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPSaslAuthenticationAttr]; found && v.(string) != "" {
		c.Config[config.SASLAuthentication] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPSaslPasswordAttr]; found && v.(string) != "" {
		c.Config[config.SASLPassword] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPSaslUserAttr]; found && v.(string) != "" {
		c.Config[config.SASLUser] = v.(string)
	}

	if v, found := smtpConfig[checkSMTPStartTLSAttr]; found {
		c.Config[config.StartTLS] = fmt.Sprintf("%t", v.(bool))
	}

	if v, found := smtpConfig[checkSMTPToAttr]; found && v.(string) != "" {
		c.Config[config.To] = v.(string)
	}

	return nil
}
