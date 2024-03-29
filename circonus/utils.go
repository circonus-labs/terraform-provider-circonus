package circonus

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/circonus-labs/go-apiclient/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// convertToHelperSchema converts the schema and injects the necessary
// parameters, notably the descriptions, in order to be valid input to
// Terraform's helper schema.
func convertToHelperSchema(descrs attrDescrs, in map[schemaAttr]*schema.Schema) map[string]*schema.Schema {
	out := make(map[string]*schema.Schema, len(in))
	for k, v := range in {
		if descr, ok := descrs[k]; ok {
			// NOTE(sean@): At some point this check needs to be uncommented and all
			// empty descriptions need to be populated.
			//
			// if len(descr) == 0 {
			// 	log.Printf("[WARN] PROVIDER BUG: Description of attribute %s empty", k)
			// }

			v.Description = string(descr)
		} else {
			log.Printf("[WARN] PROVIDER BUG: Unable to find description for attr %q", k)
		}

		out[string(k)] = v
	}

	return out
}

func failoverGroupIDToCID(groupID int) string {
	if groupID == 0 {
		return ""
	}

	return fmt.Sprintf("%s/%d", config.ContactGroupPrefix, groupID)
}

func failoverGroupCIDToID(cid api.CIDType) (int, error) {
	re := regexp.MustCompile("^" + config.ContactGroupPrefix + "/(" + config.DefaultCIDRegex + ")$")
	matches := re.FindStringSubmatch(*cid)
	if matches == nil || len(matches) < 2 {
		return -1, fmt.Errorf("Did not find a valid contact_group ID in the CID %q", *cid)
	}

	contactGroupID, err := strconv.Atoi(matches[1])
	if err != nil {
		return -1, fmt.Errorf("invalid contact_group ID: unable to find an ID in %q: %w", *cid, err)
	}

	return contactGroupID, nil
}

// flattenList returns a list of all string values to a []*string.
func flattenList(l []interface{}) []*string {
	vals := make([]*string, 0, len(l))
	for _, v := range l {
		val, ok := v.(string)
		if ok && val != "" {
			vals = append(vals, &val)
		}
	}
	return vals
}

// flattenSet flattens the values in a schema.Set and returns a []*string.
func flattenSet(s *schema.Set) []*string {
	return flattenList(s.List())
}

// importStatePassthrough is an implementation of StateFunc that can be used to
// simply pass the ID directly through. This should be used only in the case
// that an ID-only refresh is possible.  The ID is url.PathUnescape()'ed.
func importStatePassthroughUnescape(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	// Ignore any path unescape issues
	cid, _ := url.PathUnescape(d.Id())

	d.SetId(cid)

	return []*schema.ResourceData{d}, nil
}

func derefStringList(lp []*string) []string {
	l := make([]string, 0, len(lp))
	for _, sp := range lp {
		if sp != nil {
			l = append(l, *sp)
		}
	}
	return l
}

// listToSet returns a TypeSet from the given list.
func stringListToSet(stringList []string, keyName schemaAttr) []interface{} {
	m := make([]interface{}, 0, len(stringList))
	for _, v := range stringList {
		s := make(map[string]interface{}, 1)
		s[string(keyName)] = v
		m = append(m, s)
	}

	return m
}

func normalizeTimeDurationStringToSeconds(v interface{}) string {
	switch v := v.(type) {
	case string:
		d, err := time.ParseDuration(v)
		if err != nil {
			return fmt.Sprintf("<unable to normalize time duration %s: %v>", v, err)
		}

		return fmt.Sprintf("%ds", int(d.Seconds()))
	default:
		return fmt.Sprintf("<unable to normalize duration on %#v>", v)
	}
}

func indirect(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		return v
	case *string:
		if v == nil {
			return nil
		}
		return *v
	default:
		return v
	}
}

func suppressEquivalentTimeDurations(k, old, update string, d *schema.ResourceData) bool {
	d1, err := time.ParseDuration(old)
	if err != nil {
		return false
	}

	d2, err := time.ParseDuration(update)
	if err != nil {
		return false
	}

	return d1 == d2
}

func suppressWhitespace(v interface{}) string {
	return strings.TrimSpace(v.(string))
}

func jsonSort(v interface{}) string {
	var ifce interface{}
	ob := []byte(v.(string))
	err := json.Unmarshal(ob, &ifce)
	if err != nil {
		return v.(string)
	}
	os, err := json.Marshal(ifce)
	if err != nil {
		return v.(string)
	}
	return string(os)
}
