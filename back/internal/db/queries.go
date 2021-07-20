package db

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

func ListFiles() ([](*FileDB), error) {
	txn := db.Txn(false)

	it, err := txn.Get("files", "path")

	if err != nil {
		return nil, err
	}

	var result [](*FileDB)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*FileDB)
		result = append(result, p)
	}

	return result, nil
}

