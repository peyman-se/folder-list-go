package models

import (
	"fmt"
	"os"
	"path/filepath"
)

type FinderEntity struct {
	os.FileInfo       `json:"-"`
	Name              string `json:"name"`
	Size              int64  `json:"-"`
	HumanReadableSize string `json:"humanReadableSize"`
	Path              string `json:"path"`
	LastModifiedAt    string `json:"lastModifiedAt"`
}

func (entity *FinderEntity) SetHumanReadableEntitySize() {
	entity.HumanReadableSize = GetHumanReadableSize(entity.Size)
}

func (entity *FinderEntity) SetActualSize() error {
	//if it's directory we need to calculate actual size
	if entity.IsDir() {
		var size int64
		path := filepath.Join(entity.Path, entity.Name)
		err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				size += info.Size()
			}

			return err
		})

		if err != nil {
			return err
		}

		entity.Size = size
	}
	entity.SetHumanReadableEntitySize()

	return nil
}

func GetHumanReadableSize(size int64) string {
	const unit = 1000
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "kMGTPE"[exp])
}
