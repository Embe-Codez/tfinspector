package scanner

import (
	"os"
	"path/filepath"
)

func ScanDirectory(root string) ([]TerraformProject, error) {
	var projects []TerraformProject

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".tf" {
			return nil
		}

		project, err := ParseTerraformFile(path)
		if err != nil {
			return nil
		}
		if project != nil {
			projects = append(projects, *project)
		}
		return nil
	})

	return projects, err
}
