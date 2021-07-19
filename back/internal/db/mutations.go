package db

import (
	"path/filepath"

	"github.com/google/uuid"
)

func InsertFile(path string, dir bool) error {
	raw, err := txn.First("files", "path", path)

	parent := filepath.Dir(path)

	if err == nil && raw == nil {
		filedb := &FileDB{
			Id: uuid.NewString(),
			Dir: dir,
			Path: path,
			Parent: parent,
			Name: filepath.Base(path),
		}
		err = txn.Insert("files", filedb)
	}

	return err
}

