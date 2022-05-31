package circonus

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/circonus-labs/go-apiclient/config"
	"github.com/circonus-labs/terraform-provider-circonus/internal/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	// circonus_check.ssh2.* resource attribute names.
	checkSSH2PortAttr          = "port"
	checkSSH2MethodKexAttr     = "method_kex"
	checkSSH2MethodHostKeyAttr = "method_hostkey"
	checkSSH2MethodCryptCSAttr = "method_crypt_cs"
	checkSSH2MethodCryptSCAttr = "method_crypt_sc"
	checkSSH2MethodMacCSAttr   = "method_mac_cs"
	checkSSH2MethodMacSCAttr   = "method_mac_sc"
	checkSSH2MethodCompCSAttr  = "method_comp_cs"
	checkSSH2MethodCompSCAttr  = "method_comp_sc"
	checkSSH2MethodLangCSAttr  = "method_lang_cs"
	checkSSH2MethodLangSCAttr  = "method_land_sc"
)

var checkSSH2Descriptions = attrDescrs{
	checkSSH2PortAttr:          "The TCP port on which the remote server's ssh service is running",
	checkSSH2MethodKexAttr:     "The key exchange method to use",
	checkSSH2MethodHostKeyAttr: "The host key algorithm supported",
	checkSSH2MethodCryptCSAttr: "The encryption algorithm used from client to server",
	checkSSH2MethodCryptSCAttr: "The encryption algorithm used from server to client",
	checkSSH2MethodMacCSAttr:   "The message authentication code algorithm used from client to server",
	checkSSH2MethodMacSCAttr:   "The message authentication code algorithm used from server to client",
	checkSSH2MethodCompCSAttr:  "The compress algorithm used from client to server",
	checkSSH2MethodCompSCAttr:  "The compress algorithm used from server to client",
	checkSSH2MethodLangCSAttr:  "The language used from client to server",
	checkSSH2MethodLangSCAttr:  "The language used from server to client",
}

