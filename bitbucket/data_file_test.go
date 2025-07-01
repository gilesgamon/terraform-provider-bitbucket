package bitbucket

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFile_basic(t *testing.T) {
	dataSourceName := "data.bitbucket_file.test"
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	commit := os.Getenv("BITBUCKET_COMMIT")
	path := os.Getenv("BITBUCKET_PATH")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckRepo(t)
			testAccPreCheckFileCommit(t)
			testAccPreCheckFilePath(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDataFileConfig(workspace, repository, commit, path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "repo_slug", repository),
					resource.TestCheckResourceAttr(dataSourceName, "workspace", workspace),
					resource.TestCheckResourceAttrSet(dataSourceName, "commit"),
					resource.TestCheckResourceAttr(dataSourceName, "path", "README.md"),
					resource.TestCheckResourceAttrSet(dataSourceName, "content"),
					resource.TestCheckResourceAttrSet(dataSourceName, "content_b64"),
					resource.TestCheckNoResourceAttr(dataSourceName, "metadata"),
				),
			},
		},
	})
}

func testAccBitbucketDataFileConfig(workspace string, repository string, commit string, path string) string {
	return fmt.Sprintf(`

data "bitbucket_file" "test" {
	workspace = %[1]q
	repo_slug = %[2]q
	commit = %[3]q
	path = %[4]q
}
`, workspace, repository, commit, path)
}

func TestAccDataSourceFile_Metadata_basic(t *testing.T) {
	dataSourceName := "data.bitbucket_file.test"
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	commit := os.Getenv("BITBUCKET_COMMIT")
	path := os.Getenv("BITBUCKET_PATH")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckRepo(t); testAccPreCheckFilePath(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDataFileConfigMeta(workspace, repository, commit, path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "repo_slug", repository),
					resource.TestCheckResourceAttr(dataSourceName, "workspace", workspace),
					resource.TestCheckResourceAttrSet(dataSourceName, "commit"),
					resource.TestCheckResourceAttr(dataSourceName, "path", path),
					resource.TestCheckNoResourceAttr(dataSourceName, "content"),
					resource.TestCheckNoResourceAttr(dataSourceName, "content_b64"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.path", path),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.type", "commit_file"),
				),
			},
		},
	})
}

func testAccBitbucketDataFileConfigMeta(workspace string, repository string, commit string, path string) string {
	return fmt.Sprintf(`

data "bitbucket_file" "test" {
	workspace = %[1]q
	repo_slug = %[2]q
	commit = %[3]q
	path = %[4]q
	format = "meta"
}
`, workspace, repository, commit, path)
}

func TestAccDataSourceFile_Metadata_IncludeLinks(t *testing.T) {
	dataSourceName := "data.bitbucket_file.test"
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	commit := os.Getenv("BITBUCKET_COMMIT")
	path := os.Getenv("BITBUCKET_PATH")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckRepo(t); testAccPreCheckFilePath(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDataFileConfigMetaIncludeLinks(workspace, repository, commit, path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "repo_slug", repository),
					resource.TestCheckResourceAttr(dataSourceName, "workspace", workspace),
					resource.TestCheckResourceAttrSet(dataSourceName, "commit"),
					resource.TestCheckResourceAttr(dataSourceName, "path", path),
					resource.TestCheckNoResourceAttr(dataSourceName, "content"),
					resource.TestCheckNoResourceAttr(dataSourceName, "content_b64"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.path", path),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.type", "commit_file"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.link.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.link.0.self.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.link.0.self.0.href"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.link.0.history.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.link.0.history.0.href"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.link.0.meta.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.link.0.meta.0.href"),
				),
			},
		},
	})
}

func testAccBitbucketDataFileConfigMetaIncludeLinks(workspace string, repository string, commit string, path string) string {
	return fmt.Sprintf(`

data "bitbucket_file" "test" {
	workspace = %[1]q
	repo_slug = %[2]q
	commit = %[3]q
	path = %[4]q
	format = "meta"
	include_links = true
}
`, workspace, repository, commit, path)
}

func TestAccDataSourceFile_Metadata_IncludeCommit(t *testing.T) {
	dataSourceName := "data.bitbucket_file.test"
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	commit := os.Getenv("BITBUCKET_COMMIT")
	path := os.Getenv("BITBUCKET_PATH")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckRepo(t); testAccPreCheckFilePath(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDataFileConfigMetaIncludeCommit(workspace, repository, commit, path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "repo_slug", repository),
					resource.TestCheckResourceAttr(dataSourceName, "workspace", workspace),
					resource.TestCheckResourceAttrSet(dataSourceName, "commit"),
					resource.TestCheckResourceAttr(dataSourceName, "path", path),
					resource.TestCheckNoResourceAttr(dataSourceName, "content"),
					resource.TestCheckNoResourceAttr(dataSourceName, "content_b64"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.path", path),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.type", "commit_file"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.commit.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.commit.0.hash"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.commit.0.type"),
				),
			},
		},
	})
}

func testAccBitbucketDataFileConfigMetaIncludeCommit(workspace string, repository string, commit string, path string) string {
	return fmt.Sprintf(`

data "bitbucket_file" "test" {
	workspace = %[1]q
	repo_slug = %[2]q
	commit = %[3]q
	path = %[4]q
	format = "meta"
	include_commit = true
}
`, workspace, repository, commit, path)
}

