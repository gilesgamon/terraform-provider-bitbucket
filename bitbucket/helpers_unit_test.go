package bitbucket

import (
	"reflect"
	"testing"
)

func TestToRelativeEndpoint(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"absolute with query", "https://api.bitbucket.org/2.0/users?page=2", "2.0/users?page=2"},
		{"absolute no query", "https://api.bitbucket.org/2.0/repositories/ws/repo/refs/tags", "2.0/repositories/ws/repo/refs/tags"},
		{"leading slash stripped", "https://api.bitbucket.org/2.0/workspaces", "2.0/workspaces"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := toRelativeEndpoint(c.in); got != c.want {
				t.Errorf("toRelativeEndpoint(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}

func TestEncodeQueryParams(t *testing.T) {
	if got := encodeQueryParams(nil); got != "" {
		t.Errorf("expected empty string for nil params, got %q", got)
	}
	if got := encodeQueryParams(map[string]string{}); got != "" {
		t.Errorf("expected empty string for empty params, got %q", got)
	}

	// Deterministic ordering (keys sorted) and escaping of special characters.
	got := encodeQueryParams(map[string]string{"q": `name = "my repo"`, "sort": "-updated_on"})
	want := "?q=name+%3D+%22my+repo%22&sort=-updated_on"
	if got != want {
		t.Errorf("encodeQueryParams mismatch:\n got: %q\nwant: %q", got, want)
	}
}

func TestWorkspaceRunnerID(t *testing.T) {
	ws, uuid, err := workspaceRunnerId("gob/{1234}")
	if err != nil || ws != "gob" || uuid != "{1234}" {
		t.Fatalf("unexpected: ws=%q uuid=%q err=%v", ws, uuid, err)
	}
	if _, _, err := workspaceRunnerId("no-slash"); err == nil {
		t.Error("expected error for malformed id")
	}
}

func TestRepositoryRunnerID(t *testing.T) {
	ws, repo, uuid, err := repositoryRunnerId("gob/app/{1234}")
	if err != nil || ws != "gob" || repo != "app" || uuid != "{1234}" {
		t.Fatalf("unexpected: ws=%q repo=%q uuid=%q err=%v", ws, repo, uuid, err)
	}
	if _, _, _, err := repositoryRunnerId("gob/app"); err == nil {
		t.Error("expected error for missing uuid segment")
	}
}

func TestProjectDeployKeyID(t *testing.T) {
	ws, key, id, err := projectDeployKeyID("gob/PROJ/42")
	if err != nil || ws != "gob" || key != "PROJ" || id != "42" {
		t.Fatalf("unexpected: ws=%q key=%q id=%q err=%v", ws, key, id, err)
	}
	if _, _, _, err := projectDeployKeyID("gob//42"); err == nil {
		t.Error("expected error for empty project key segment")
	}
}

func TestStringifyMap(t *testing.T) {
	if got := stringifyMap(nil); len(got) != 0 {
		t.Errorf("expected empty map for nil input, got %v", got)
	}

	in := map[string]interface{}{
		"id":      "abc",
		"nested":  map[string]interface{}{"href": "https://example.com"},
		"enabled": true,
	}
	got := stringifyMap(in)
	if got["id"] != "abc" {
		t.Errorf("string value should pass through, got %v", got["id"])
	}
	if got["enabled"] != "true" {
		t.Errorf("bool should be JSON-encoded to \"true\", got %v", got["enabled"])
	}
	if got["nested"] != `{"href":"https://example.com"}` {
		t.Errorf("nested object should be JSON-encoded, got %v", got["nested"])
	}
}

func TestFlattenUserWorkspaces(t *testing.T) {
	var items []WorkspaceAccess
	item := WorkspaceAccess{Administrator: true}
	item.Workspace.UUID = "{uuid}"
	item.Workspace.Slug = "gob"
	item.Workspace.Name = "Gob Bluth"
	item.Workspace.Type = "workspace"
	items = append(items, item)

	got := flattenUserWorkspaces(items)
	want := []interface{}{
		map[string]interface{}{
			"administrator": true,
			"uuid":          "{uuid}",
			"slug":          "gob",
			"name":          "Gob Bluth",
			"type":          "workspace",
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("flattenUserWorkspaces mismatch:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestFlattenFileConflicts(t *testing.T) {
	in := []FileConflict{{Type: "file_conflict", Path: "a.txt", Scenario: "content", Message: "boom"}}
	got := flattenFileConflicts(in)
	if len(got) != 1 {
		t.Fatalf("expected 1 conflict, got %d", len(got))
	}
	m := got[0].(map[string]interface{})
	if m["path"] != "a.txt" || m["scenario"] != "content" || m["message"] != "boom" || m["type"] != "file_conflict" {
		t.Errorf("unexpected flattened conflict: %#v", m)
	}
}
