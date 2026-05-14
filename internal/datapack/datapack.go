package datapack

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Datapack struct {
	Namespace string
	ItemTags  map[string][]string
}

type itemTag struct {
	Replace bool     `json:"replace"`
	Values  []string `json:"values"`
}

func (d *Datapack) Save(path string) error {
	itemTagsPath := filepath.Join(path, "data", d.Namespace, "tags", "item")

	err := os.MkdirAll(itemTagsPath, os.ModePerm)
	if err != nil {
		return err
	}

	for tag, values := range d.ItemTags {
		tagPath := filepath.Join(itemTagsPath, tag+".json")
		tagData, err := json.MarshalIndent(itemTag{
			Replace: false,
			Values:  values,
		}, "", "  ")
		if err != nil {
			return err
		}
		err = os.WriteFile(tagPath, tagData, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
