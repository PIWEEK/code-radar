package db

func ListFiles() ([](*FileDB), error) {
	t := db.Txn(false)

	it, err := t.Get("files", "parent")

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

