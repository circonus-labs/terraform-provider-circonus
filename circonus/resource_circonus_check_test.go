package circonus

import (
	"fmt"

	api "github.com/circonus-labs/go-apiclient"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func testAccCheckDestroyCirconusCheckBundle(s *terraform.State) error {
	c := testAccProvider.Meta().(*providerContext)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "circonus_check" {
			continue
		}

		// for getting resource ids for acceptance tests
		// if t, ok := rs.Primary.Attributes["type"]; ok {
		// 	fmt.Println(t)
		// 	if t == "http" {
		// 		fmt.Println(rs.Primary.String())
		// 	}
		// }

		cid := rs.Primary.ID
		exists, err := checkCheckBundleExists(c, api.CIDType(&cid))
		if err != nil {
			return fmt.Errorf("Error checking check bundle %s", err)
		}

		if exists {
			return fmt.Errorf("check bundle still exists after destroy")
		}
	}

	return nil
}

func checkCheckBundleExists(c *providerContext, checkBundleID api.CIDType) (bool, error) {
	cb, err := c.client.FetchCheckBundle(checkBundleID)
	if err != nil {
		return false, err
	}

	if api.CIDType(&cb.CID) == checkBundleID {
		return true, nil
	}

	return false, nil
}
