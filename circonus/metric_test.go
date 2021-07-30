package circonus

import "testing"

func Test_MetricChecksum(t *testing.T) {
	m := interfaceMap{
		string(metricActiveAttr): true,
		string(metricNameAttr):   "asdf",
		string(metricTypeAttr):   "json",
	}

	if csum := metricChecksum(m); csum != 4074128010 {
		t.Fatalf("Checksum mismatch")
	}
}
