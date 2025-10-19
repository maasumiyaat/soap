package database

import (
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

var DB *bolt.DB
var UserBucket = []byte("Users")

func InitDB(dbPath string) error {
	var err error
	DB, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}

	return DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(UserBucket)
		if err != nil {
			return err
		}
		log.Println("bbolt database initialized successfully.")
		return nil
	})
}
