package repository

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mpppk/twitter/internal/io"

	"github.com/ChimeraCoder/anaconda"
	bolt "github.com/mpppk/bbolt"
	"golang.org/x/xerrors"
)

type Repository struct {
	db     *bolt.DB
	logger Logger
}

type Logger interface {
	Println(i ...interface{})
	Printf(format string, i ...interface{})
}

func New(dbPath string, logger Logger) (*Repository, error) {
	db, err := bolt.Open(dbPath, 0666, nil)
	if err != nil {
		return nil, xerrors.Errorf("failed to open db file from %s: %w", dbPath, err)
	}

	if err := createBucketIfNotExists(db, "tweets"); err != nil {
		return nil, xerrors.Errorf("failed to create bucket which named %s: %w", "tweets", err)
	}
	if err := createBucketIfNotExists(db, "maxID"); err != nil {
		return nil, xerrors.Errorf("failed to create bucket which named %s: %w", "maxID", err)
	}

	return &Repository{db: db, logger: logger}, nil
}

func (r *Repository) Close() error {
	err := r.db.Close()
	if err != nil {
		return xerrors.Errorf("failed to close repository: %w", err)
	}
	return nil
}

func (r *Repository) DownloadImageFromDB(downloadDir string) error {
	return r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tweets"))
		if b == nil {
			return fmt.Errorf("failed to retrieve bucket which named %s", "tweets")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var tweet anaconda.Tweet
			if err := json.Unmarshal(v, &tweet); err != nil {
				r.logger.Println(xerrors.Errorf("failed to unmarshal tweet json: %w", err))
				continue
			}

			for i, entityMedia := range tweet.Entities.Media {
				downloadPath, err := io.DownloadEntityMedia(&tweet, &entityMedia, i, downloadDir)
				if err != nil {
					r.logger.Println(xerrors.Errorf("failed to download entity media from %s: %w", entityMedia.Media_url_https, err))
					continue
				}
				if downloadPath != "" {
					r.logger.Printf("media is downloaded to %s\n", downloadPath)
					time.Sleep(10 * time.Second)
				}
			}
			for i, entityMedia := range tweet.ExtendedEntities.Media {
				downloadPath, err := io.DownloadEntityMedia(&tweet, &entityMedia, i, downloadDir)
				if err != nil {
					r.logger.Println(xerrors.Errorf("failed to download extended entity media from %s: %w", entityMedia.Media_url_https, err))
					continue
				}
				if downloadPath != "" {
					r.logger.Printf("extended media is downloaded to %s\n", downloadPath)
					time.Sleep(10 * time.Second)
				}
			}
		}
		return nil
	})
}

func (r *Repository) GetMaxID() (int64, error) {
	var maxID int64
	err := r.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("maxID"))
		maxIDBytes := b.Get([]byte("maxID"))
		mid, bufSize := binary.Varint(maxIDBytes)
		if mid == 0 && bufSize == 0 {
			return fmt.Errorf("buf too small")
		}
		if mid == 0 && bufSize < 0 {
			return fmt.Errorf("value larger than 64 bits (overflow)")
		}
		maxID = mid
		return nil
	})
	if err != nil {
		return 0, xerrors.Errorf("failed to retrieve tweet max id from db: %w", err)
	}
	return maxID, nil
}

func (r *Repository) SaveMaxId(maxId int64) (bool, error) {
	currentMaxID, err := r.GetMaxID()
	if err == nil && currentMaxID <= maxId {
		return false, nil
	}

	idBytes := make([]byte, binary.MaxVarintLen64)
	return true, r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("maxID"))
		binary.PutVarint(idBytes, maxId)
		return b.Put(
			[]byte("maxID"),
			idBytes,
		)
	})
}

func (r *Repository) SaveTweet(tweet *anaconda.Tweet) error {
	idBytes := make([]byte, binary.MaxVarintLen64)
	tweetJsonBytes, err := json.Marshal(tweet)
	if err != nil {
		return err
	}
	return r.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tweets"))
		binary.PutVarint(idBytes, tweet.Id)
		return b.Put(
			idBytes,
			tweetJsonBytes,
		)
	})
}