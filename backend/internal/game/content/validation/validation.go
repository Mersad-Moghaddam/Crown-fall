package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"sort"
)

type Document struct {
	SchemaVersion   int      `json:"schemaVersion"`
	ID              string   `json:"id"`
	LocalizationKey string   `json:"localizationKey"`
	RoleIDs         []string `json:"roleIds"`
}

func Load(filesystem fs.FS, paths []string) ([]Document, error) {
	sort.Strings(paths)
	documents := make([]Document, 0, len(paths))
	seen := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		data, err := fs.ReadFile(filesystem, path)
		if err != nil {
			return nil, err
		}
		var document Document
		if err := json.Unmarshal(data, &document); err != nil {
			return nil, fmt.Errorf("decode %s: %w", path, err)
		}
		if document.SchemaVersion != 1 || document.ID == "" || document.LocalizationKey == "" {
			return nil, fmt.Errorf("invalid content document %s", path)
		}
		if _, exists := seen[document.ID]; exists {
			return nil, fmt.Errorf("duplicate content id %s", document.ID)
		}
		seen[document.ID] = struct{}{}
		documents = append(documents, document)
	}
	return documents, nil
}

func ValidateReferences(documents []Document) error {
	ids := make(map[string]struct{}, len(documents))
	for _, document := range documents {
		ids[document.ID] = struct{}{}
	}
	for _, document := range documents {
		for _, roleID := range document.RoleIDs {
			if _, exists := ids[roleID]; !exists {
				return fmt.Errorf("%s references unknown role %s", document.ID, roleID)
			}
		}
	}
	if len(documents) == 0 {
		return errors.New("content set is empty")
	}
	return nil
}
