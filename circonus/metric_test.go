package circonus

import (
	"testing"
)

func Test_MetricChecksum(t *testing.T) {
	m := interfaceMap{
		string(metricActiveAttr): true,
		string(metricNameAttr):   "asdf",
		string(metricTypeAttr):   "json",
	}

	csum := metricChecksum(m)
	if csum != 4074128010 {
		t.Fatalf("Checksum mismatch")
	}
}