var schemaCheckSSH2 = &schema.Schema{
	Type:     schema.TypeSet,
	Optional: true,
	MaxItems: 1,
	MinItems: 1,
	Set:      checkSSH2ConfigChecksum,
	Elem: &schema.Resource{
		Schema: convertToHelperSchema(checkSSH2Descriptions, map[schemaAttr]*schema.Schema{
			checkSSH2PortAttr: {
				Type:     schema.TypeInt,
				Default:  defaultCheckSSH2Port,
				Optional: true,
				ValidateFunc: validateFuncs(
					validateIntMin(checkSSH2PortAttr, 0),
					validateIntMax(checkSSH2PortAttr, 65535),
				),
			},
			checkSSH2MethodKexAttr: {
				Type:         schema.TypeString,
				Default:      defaultCheckSSH2MethodKex,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodKexAttr, `^diffie-hellman-(?:group1-sha1|group14-sha1|group16-sha512|group18-sha512)$`),
			},
			checkSSH2MethodHostKeyAttr: {
				Type:         schema.TypeString,
				Default:      defaultCheckSSH2MethodHostKey,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodHostKeyAttr, `^(?:ssh-dss|ssh-rsa|ecdsa-sha2-nistp256|ssh-ed25519)$`),
			},
			checkSSH2MethodCryptCSAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodCryptCSAttr, `^(?:chacha20-poly1305@openssh.com|aes256-gcm@openssh.com|aes128-gcm@openssh.com|aes256-ctr|aes192-ctr|aes128-ctr|aes256-cbc|aes192-cbc|aes128-cbc|rijndael128-cbc|rijndael192-cbc|rijndael256-cbc|blowfish-cbc|blowfish-ecb|blowfish-cfb|blowfish-ofb|blowfish-ctr|twofish128-ctr|twofish128-cbc|twofish192-ctr|twofish192-cbc|twofish256-ctr|twofish256-cbc|twofish-cbc|twofish-ecb|twofish-cfb|twofish-ofb|arcfour256|arcfour128|arcfour|cast128-cbc|cast128-ecb|cast128-cfb|cast128-ofb|idea-cbc|idea-ecb|idea-cfb|idea-ofb|3des-cbc|3des-ecb|3des-cfb|3des-ofb|3des-ctr|none)$`),
			},
			checkSSH2MethodCryptSCAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodCryptSCAttr, `^(?:chacha20-poly1305@openssh.com|aes256-gcm@openssh.com|aes128-gcm@openssh.com|aes256-ctr|aes192-ctr|aes128-ctr|aes256-cbc|aes192-cbc|aes128-cbc|rijndael128-cbc|rijndael192-cbc|rijndael256-cbc|blowfish-cbc|blowfish-ecb|blowfish-cfb|blowfish-ofb|blowfish-ctr|twofish128-ctr|twofish128-cbc|twofish192-ctr|twofish192-cbc|twofish256-ctr|twofish256-cbc|twofish-cbc|twofish-ecb|twofish-cfb|twofish-ofb|arcfour256|arcfour128|arcfour|cast128-cbc|cast128-ecb|cast128-cfb|cast128-ofb|idea-cbc|idea-ecb|idea-cfb|idea-ofb|3des-cbc|3des-ecb|3des-cfb|3des-ofb|3des-ctr|none)$`),
			},
			checkSSH2MethodMacCSAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodMacCSAttr, `^(?:hmac-sha2-512-etm@openssh.com|hmac-sha2-256-etm@openssh.com|umac-128-etm@openssh.com|umac-64-etm@openssh.com|hmac-sha2-512|hmac-sha2-256|hmac-sha1|hmac-sha1-96|hmac-md5|hmac-md5-96|hmac-ripemd160|hmac-ripemd160@openssh.com|none)$`),
			},
			checkSSH2MethodMacSCAttr: {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodMacSCAttr, `^(?:hmac-sha2-512-etm@openssh.com|hmac-sha2-256-etm@openssh.com|umac-128-etm@openssh.com|umac-64-etm@openssh.com|hmac-sha2-512|hmac-sha2-256|hmac-sha1|hmac-sha1-96|hmac-md5|hmac-md5-96|hmac-ripemd160|hmac-ripemd160@openssh.com|none)$`),
			},
			checkSSH2MethodCompCSAttr: {
				Type:         schema.TypeString,
				Default:      defaultCheckSSH2MethodCompCS,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodCompCSAttr, `^(?:zlib|none)$`),
			},
			checkSSH2MethodCompSCAttr: {
				Type:         schema.TypeString,
				Default:      defaultCheckSSH2MethodCompSC,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodCompSCAttr, `^(?:zlib|none)$`),
			},
			checkSSH2MethodLangCSAttr: {
				Type:         schema.TypeString,
				Default:      defaultCheckSSH2MethodLangCS,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodLangCSAttr, `^(?:|\w+)$`),
			},
			checkSSH2MethodLangSCAttr: {
				Type:         schema.TypeString,
				Default:      defaultCheckSSH2MethodLangSC,
				Optional:     true,
				ValidateFunc: validateRegexp(checkSSH2MethodLangSCAttr, `^(?:|\w+)$`),
			},
		}),
	},
}

