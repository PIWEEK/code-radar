package db

import (
	"github.com/hashicorp/go-memdb"
)

var txn *memdb.Txn

func StartTransaction() {
	txn = db.Txn(true)
}

func Commit() {
	txn.Commit()
	txn = nil
}

func Rollback() {
	txn.Abort()
	txn = nil
}
