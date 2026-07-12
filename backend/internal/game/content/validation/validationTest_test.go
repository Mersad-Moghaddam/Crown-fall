package validation

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestDeterministicLoadingAndReferences(t *testing.T) {
	files := fstest.MapFS{
		"role.json":     {Data: []byte(`{"schemaVersion":1,"id":"role.one","localizationKey":"role.one"}`)},
		"scenario.json": {Data: []byte(`{"schemaVersion":1,"id":"scenario.one","localizationKey":"scenario.one","roleIds":["role.one"]}`)},
	}
	documents, err := Load(files, []string{"scenario.json", "role.json"})
	if err != nil {
		t.Fatal(err)
	}
	if documents[0].ID != "role.one" {
		t.Fatal("paths were not loaded deterministically")
	}
	if err := ValidateReferences(documents); err != nil {
		t.Fatal(err)
	}
}

func TestRejectsDuplicateIDs(t *testing.T) {
	files := fstest.MapFS{"one.json": {Data: []byte(`{"schemaVersion":1,"id":"same","localizationKey":"one"}`)}, "two.json": {Data: []byte(`{"schemaVersion":1,"id":"same","localizationKey":"two"}`)}}
	if _, err := Load(files, []string{"one.json", "two.json"}); err == nil || !strings.Contains(err.Error(), "duplicate") {
		t.Fatalf("expected duplicate error, got %v", err)
	}
}

func TestRejectsUnsupportedSchema(t *testing.T) {
	files := fstest.MapFS{"role.json": {Data: []byte(`{"schemaVersion":2,"id":"role","localizationKey":"role"}`)}}
	if _, err := Load(files, []string{"role.json"}); err == nil || !strings.Contains(err.Error(), "invalid content") {
		t.Fatalf("expected schema error, got %v", err)
	}
}

func TestRejectsMissingReference(t *testing.T) {
	documents := []Document{{SchemaVersion: 1, ID: "scenario", LocalizationKey: "scenario", RoleIDs: []string{"missing"}}}
	if err := ValidateReferences(documents); err == nil || !strings.Contains(err.Error(), "unknown role") {
		t.Fatalf("expected reference error, got %v", err)
	}
}