func TestAccDataSourceFile_Metadata_IncludeCommitLinks(t *testing.T) {
	dataSourceName := "data.bitbucket_file.test"
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	commit := os.Getenv("BITBUCKET_COMMIT")
	path := os.Getenv("BITBUCKET_PATH")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckRepo(t); testAccPreCheckFilePath(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDataFileConfigMetaIncludeCommitLinks(workspace, repository, commit, path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "repo_slug", repository),
					resource.TestCheckResourceAttr(dataSourceName, "workspace", workspace),
					resource.TestCheckResourceAttrSet(dataSourceName, "commit"),
					resource.TestCheckResourceAttr(dataSourceName, "path", path),
					resource.TestCheckNoResourceAttr(dataSourceName, "content"),
					resource.TestCheckNoResourceAttr(dataSourceName, "content_b64"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.path", path),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.type", "commit_file"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.commit.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.commit.0.hash"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.commit.0.type"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.commit.0.link.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.commit.0.link.0.self.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.commit.0.link.0.self.0.href"),
					resource.TestCheckResourceAttr(dataSourceName, "metadata.0.commit.0.link.0.html.#", "1"),
					resource.TestCheckResourceAttrSet(dataSourceName, "metadata.0.commit.0.link.0.html.0.href"),
				),
			},
		},
	})
}

func testAccBitbucketDataFileConfigMetaIncludeCommitLinks(workspace string, repository string, commit string, path string) string {
	return fmt.Sprintf(`

data "bitbucket_file" "test" {
	workspace = %[1]q
	repo_slug = %[2]q
	commit = %[3]q
	path = %[4]q
	format = "meta"
	include_commit = true
	include_commit_links = true
}
`, workspace, repository, commit, path)
}

func TestAccDataSourceFile_Metadata_RequestCommitLinksNoCommit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				data "bitbucket_file" "test" {
					workspace = "test"
					repo_slug = "test"
					commit = "test"
					path = "test"
					format = "meta"
					include_commit_links = true
				}
				`,
				ExpectError: regexp.MustCompile(".*include_commit_links cannot be true if include_commit is not set to true.*"),
			},
		},
	})
}

func TestAccDataSourceFile_BadCommit(t *testing.T) {
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	commit := os.Getenv("BITBUCKET_COMMIT")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckRepo(t); testAccPreCheckFileCommit(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBitbucketDataFileConfigBadCommit(workspace, repository, commit),
				ExpectError: regexp.MustCompile(".*Commit not found.*"),
			},
		},
	})
}

func testAccBitbucketDataFileConfigBadCommit(workspace string, repository string, commit string) string {
	return fmt.Sprintf(`
data "bitbucket_file" "test" {
	workspace = %[1]q
	repo_slug = %[2]q
	commit = "badcommit"
	path = "README.md"
}
`, workspace, repository, commit)
}

func TestAccDataSourceFile_BadFile(t *testing.T) {
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	commit := os.Getenv("BITBUCKET_COMMIT")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckRepo(t); testAccPreCheckFileCommit(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccBitbucketDataFileConfigBadFile(workspace, repository, commit),
				ExpectError: regexp.MustCompile(".*No such file.*"),
			},
		},
	})
}

func testAccBitbucketDataFileConfigBadFile(workspace string, repository string, commit string) string {
	return fmt.Sprintf(`
data "bitbucket_file" "test" {
	workspace = %[1]q
	repo_slug = %[2]q
	commit = %[3]q
	path = "badfile"
}
`, workspace, repository, commit)
}

func TestAccDataSourceFile_LookupCommit(t *testing.T) {
	dataSourceName := "data.bitbucket_file.test"
	workspace := os.Getenv("BITBUCKET_TEAM")
	repository := os.Getenv("BITBUCKET_REPO")
	path := os.Getenv("BITBUCKET_PATH")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccPreCheckRepo(t); testAccPreCheckFilePath(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBitbucketDataFileConfigLookupCommit(workspace, repository, path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "repo_slug", repository),
					resource.TestCheckResourceAttr(dataSourceName, "workspace", workspace),
					resource.TestCheckResourceAttrSet(dataSourceName, "commit"),
					resource.TestCheckResourceAttr(dataSourceName, "path", path),
					resource.TestCheckResourceAttrSet(dataSourceName, "content"),
					resource.TestCheckResourceAttrSet(dataSourceName, "content_b64"),
					resource.TestCheckNoResourceAttr(dataSourceName, "metadata"),
				),
			},
		},
	})
}

func testAccBitbucketDataFileConfigLookupCommit(workspace string, repository string, path string) string {
	return fmt.Sprintf(`

data "bitbucket_repository" "test" {
	workspace = %[1]q
	repo_slug = %[2]q
}

data "bitbucket_file" "test" {
	workspace = data.bitbucket_repository.test.workspace
	repo_slug = data.bitbucket_repository.test.repo_slug
	commit = data.bitbucket_repository.test.main_branch
	path = %[3]q
}
`, workspace, repository, path)
}
