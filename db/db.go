package db

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	db *badger.DB
}

var once sync.Once
var db *DB

func NewDB(dbname string) *DB {
	once.Do(func() {
		var err error
		dbPointer := new(DB)
		dbname = fmt.Sprintf("./database/%s", dbname)
		dbPointer.db, err = badger.Open(badger.DefaultOptions(dbname))
		if err != nil {
			log.Fatal(err)
		}
		db = dbPointer
	})
	return db
}
func Close() {
	db.db.Close()
}

func (db *DB) Add(key string, value []byte) {
	db.db.Update(func(txn *badger.Txn) error {
		txn.Set([]byte(key), []byte(value))
		return nil
	})
}

func (db *DB) Get(key string) ([]byte, bool) {
	var buf bytes.Buffer
	ok := false
	err := db.db.View(func(txn *badger.Txn) error {
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

func (db *DB) Remove(key string) {
	db.db.Update(func(txn *badger.Txn) error {
		txn.Delete([]byte(key))
		return nil
	})
}
