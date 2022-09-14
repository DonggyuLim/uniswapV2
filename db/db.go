package db

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	dbName  = "Uniswap"
	factory = "factory"
)

var db *bolt.DB

func NewDB() {
	if db == nil {
		dbPointer, err := bolt.Open(dbName, 0600, nil)
		if err != nil {
			log.Fatal(err)
		}
		db = dbPointer
		err = db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte(factory))
			return err
		})
		if err != nil {
			log.Fatal(err)
		}
	}

}

func Close() {
	db.Close()
}

func Add(bucketName, key string, value []byte) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		err := bucket.Put([]byte(key), value)
		return err
	})

	return err
}

func Get(bucketName, key string) ([]byte, error) {
	var data []byte
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		data = bucket.Get([]byte(key))
		return nil
	})
	if err != nil {
		return data, err
	}
	return data, nil
}

func Remove(bucketName, key string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(bucketName))
		err := bucket.Delete([]byte(key))
		return err
	})
	return err
}

func RemoveBucket(bucketName string) error {
	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket([]byte(bucketName))
		return err
	})
	return err
}

// func Add(key string, value []byte) {
// 	db.Update(func(txn *bolt.Tx) error {
// 		bucket := txn.Bucket([]byte())
// 		err := bucket.Put([]byte)
// 	})
// }

// func Get(key string) ([]byte, bool) {
// 	var buf bytes.Buffer
// 	ok := false
// 	err := db.View(func(txn *badger.Txn) error {
// 		item, err := txn.Get([]byte(key))
// 		if err == nil {
// 			item.Value(func(val []byte) error {
// 				buf.Write(val)
// 				ok = true
// 				return nil
// 			})
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return buf.Bytes(), ok
// }

// func Remove(key string) {
// 	db.Update(func(txn *badger.Txn) error {
// 		txn.Delete([]byte(key))
// 		return nil
// 	})
// }
