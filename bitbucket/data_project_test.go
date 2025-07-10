package bitbucket

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceProject_basic(t *testing.T) {
	datasourceName := "data.bitbucket_project.test"
	testTeam := os.Getenv("BITBUCKET_TEAM")
	projectName := os.Getenv("BITBUCKET_PROJECT")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckProject(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDSProjectConfig(testTeam, projectName),
				Check: resource.ComposeTestCheckFunc(
					// testAccCheckBitbucketProjectExists(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "has_publicly_visible_repos", "false"),
					resource.TestCheckResourceAttr(datasourceName, "key", projectName),
					resource.TestCheckResourceAttrSet(datasourceName, "name"),
					resource.TestCheckResourceAttr(datasourceName, "description", ""),
					resource.TestCheckResourceAttrSet(datasourceName, "is_private"),
					resource.TestCheckResourceAttrSet(datasourceName, "uuid"),
					resource.TestCheckResourceAttr(datasourceName, "owner.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "owner.0.display_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "owner.0.uuid"),
					resource.TestCheckResourceAttr(datasourceName, "link.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "link.0.avatar.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "link.0.avatar.0.href"),
				),
			},
		},
	})
}

func testAccBitbucketDSProjectConfig(team, key string) string {
	return fmt.Sprintf(`
data "bitbucket_project" "test" {
  workspace = %[1]q
  key   = %[2]q
}
`, team, key)
}

func TestAccDataSourceProject_Badworkspace(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckProject(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "bitbucket_project" "test" {
					workspace = "badworkspace"
					key   = %[1]q
				}
				`, os.Getenv("BITBUCKET_PROJECT")),
				ExpectError: regexp.MustCompile(".*You may not have access to this project or it no longer exists .*"),
			},
		},
	})
}

func TestAccDataSourceProject_UnknownProject(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckProject(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "bitbucket_project" "test" {
					workspace = %[1]q
					key   = "BADKEY"
				}
				`, os.Getenv("BITBUCKET_TEAM")),
				ExpectError: regexp.MustCompile(".*You may not have access to this project or it no longer exists .*"),
			},
		},
	})
}
