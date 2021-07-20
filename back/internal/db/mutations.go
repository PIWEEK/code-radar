package db

import (
	"path/filepath"
	"github.com/google/uuid"
)

func UpdateFile(path string, dir bool, linesAdd int, linesDelete int, rating float32) error {
	var err error

	file, err := GetFile(path)

	if err != nil {
		return err
	}

	parent := filepath.Dir(path)

	if parent == "." {
		parent = ""
	} else {
		err = UpdateFile(parent, true, linesAdd, linesDelete, rating)

		if err != nil {
			return err
		}
	}

	var filedb FileDB

	if file == nil {
		filedb = FileDB{
			Id: uuid.NewString(),
			Dir: dir,
			Path: path,
			Parent: parent,
			Name: filepath.Base(path),
			Lines: linesAdd,
			Rating: rating,
		}
	} else {
		// Create a copy
		filedb = (*file)
		filedb.Lines += linesAdd
		filedb.Lines -= linesDelete
		filedb.Rating = rating
	}

	err = txn.Insert("files", &filedb)

	return err
}

func MoveFile(path string, newPath string) error {
	return nil
}

func DeleteFile(path string) error {
	return nil
}
