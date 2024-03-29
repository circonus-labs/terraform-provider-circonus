package circonus

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type (
	circonusTag  string
	circonusTags []circonusTag
)

// tagMakeConfigSchema returns a schema pointer to the necessary tag structure.
func tagMakeConfigSchema(tagAttrName schemaAttr) *schema.Schema { //nolint:unparam
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validateTag,
		},
	}
}

func (t circonusTag) Category() string {
	tagInfo := strings.SplitN(string(t), ":", 2)
	switch len(tagInfo) {
	case 1:
		return strings.ToLower(string(t))
	case 2:
		return strings.ToLower(tagInfo[0])
	default:
		log.Printf("[ERROR]: Invalid category on tag %q", string(t))
		return ""
	}
}

func (t circonusTag) Value() string {
	tagInfo := strings.SplitN(string(t), ":", 2)
	switch len(tagInfo) {
	case 1:
		return ""
	case 2:
		return strings.ToLower(tagInfo[1])
	default:
		log.Printf("[ERROR]: Invalid value on tag %q", string(t))
		return ""
	}
}

func tagsToState(tags circonusTags) *schema.Set {
	tagSet := schema.NewSet(schema.HashString, nil)
	for i := range tags {
		tagSet.Add(strings.ToLower(string(tags[i])))
	}
	return tagSet
}

func apiToTags(apiTags []string) circonusTags {
	tags := make(circonusTags, 0, len(apiTags))
	for _, v := range apiTags {
		tags = append(tags, circonusTag(v))
	}
	return tags
}
