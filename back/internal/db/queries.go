package db

import (
	"github.com/hashicorp/go-memdb"
)

type FileList = [](*FileDB)

func GetFile(path string) (*FileDB, error) {
	raw, err := txn.First("files", "path", path)

	if err != nil {
		return nil, err
	}

	if raw == nil {
		return nil, nil
	}

	return raw.(*FileDB), nil
}

func ProcessResult(it memdb.ResultIterator) FileList {
	var result [](*FileDB)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*FileDB)
		result = append(result, p)
	}
	return result
}

func All() (FileList, error) {
	txn := db.Txn(false)

	it, err := txn.Get("files", "path")

	if err != nil {
		return nil, err
	}

	return ProcessResult(it), nil
}

func ListAllFiles() (FileList, error) {
	it, err := txn.Get("files", "dir", false)

	if err != nil {
		return nil, err
	}

	return ProcessResult(it), nil
}

func ListAllDirectories() (FileList, error) {
	it, err := txn.Get("files", "dir", true)

	if err != nil {
		return nil, err
	}

	return ProcessResult(it), nil
}

func ListFolderFiles(parent string) (FileList, error) {
	it, err := txn.Get("files", "parent", parent)

	if err != nil {
		return nil, err
	}

	return ProcessResult(it), nil
}
