package circonus

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider

	testAccAccount string

	testAccBroker1 string
	testAccBroker2 string
	testAccBroker3 string
	testAccBroker4 string

	testAccContactGroup1 string
	testAccContactGroup2 string
	testAccContactGroup3 string
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"circonus": testAccProvider,
	}

	if testAccAccount = os.Getenv("CIRCONUS_TEST_ACCOUNT"); testAccAccount == "" {
		testAccAccount = "/account/4536"
	}

	if testAccBroker1 = os.Getenv("CIRCONUS_TEST_BROKER_1"); testAccBroker1 == "" {
		testAccBroker1 = "/broker/1"
	}

	if testAccBroker2 = os.Getenv("CIRCONUS_TEST_BROKER_2"); testAccBroker2 == "" {
		testAccBroker2 = "/broker/275"
	}

	if testAccBroker3 = os.Getenv("CIRCONUS_TEST_BROKER_3"); testAccBroker3 == "" {
		testAccBroker3 = "/broker/35"
	}

	if testAccBroker4 = os.Getenv("CIRCONUS_TEST_BROKER_4"); testAccBroker4 == "" {
		testAccBroker4 = "/broker/1490"
	}

	if testAccContactGroup1 = os.Getenv("CIRCONUS_TEST_CONTACT_GROUP_1"); testAccContactGroup1 == "" {
		testAccContactGroup1 = "/contact_group/4661"
	}

	if testAccContactGroup2 = os.Getenv("CIRCONUS_TEST_CONTACT_GROUP_2"); testAccContactGroup2 == "" {
		testAccContactGroup2 = "/contact_group/4679"
	}

	if testAccContactGroup3 = os.Getenv("CIRCONUS_TEST_CONTACT_GROUP_3"); testAccContactGroup3 == "" {
		testAccContactGroup3 = "/contact_group/4680"
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if apiToken := os.Getenv("CIRCONUS_API_TOKEN"); apiToken == "" {
		t.Fatal("CIRCONUS_API_TOKEN must be set for acceptance tests")
	}
}
