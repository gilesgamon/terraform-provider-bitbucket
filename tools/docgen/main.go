package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-bitbucket/bitbucket"
)

func typeName(t schema.ValueType) string {
	switch t {
	case schema.TypeBool:
		return "Boolean"
	case schema.TypeInt:
		return "Integer"
	case schema.TypeFloat:
		return "Float"
	case schema.TypeString:
		return "String"
	case schema.TypeList:
		return "List"
	case schema.TypeMap:
		return "Map"
	case schema.TypeSet:
		return "Set"
	default:
		return "String"
	}
}

func humanize(name string) string {
	s := strings.TrimPrefix(name, "bitbucket_")
	s = strings.ReplaceAll(s, "_", " ")
	return s
}

func placeholder(name string, t schema.ValueType) string {
	switch t {
	case schema.TypeBool:
		return "true"
	case schema.TypeInt, schema.TypeFloat:
		return "1"
	case schema.TypeList, schema.TypeSet:
		return "[]"
	case schema.TypeMap:
		return "{}"
	}
	switch name {
	case "workspace":
		return "\"example-workspace\""
	case "repo_slug", "repository":
		return "\"example-repo\""
	case "project_key":
		return "\"PROJ\""
	case "pull_request_id":
		return "\"1\""
	case "commit":
		return "\"a1b2c3d4\""
	default:
		return fmt.Sprintf("\"%s\"", name)
	}
}

// classify returns required args, optional args, and computed attributes.
func classify(s map[string]*schema.Schema) (req, opt, comp []string) {
	for name, attr := range s {
		switch {
		case attr.Required:
			req = append(req, name)
		case attr.Optional:
			opt = append(opt, name)
		default:
			comp = append(comp, name)
		}
	}
	sort.Strings(req)
	sort.Strings(opt)
	sort.Strings(comp)
	return
}

func describe(name string, attr *schema.Schema) string {
	if attr.Description != "" {
		return attr.Description
	}
	return fmt.Sprintf("The %s.", strings.ReplaceAll(name, "_", " "))
}

func renderDoc(typeName2, kind string, res *schema.Resource) string {
	name := typeName2
	human := humanize(name)
	escaped := strings.ReplaceAll(name, "_", "\\_")
	sidebarName := strings.ReplaceAll(strings.TrimPrefix(name, "bitbucket_"), "_", "-")

	var blurb string
	if kind == "resource" {
		blurb = fmt.Sprintf("Provides a Bitbucket %s resource.", human)
	} else {
		blurb = fmt.Sprintf("Provides information about Bitbucket %s.", human)
	}

	req, opt, comp := classify(res.Schema)

	var b strings.Builder
	fmt.Fprintf(&b, "---\n")
	fmt.Fprintf(&b, "layout: \"bitbucket\"\n")
	fmt.Fprintf(&b, "page_title: \"Bitbucket: %s\"\n", name)
	fmt.Fprintf(&b, "sidebar_current: \"docs-bitbucket-%s-%s\"\n", kind, sidebarName)
	fmt.Fprintf(&b, "description: |-\n  %s\n", blurb)
	fmt.Fprintf(&b, "---\n\n")
	fmt.Fprintf(&b, "# %s\n\n", escaped)
	fmt.Fprintf(&b, "%s\n\n", blurb)

	// Example usage
	fmt.Fprintf(&b, "## Example Usage\n\n")
	fmt.Fprintf(&b, "```hcl\n")
	if kind == "resource" {
		fmt.Fprintf(&b, "resource \"%s\" \"example\" {\n", name)
	} else {
		fmt.Fprintf(&b, "data \"%s\" \"example\" {\n", name)
	}
	for _, r := range req {
		fmt.Fprintf(&b, "  %s = %s\n", r, placeholder(r, res.Schema[r].Type))
	}
	fmt.Fprintf(&b, "}\n")
	fmt.Fprintf(&b, "```\n\n")

	// Argument reference
	fmt.Fprintf(&b, "## Argument Reference\n\n")
	if len(req)+len(opt) == 0 {
		fmt.Fprintf(&b, "This %s takes no arguments.\n\n", kind)
	} else {
		fmt.Fprintf(&b, "The following arguments are supported:\n\n")
		for _, r := range req {
			fmt.Fprintf(&b, "* `%s` - (Required) %s\n", r, describe(r, res.Schema[r]))
		}
		for _, o := range opt {
			fmt.Fprintf(&b, "* `%s` - (Optional) %s\n", o, describe(o, res.Schema[o]))
		}
		fmt.Fprintf(&b, "\n")
	}

	// Attributes reference
	fmt.Fprintf(&b, "## Attributes Reference\n\n")
	if len(comp) == 0 {
		fmt.Fprintf(&b, "Only the arguments listed above are exposed as attributes.\n")
	} else {
		fmt.Fprintf(&b, "In addition to the arguments above, the following attributes are exported:\n\n")
		fmt.Fprintf(&b, "* `id` - The identifier of the %s.\n", human)
		for _, c := range comp {
			if c == "id" {
				continue
			}
			attr := res.Schema[c]
			t := attr.Type
			nested := nestedResource(attr)
			suffix := ""
			if (t == schema.TypeList || t == schema.TypeSet) && nested != nil {
				suffix = " Each item contains:"
			} else if t == schema.TypeList || t == schema.TypeSet {
				suffix = " (a " + strings.ToLower(typeName(t)) + " of values)"
			}
			fmt.Fprintf(&b, "* `%s` - %s%s\n", c, describe(c, attr), suffix)
			if nested != nil {
				keys := make([]string, 0, len(nested.Schema))
				for k := range nested.Schema {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Fprintf(&b, "    * `%s` - %s\n", k, describe(k, nested.Schema[k]))
				}
			}
		}
	}

	return b.String()
}

// nestedResource returns the *schema.Resource used as the element of a
// list/set attribute, or nil when the element is a primitive.
func nestedResource(attr *schema.Schema) *schema.Resource {
	if attr == nil {
		return nil
	}
	if r, ok := attr.Elem.(*schema.Resource); ok {
		return r
	}
	return nil
}

func docSet(dir string) map[string]bool {
	out := map[string]bool{}
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".md") {
			out[strings.TrimSuffix(e.Name(), ".md")] = true
		}
	}
	return out
}

func main() {
	p := bitbucket.Provider()

	type target struct {
		kind string
		dir  string
		m    map[string]*schema.Resource
	}
	targets := []target{
		{"resource", "docs/resources", p.ResourcesMap},
		{"data-source", "docs/data-sources", p.DataSourcesMap},
	}

	created := 0
	for _, t := range targets {
		existing := docSet(t.dir)
		for name, res := range t.m {
			short := strings.TrimPrefix(name, "bitbucket_")
			if existing[short] {
				continue
			}
			// kind for frontmatter/sidebar: "data" or "resource"
			fmKind := "resource"
			if t.kind == "data-source" {
				fmKind = "data"
			}
			content := renderDoc(name, fmKind, res)
			path := filepath.Join(t.dir, short+".md")
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				fmt.Fprintln(os.Stderr, "error writing", path, err)
				os.Exit(1)
			}
			created++
			fmt.Println("generated", path)
		}
	}
	fmt.Printf("done: %d docs generated\n", created)
}
