package repository

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/xerrors"

	bolt "github.com/mpppk/bbolt"
)

var MaxIDKey = "maxID"
var MinIDKey = "minID"
var tweetsBucketName = "tweets"
var idBucketName = "maxID"

func IsValidKey(key string) bool {
	return key == MaxIDKey || key == MinIDKey
}

func createBucketIfNotExists(db *bolt.DB, bucketName string) error {
	return db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(bucketName)); err != nil {
			return fmt.Errorf("failed to create %s bucket: %s", bucketName, err)
		}
		return nil
	})
}

func saveInt64(db *bolt.DB, bucketName, key string, value int64) error {
	idBytes := make([]byte, binary.MaxVarintLen64)
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		binary.PutVarint(idBytes, value)
		return b.Put(
			[]byte(key),
			idBytes,
		)
	})
}

func getInt64(db *bolt.DB, bucketName, key string) (int64, error) {
	var value int64
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		maxIDBytes := b.Get([]byte(key))
		mid, bufSize := binary.Varint(maxIDBytes)
		if mid == 0 && bufSize == 0 {
			return fmt.Errorf("buf too small")
		}
		if mid == 0 && bufSize < 0 {
			return fmt.Errorf("value larger than 64 bits (overflow)")
		}
		value = mid
		return nil
	})
	if err != nil {
		return 0, xerrors.Errorf(
			"failed to retrieve int64 value from db: bucket:%s key:%s : %w",
			bucketName,
			key,
			err)
	}
	return value, nil
}
