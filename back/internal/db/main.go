package db

import (
	"github.com/hashicorp/go-memdb"
)

var db *memdb.MemDB

func InitializeDB() error {
	var err error
	schema := CreateSchema()
	db, err = memdb.NewMemDB(schema)
	return err
}
