package bitbucket

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceRepository_basic(t *testing.T) {
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	datasourceName := "data.bitbucket_repository.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckRepo(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDataRepositoryConfig(workspace, repository),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBitbucketRepositoryExists(datasourceName),
					resource.TestCheckResourceAttr(datasourceName, "name", repository),
					resource.TestCheckResourceAttr(datasourceName, "full_name", fmt.Sprintf("%s/%s", workspace, repository)),
					resource.TestCheckResourceAttr(datasourceName, "scm", "git"),
					resource.TestCheckResourceAttrSet(datasourceName, "has_wiki"),
					resource.TestCheckResourceAttrSet(datasourceName, "uuid"),
					resource.TestCheckResourceAttr(datasourceName, "fork_policy", "allow_forks"),
					resource.TestCheckResourceAttr(datasourceName, "language", ""),
					resource.TestCheckResourceAttrSet(datasourceName, "has_issues"),
					resource.TestCheckResourceAttrSet(datasourceName, "is_private"),
					resource.TestCheckResourceAttrSet(datasourceName, "description"),
					resource.TestCheckResourceAttrSet(datasourceName, "main_branch"),
					resource.TestCheckResourceAttr(datasourceName, "link.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "link.0.avatar.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "link.0.avatar.0.href"),
					resource.TestCheckResourceAttr(datasourceName, "project.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "project.0.name"),
					resource.TestCheckResourceAttr(datasourceName, "owner.#", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "owner.0.display_name"),
					resource.TestCheckResourceAttrSet(datasourceName, "owner.0.uuid"),
				),
			},
		},
	})
}
func testAccBitbucketDataRepositoryConfig(workspace string, repoName string) string {
	return fmt.Sprintf(`

data "bitbucket_workspace" "test" {
  workspace = %[1]q
}

data "bitbucket_repository" "test" {
  workspace = data.bitbucket_workspace.test.slug
  repo_slug = %[2]q
}
`, workspace, repoName)
}

func TestAccDataSourceRepository_UnknownWorkspace(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); testAccPreCheckRepo(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "bitbucket_repository" "test" {
					workspace = "badworkspace"
					repo_slug   = %[1]q
				}
				`, os.Getenv("BITBUCKET_REPO")),
				ExpectError: regexp.MustCompile(".*You may not have access to this repository or it no longer exists in this workspace..*"),
			},
		},
	})
}

func TestAccDataSourceRepository_UnknownRepo(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBitbucketProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "bitbucket_repository" "test" {
					workspace = %[1]q
					repo_slug   = "badproject"
				}
				`, os.Getenv("BITBUCKET_TEAM")),
				ExpectError: regexp.MustCompile(".*You may not have access to this repository or it no longer exists in this workspace..*"),
			},
		},
	})
}
