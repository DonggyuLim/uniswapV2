package db

import (
	"bytes"
	"log"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	db *badger.DB
}

var once sync.Once
var td *DB

func newDB() *DB {
	once.Do(func() {
		var err error
		td := new(DB)
		td.db, err = badger.Open(badger.DefaultOptions("/tmp/token"))
		if err != nil {
			log.Fatal(err)
		}
	})
	return td
}

func (td *DB) Add(key string, value []byte) {
	td.db.Update(func(txn *badger.Txn) error {
		txn.Set([]byte(key), []byte(value))
		return nil
	})
}

func (td *DB) Get(key string) ([]byte, bool) {
	var buf bytes.Buffer
	ok := false
	err := td.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err == nil {
			item.Value(func(val []byte) error {
				buf.Write(val)
				ok = true
				return nil
			})
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes(), ok
}

func (td *DB) Remove(key string) {
	td.db.Update(func(txn *badger.Txn) error {
		txn.Delete([]byte(key))
		return nil
	})
}
