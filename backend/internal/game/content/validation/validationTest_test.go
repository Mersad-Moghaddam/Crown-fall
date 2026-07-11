package validation

import (
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