// checkAPIToStateSSH2 reads the Config data out of circonusCheck.CheckBundle into
// the statefile.
func checkAPIToStateSSH2(c *circonusCheck, d *schema.ResourceData) error {
	ssh2Config := make(map[string]interface{}, len(c.Config))

	// swamp is a sanity check: it must be empty by the time this method returns
	swamp := make(map[config.Key]string, len(c.Config))
	for k, s := range c.Config {
		swamp[k] = s
	}

	saveStringConfigToState := func(apiKey config.Key, attrName schemaAttr) {
		if s, ok := c.Config[apiKey]; ok && s != "" {
			ssh2Config[string(attrName)] = s
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
			ssh2Config[string(attrName)] = int(i)
		}

		delete(swamp, apiKey)
	}

	saveIntConfigToState(config.Port, checkSSH2PortAttr)
	saveStringConfigToState(config.MethodKeyExchange, checkSSH2MethodKexAttr)
	saveStringConfigToState(config.MethodHostKey, checkSSH2MethodHostKeyAttr)
	saveStringConfigToState(config.MethodCryptCS, checkSSH2MethodCryptCSAttr)
	saveStringConfigToState(config.MethodCryptSC, checkSSH2MethodCryptSCAttr)
	saveStringConfigToState(config.MethodMacCS, checkSSH2MethodMacCSAttr)
	saveStringConfigToState(config.MethodMacSC, checkSSH2MethodMacSCAttr)
	saveStringConfigToState(config.MethodCompCS, checkSSH2MethodCompCSAttr)
	saveStringConfigToState(config.MethodCompSC, checkSSH2MethodCompSCAttr)
	// saveStringConfigToState(config.MethodLangCS, checkSSH2MethodLangCSAttr)
	// saveStringConfigToState(config.MethodLangSC, checkSSH2MethodLangSCAttr)

	whitelistedConfigKeys := map[config.Key]struct{}{
		config.ReverseSecretKey:      {},
		config.SubmissionURL:         {},
		config.Key("method_lang_cs"): {},
		config.Key("method_lang_sc"): {},
	}

	for k := range swamp {
		if _, ok := whitelistedConfigKeys[k]; ok {
			delete(c.Config, k)
		}

		if _, ok := whitelistedConfigKeys[k]; !ok {
			return fmt.Errorf("PROVIDER BUG: API Config not empty: %#v", swamp)
		}
	}

	if err := d.Set(checkSSH2Attr, schema.NewSet(checkSSH2ConfigChecksum, []interface{}{ssh2Config})); err != nil {
		return fmt.Errorf("Unable to store check %q attribute: %w", checkSSH2Attr, err)
	}

	return nil
}

// checkSSH2ConfigChecksum creates a stable hash of the normalized values found
// in a user's Terraform config.
func checkSSH2ConfigChecksum(v interface{}) int {
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
	writeInt(checkSSH2PortAttr)
	writeString(checkSSH2MethodKexAttr)
	writeString(checkSSH2MethodHostKeyAttr)
	writeString(checkSSH2MethodCryptCSAttr)
	writeString(checkSSH2MethodCryptSCAttr)
	writeString(checkSSH2MethodMacCSAttr)
	writeString(checkSSH2MethodMacSCAttr)
	writeString(checkSSH2MethodCompCSAttr)
	writeString(checkSSH2MethodCompSCAttr)
	writeString(checkSSH2MethodLangCSAttr)
	writeString(checkSSH2MethodLangSCAttr)

	s := b.String()
	return hashcode.String(s)
}

func checkConfigToAPISSH2(c *circonusCheck, l interfaceList) error { //nolint:unparam
	c.Type = string(apiCheckTypeSSH2)

	for _, mapRaw := range l {
		ssh2Config := newInterfaceMap(mapRaw)

		if v, found := ssh2Config[checkSSH2PortAttr]; found {
			i := v.(int)
			if i != 0 {
				c.Config[config.Port] = fmt.Sprintf("%d", i)
			}
		}

		if v, found := ssh2Config[checkSSH2MethodKexAttr]; found {
			c.Config[config.MethodKeyExchange] = v.(string)
		}

		if v, found := ssh2Config[checkSSH2MethodHostKeyAttr]; found {
			c.Config[config.MethodHostKey] = v.(string)
		}

		if v, found := ssh2Config[checkSSH2MethodCryptCSAttr]; found {
			c.Config[config.MethodCryptCS] = v.(string)
		}

		if v, found := ssh2Config[checkSSH2MethodCryptSCAttr]; found {
			c.Config[config.MethodCryptSC] = v.(string)
		}

		if v, found := ssh2Config[checkSSH2MethodMacCSAttr]; found {
			c.Config[config.MethodMacCS] = v.(string)
		}

		if v, found := ssh2Config[checkSSH2MethodMacSCAttr]; found {
			c.Config[config.MethodMacSC] = v.(string)
		}

		if v, found := ssh2Config[checkSSH2MethodCompCSAttr]; found {
			c.Config[config.MethodCompCS] = v.(string)
		}

		if v, found := ssh2Config[checkSSH2MethodCompSCAttr]; found {
			c.Config[config.MethodCompSC] = v.(string)
		}

		/*
			if v, found := ssh2Config[checkSSH2MethodLangCSAttr]; found {
				c.Config[config.MethodLangCS] = v.(string)
			}

			if v, found := ssh2Config[checkSSH2MethodLangSCAttr]; found {
				c.Config[config.MethodLangSC] = v.(string)
			}
		*/
	}

	return nil
}
