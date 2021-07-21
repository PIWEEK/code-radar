package db

import (
	"github.com/hashicorp/go-memdb"
)

func CreateSchema() *memdb.DBSchema {
	schema := &memdb.DBSchema {
		Tables: map[string]*memdb.TableSchema {
			"files": &memdb.TableSchema {
				Name: "files",
				Indexes: map[string]*memdb.IndexSchema {
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.UUIDFieldIndex{Field: "Id"},
					},
					"path": &memdb.IndexSchema{
						Name:    "path",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Path"},
					},
					"parent": &memdb.IndexSchema{
						Name:    "parent",
						Unique:  false,
						AllowMissing: true,
						Indexer: &memdb.StringFieldIndex{Field: "Parent"},
					},
					"dir": &memdb.IndexSchema{
						Name:    "dir",
						Unique:  false,
						Indexer: &memdb.BoolFieldIndex{Field: "Dir"},
					},
				},
			},
		},
	}
	return schema
}
