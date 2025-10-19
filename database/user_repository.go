package database

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/maasumiyaat/soap/model"
	bolt "go.etcd.io/bbolt"
)

func GetUserByID(id int) (*model.User, error) {
	var user model.User
	userIDStr := strconv.Itoa(id)

	err := DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(UserBucket)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", UserBucket)
		}

		v := bucket.Get([]byte(userIDStr))
		if v == nil {
			return fmt.Errorf("user with ID %d not found", id)
		}

		return json.Unmarshal(v, &user)
	})

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func SaveUser(user *model.User) error {
	return DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(UserBucket)
		if err != nil {
			return err
		}
		if user.ID == 0 {
			id, _ := bucket.NextSequence()
			user.ID = int(id)
		}

		buf, err := json.Marshal(user)
		if err != nil {
			return err
		}

		key := []byte(strconv.Itoa(user.ID))
		return bucket.Put(key, buf)
	})
}
