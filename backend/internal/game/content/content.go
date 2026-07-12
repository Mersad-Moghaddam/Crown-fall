package content

import (
	"embed"
	"fmt"
	"io/fs"

	"crownfall/backend/internal/game/content/validation"
)

//go:embed roles/*.json scenarios/*.json events/*.json
var files embed.FS

func Validate() error {
	paths := make([]string, 0)
	for _, pattern := range []string{"roles/*.json", "scenarios/*.json", "events/*.json"} {
		matches, err := fs.Glob(files, pattern)
		if err != nil {
			return err
		}
		paths = append(paths, matches...)
	}
	documents, err := validation.Load(files, paths)
	if err != nil {
		return fmt.Errorf("load game content: %w", err)
	}
	if err := validation.ValidateReferences(documents); err != nil {
		return fmt.Errorf("validate game content: %w", err)
	}
	return nil
}
