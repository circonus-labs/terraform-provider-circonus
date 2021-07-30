package circonus

// NOTE(sean): One of the objectives of the use of types is to ensure that based
// on aesthetics alone are very few locations where type assertions or casting
// in the main resource files is required (mainly when interacting with the
// external API structs).  As a rule of thumb, all type assertions should happen
// in the utils file and casting is only done at assignment time when storing a
// result to a struct.  Said differently, contained tedium should enable
// compiler enforcement of types and easy verification.

type (
	apiCheckType      string
	attrDescr         string
	attrDescrs        map[schemaAttr]attrDescr
	schemaAttr        string
	metricID          string
	validString       string
	validStringValues []validString
)
