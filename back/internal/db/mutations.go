package db

import (
	"time"
	"path/filepath"
	"github.com/google/uuid"
)

func UpdateFile(path string, dir bool, linesAdd int, linesDelete int, user string, date time.Time) error {
	var err error

	file, err := GetFile(path)

	if err != nil {
		return err
	}

	var parent string

	if path != "." {
		parent = filepath.Dir(path)
		err = UpdateFile(parent, true, linesAdd, linesDelete, user, date)

		if err != nil {
			return err
		}
	}

	var filedb *FileDB

	if file == nil {
		filedb = &FileDB{
			Id: uuid.NewString(),
			Dir: dir,
			Path: path,
			Parent: parent,
			Name: filepath.Base(path),
			Lines: linesAdd,
			History: []FileHistory{
				FileHistory{
					User: user,
					Added: linesAdd,
					Deleted: linesDelete,
					Date: date,
				},
			},
		}
		err = txn.Insert("files", filedb)
	} else {
		// Create a copy
		filedb = file.Copy()
		filedb.Lines += linesAdd
		filedb.Lines -= linesDelete
		filedb.History = append(filedb.History, FileHistory{
			User: user,
			Added: linesAdd,
			Deleted: linesDelete,
			Date: date,
		})

		if (filedb.Dir && filedb.Lines == 0) {
			err = txn.Delete("files", filedb)
		} else {
			err = txn.Insert("files", filedb)
		}
	}

	return err
}

func MoveFile(path string, newPath string, user string, date time.Time) error {
	file, err := GetFile(path)

	if file != nil && err == nil {
		if (file.Parent != "") {
			// Remove lines from parent
			err = UpdateFile(file.Parent, true, 0, file.Lines, user, date)

			if err != nil {
				return err
			}
		}

		newParent := filepath.Dir(newPath)

		if newParent != "." {
			err = UpdateFile(newParent, true, file.Lines, 0, user, date)
		}

		if err != nil {
			return err
		}

		// copy the file
		newFile := file.Copy()
		newFile.Path = newPath
		newFile.Parent = newParent

		err = txn.Insert("files", newFile)
	}

	return err
}

func DeleteFile(path string, user string, date time.Time) error {
	file, err := GetFile(path)

	if file != nil && err == nil {
		if (file.Parent != "") {
			err = UpdateFile(file.Parent, true, 0, file.Lines, user, date)

			if err != nil {
				return err
			}
		}

		err = txn.Delete("files", file)
	}

	return err
}

func SaveFile(file *FileDB) error {
	return txn.Insert("files", file)
}
